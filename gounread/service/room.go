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

	newRoomID := uuid.NewString()
	err := s.Repo.CreateRoom(newRoomID, users)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, u := range users {
		err := s.Repo.UpdateMessageReadTime(newRoomID, u, now)
		if err != nil {
			return fmt.Errorf("create room update message read time failed. %v", err)
		}
	}

	return nil
}

func (s *Service) GetUsersByRoomID(roomID string) ([]string, error) {
	return s.Repo.GetUsersByRoomID(roomID)
}
