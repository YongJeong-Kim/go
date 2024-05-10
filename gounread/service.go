package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/scylladb/gocqlx/v2"
	"time"
)

type Service struct {
	Session gocqlx.Session
	Redis   *redis.Client
}

type Servicer interface {
	GetRoomsByUserID(userID string) []*GetRoomsByUserIDResult
	SendMessage(*SendMessageParam) error
	SetRecentMessage(roomID, recentMessage string) error
	AddMessageUnread(roomID, userID string) error
}

func NewService(session gocqlx.Session, redis *redis.Client) *Service {
	return &Service{
		Session: session,
		Redis:   redis,
	}
}

type GetRoomsByUserIDResult struct {
	ID            string    `db:"id" json:"id"`
	RecentMessage string    `db:"recent_message" json:"recent_message"`
	Time          time.Time `db:"time" json:"time"`
}

func (s *Service) GetRoomsByUserID(userID string) []*GetRoomsByUserIDResult {
	q := `SELECT id, recent_message, time FROM room WHERE users CONTAINS ?`
	rooms := s.Session.Query(q, []string{"id", "recent_message", "time"}).Bind(userID).Iter()
	defer rooms.Close()

	var result []*GetRoomsByUserIDResult
	for {
		var r GetRoomsByUserIDResult
		if !rooms.StructScan(&r) {
			break
		}
		result = append(result, &r)
	}
	return result
}

type SendMessageParam struct {
	RoomID  string `json:"room_id"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
	After   func() error
}

func (s *Service) SendMessage(param *SendMessageParam) error {
	q := `INSERT INTO message (room_id, sender, msg, sent) VALUES (?, ?, ?, toTimestamp(now()))`
	err := s.Session.Query(q, []string{"asd"}).Bind(param.RoomID, param.Sender, param.Message).ExecRelease()
	if err != nil {
		return fmt.Errorf("send message failed. %v", err)
	}
	return param.After()
}

func (s *Service) SetRecentMessage(roomID, recentMessage string) error {
	q := `UPDATE room SET recent_message = ? WHERE id = ?`
	err := s.Session.Query(q, []string{}).Bind(recentMessage, roomID).ExecRelease()
	if err != nil {
		return fmt.Errorf("add message unread failed. %v", err)
	}
	return nil
}

func (s *Service) AddMessageUnread(roomID, userID string) error {
	q := `UPDATE message_unread SET unread = unread + 1 WHERE room_id = ? AND user_id = ?`
	err := s.Session.Query(q, []string{}).Bind(roomID, userID).ExecRelease()
	if err != nil {
		return fmt.Errorf("add message unread failed. %v", err)
	}
	return nil
}
