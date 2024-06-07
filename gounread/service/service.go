package service

import (
	"github.com/scylladb/gocqlx/v2"
)

type Service struct {
	Session gocqlx.Session
}

type Servicer interface {
	GetRoomsByUserID(userID string) []*GetRoomsByUserIDResult
	SendMessage(*SendMessageParam) error
	SetRecentMessage(roomID, recentMessage string) error
	IncrementUnreadMessage(roomID, userID string, inc int) error
	GetRoomUsersByRoomID(roomID string) ([]string, error)
	GetUsersByRoomID(roomID string) []string
	GetUnreadCount(roomID, userID string) (int, error)
	GetRecentMessages(roomID string, limit int) []*GetRecentMessagesResult
}

func NewService(session gocqlx.Session) *Service {
	return &Service{
		Session: session,
	}
}
