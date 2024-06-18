package repository

import (
	"github.com/scylladb/gocqlx/v2"
	"time"
)

type Message interface {
	CreateMessage(param *CreateMessageParam) error
	GetAllRoomsReadMessageTime(userID string) []*GetAllRoomsReadMessageTimeResult
	GetMessageReadTime(roomID, userID string) (time.Time, error)
	GetMessageCountByRoomIDAndSent(roomID string, readTime time.Time) (int, error)
	GetRecentMessageByRoomID(roomID string) (*GetRecentMessageByRoomIDResult, error)
	GetRecentMessages(roomID string, limit int) []*GetRecentMessagesResult
	GetUnreadMessageCount(roomID string, t time.Time) (*int, error)
	GetMessagesByRoomIDAndTime(roomID string, start time.Time, end time.Time) []*GetMessagesByRoomIDAndTimeResult
	UpdateMessageReadTime(roomID string, userID string, now time.Time) error
	UpdateRecentMessage(roomID, recentMessage string) error
	UpdateUnreadMessageBatch(param *UpdateUnreadMessageBatchParam) error
}

type Room interface {
	CreateRoom(roomID string, users []string) error
	GetRoomsByUserID(userID string) ([]*GetRoomsByUserIDResult, error)
	GetUsersByRoomID(roomID string) ([]string, error)
}

type Repository struct {
	Session gocqlx.Session
}

type Repositorier interface {
	Message
	Room
}

func NewRepository(session gocqlx.Session) *Repository {
	return &Repository{
		Session: session,
	}
}
