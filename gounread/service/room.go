package service

import (
	"fmt"
	"github.com/google/uuid"
	"gounread/repository"
	"time"
)

func (s *Service) GetRoomsByUserID(userID string) ([]*repository.GetRoomsByUserIDResult, error) {
	return s.Repo.GetRoomsByUserID(userID)
}

func (s *Service) CreateRoom(users []string) error {
	if len(users) < 2 {
		return fmt.Errorf("minimal user count 2. but: %d", len(users))
	}

	for _, u := range users {
		err := uuid.Validate(u)
		if err != nil {
			return fmt.Errorf("invalid user id: %v. %s", err, u)
		}
	}

	err := s.Repo.CreateRoom(users)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetUsersByRoomID(roomID string) ([]string, error) {
	return s.Repo.GetUsersByRoomID(roomID)
}

func (s *Service) JoinRoom(roomID, userID string) ([]*repository.GetMessagesByRoomIDAndTimeResult, error) {
	previousReadTime, err := s.Repo.GetMessageReadTime(roomID, userID)
	if err != nil {
		return nil, err
	}

	// select between prevvious read time and now
	now := time.Now().UTC()
	unreadMessages := s.Repo.GetMessagesByRoomIDAndTime(roomID, previousReadTime, now)

	// update message delete unread user
	err = s.Repo.UpdateUnreadMessageBatch(&repository.UpdateUnreadMessageBatchParam{
		UserID:   userID,
		Messages: unreadMessages,
	})
	if err != nil {
		return nil, err
	}

	// new select between previous read time and now
	updatedMessages := s.Repo.GetMessagesByRoomIDAndTime(roomID, previousReadTime, now)

	// update message read time
	err = s.Repo.UpdateMessageReadTime(roomID, userID, now)
	if err != nil {
		return nil, err
	}
	return updatedMessages, nil
}
