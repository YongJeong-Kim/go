package service

import (
	"gounread/repository"
	"time"
)

type Message interface {
	SendMessage(param *repository.CreateMessageParam) error
	UpdateRecentMessage(roomID, recentMessage string) error
	GetRecentMessages(roomID string, limit int) []*repository.GetRecentMessagesResult
	ReadMessage(roomID, userID string) error
	GetAllRoomsReadMessageTime(userID string) []*repository.GetAllRoomsReadMessageTimeResult
	GetRoomsUnreadMessageCount(times []*repository.GetAllRoomsReadMessageTimeResult) ([]*GetRoomsUnreadMessageCountResult, error)
	GetRecentMessageByRoomID(roomID string) (*repository.GetRecentMessageByRoomIDResult, error)
	GetMessageReadTime(roomID, userID string) (time.Time, error)
	GetUnreadMessages(roomID string, t time.Time) []*repository.GetMessagesByRoomIDAndTimeResult
	GetUnreadMessageCount(roomID string, t time.Time) (*int, error)
}

type Room interface {
	CreateRoom(users []string) error
	GetRoomsByUserID(userID string) []*repository.GetRoomsByUserIDResult
	GetUsersByRoomID(roomID string) ([]string, error)
	JoinRoom(roomID, userID string) ([]*repository.GetMessagesByRoomIDAndTimeResult, error)
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
