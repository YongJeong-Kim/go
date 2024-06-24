package service

import (
	"gounread/repository"
	"time"
)

type Message interface {
	CreateMessage(param *repository.CreateMessageParam) error
	GetAllRoomsReadMessageTime(userID string) ([]*repository.GetAllRoomsReadMessageTimeResult, error)
	GetRecentMessageByRoomID(roomID string) (*repository.GetRecentMessageByRoomIDResult, error)
	GetRecentMessages(roomID string, limit int) ([]*repository.GetRecentMessagesResult, error)
	GetRoomsUnreadMessageCount(times []*repository.GetAllRoomsReadMessageTimeResult) ([]*GetRoomsUnreadMessageCountResult, error)
	GetUnreadMessageCount(roomID, userID string) (int, error)
	GetUnreadMessages(roomID string, start time.Time, end time.Time) ([]*repository.GetMessagesByRoomIDAndTimeResult, error)
	ReadMessage(roomID, userID string) (time.Time, time.Time, error)
	UpdateMessageReadTime(roomID, userID string, t time.Time) error
	UpdateRecentMessage(roomID, recentMessage string) error
}

type Room interface {
	CreateRoom(users []string) error
	GetRoomsByUserID(userID string) ([]*repository.GetRoomsByUserIDResult, error)
	GetUsersByRoomID(roomID string) ([]string, error)
}

type Servicer interface {
	Message
	Room
}

type Service struct {
	Repo repository.Repositorier
}

func NewService(repo repository.Repositorier) *Service {
	return &Service{
		Repo: repo,
	}
}
