package service

import (
	"github.com/scylladb/gocqlx/v2"
	"time"
)

type Message interface {
	GetRecentMessageByRoomID(roomID string) (*GetRecentMessageByRoomIDResult, error)
	GetMessageReadTime(roomID, userID string) (*time.Time, error)
	GetUnreadMessageCount(roomID string, t time.Time) (*int, error)
	GetRecentMessages(roomID string, limit int) []*GetRecentMessagesResult
	GetAllRoomsReadMessageTime(userID string) []*GetAllRoomsReadMessageTimeResult
	GetRoomsUnreadMessageCount(times []*GetAllRoomsReadMessageTimeResult) []*GetRoomsUnreadMessageCountResult
	ReadMessage(roomID, userID string) error
	SendMessage(*SendMessageParam) error
	UpdateRecentMessage(roomID, recentMessage string) error
}

type Room interface {
	CreateRoom(users []string) error
	GetRoomsByUserID(userID string) []*GetRoomsByUserIDResult
}

type Servicer interface {
	Message
	Room
}

type Service struct {
	Session gocqlx.Session
}

func NewService(session gocqlx.Session) *Service {
	return &Service{
		Session: session,
	}
}
