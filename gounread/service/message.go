package service

import (
	"fmt"
	"time"
)

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
	return nil
}

func (s *Service) SetRecentMessage(roomID, recentMessage string) error {
	q := `UPDATE room SET recent_message = ?, time = toTimestamp(now()) WHERE id = ?`
	err := s.Session.Query(q, []string{}).Bind(recentMessage, roomID).ExecRelease()
	if err != nil {
		return fmt.Errorf("add message unread failed. %v", err)
	}
	return nil
}

func (s *Service) IncrementUnreadMessage(roomID, userID string, inc int) error {
	q := `UPDATE message_unread SET unread = unread + ? WHERE room_id = ? AND user_id = ?`
	err := s.Session.Query(q, []string{}).Bind(inc, roomID, userID).ExecRelease()
	if err != nil {
		return fmt.Errorf("add message unread failed. %v", err)
	}
	return nil
}

func (s *Service) GetUnreadCount(roomID, userID string) (int, error) {
	var result int
	q := `SELECT unread FROM message_unread WHERE room_id = ? AND user_id = ? LIMIT 1`
	err := s.Session.Query(q, []string{}).Bind(roomID, userID).GetRelease(&result)
	if err != nil {
		return 0, fmt.Errorf("reset unread failed. %v", err)
	}

	return result, nil
}

type GetRecentMessagesResult struct {
	Sent   time.Time `db:"sent" json:"sent"`
	Msg    string    `db:"msg" json:"msg"`
	Sender string    `db:"sender" json:"sender"`
}

func (s *Service) GetRecentMessages(roomID string, limit int) []*GetRecentMessagesResult {
	q := `SELECT sent, msg, sender FROM message WHERE room_id = ? LIMIT ?`
	messages := s.Session.Query(q, []string{}).Bind(roomID, limit).Iter()

	var result []*GetRecentMessagesResult
	for {
		var m GetRecentMessagesResult
		if !messages.StructScan(&m) {
			break
		}
		result = append(result, &m)
	}

	return result
}
