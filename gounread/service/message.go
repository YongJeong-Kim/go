package service

import (
	"fmt"
	"gounread/repository"
	"time"
)

type SendMessageParam struct {
	RoomID  string
	Sender  string
	Sent    time.Time
	Message string
}

type Payload struct {
	RoomID  string    `json:"room_id"`
	Sender  string    `json:"sender"`
	Sent    time.Time `json:"sent"`
	Message string    `json:"message"`
	Unread  []string  `json:"unread"`
}

func (s *Service) SendMessage(param *SendMessageParam) (*Payload, error) {
	users, err := s.Repo.GetUsersByRoomID(param.RoomID)
	if err != nil {
		return nil, err
	}

	for i := range users {
		if param.Sender == users[i] {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}

	err = s.Repo.CreateMessage(&repository.CreateMessageParam{
		RoomID:      param.RoomID,
		Sender:      param.Sender,
		Message:     param.Message,
		Sent:        param.Sent,
		UnreadUsers: users,
	})
	if err != nil {
		return nil, err
	}

	err = s.Repo.UpdateRecentMessage(param.RoomID, param.Message)
	if err != nil {
		return nil, err
	}

	err = s.Repo.UpdateMessageReadTime(param.RoomID, param.Sender, param.Sent)
	if err != nil {
		return nil, err
	}

	return &Payload{
		RoomID:  param.RoomID,
		Sender:  param.Sender,
		Sent:    param.Sent,
		Message: param.Message,
		Unread:  users,
	}, nil
}

func (s *Service) GetRecentMessages(roomID string, limit int) ([]*repository.GetRecentMessagesResult, error) {
	return s.Repo.GetRecentMessages(roomID, limit)
}

func (s *Service) ReadMessage(roomID, userID string) ([]*repository.GetMessagesByRoomIDAndTimeResult, error) {
	t, err := s.Repo.GetMessageReadTime(roomID, userID)
	if err != nil {
		return nil, fmt.Errorf("read message get message read time error. %v", err)
	}

	now := time.Now().UTC()
	messages, err := s.Repo.GetMessagesByRoomIDAndTime(roomID, t, now)
	if err != nil {
		return nil, fmt.Errorf("get messages by room id and time error. %v", err)
	}
	if messages != nil {
		err = s.Repo.UpdateUnreadMessageBatch(&repository.UpdateUnreadMessageBatchParam{
			UserID:   userID,
			Messages: messages,
		})
		if err != nil {
			return nil, fmt.Errorf("UpdateUnreadMessageBatch error from ReadMessage. %v", err)
		}

		err = s.Repo.UpdateMessageReadTime(roomID, userID, now)
		if err != nil {
			return nil, fmt.Errorf("UpdateMessageReadTime error from ReadMessage. %v", err)
		}
	}

	// get messages after update message read time
	messages, err = s.Repo.GetMessagesByRoomIDAndTime(roomID, t, now)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
