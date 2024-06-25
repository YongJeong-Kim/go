package service

import (
	"gounread/repository"
)

type Message interface {
	SendMessage(param *SendMessageParam) (*Payload, error)
	GetRecentMessages(roomID string, limit int) ([]*repository.GetRecentMessagesResult, error)
	ReadMessage(roomID, userID string) ([]*repository.GetMessagesByRoomIDAndTimeResult, error)
}

type Room interface {
	CreateRoom(users []string) error
	GetRoomsByUserID(userID string) ([]*repository.GetRoomsByUserIDResult, error)
	ListRoomsByUserID(userID string) ([]*ListRoomsByUserIDResult, error)
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
