package service

import (
	"github.com/scylladb/gocqlx/v2"
)

type Message interface {
	SendMessage(*SendMessageParam) error
	SetRecentMessage(roomID, recentMessage string) error
	GetRecentMessages(roomID string, limit int) []*GetRecentMessagesResult
	GetAllRoomsReadMessageTime(userID string) []*GetAllRoomsReadMessageTimeResult
	GetRoomsUnreadMessageCount(times []*GetAllRoomsReadMessageTimeResult) []*GetRoomsUnreadMessageCountResult
	ReadMessage(roomID, userID string) error
}

type Room interface {
	GetRoomsByUserID(userID string) []*GetRoomsByUserIDResult
	GetRoomStatusInLobby(roomID, userID string) (*GetRoomStatusInLobbyResult, error)
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
