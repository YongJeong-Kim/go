package service

import (
	"fmt"
	"gounread/repository"
	"time"
)

func (s *Service) CreateMessage(param *repository.CreateMessageParam) error {
	return s.Repo.CreateMessage(param)
}

func (s *Service) UpdateRecentMessage(roomID, recentMessage string) error {
	return s.Repo.UpdateRecentMessage(roomID, recentMessage)
}

func (s *Service) GetRecentMessages(roomID string, limit int) []*repository.GetRecentMessagesResult {
	return s.Repo.GetRecentMessages(roomID, limit)
}

func (s *Service) ReadMessage(roomID, userID string) (time.Time, time.Time, error) {
	t, err := s.Repo.GetMessageReadTime(roomID, userID)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("read message get message read time error. %v", err)
	}

	now := time.Now().UTC()
	messages := s.Repo.GetMessagesByRoomIDAndTime(roomID, t, now)

	err = s.Repo.UpdateUnreadMessageBatch(&repository.UpdateUnreadMessageBatchParam{
		UserID:   userID,
		Messages: messages,
	})
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("UpdateUnreadMessageBatch error from ReadMessage. %v", err)
	}

	err = s.Repo.UpdateMessageReadTime(roomID, userID, now)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("UpdateMessageReadTime error from ReadMessage. %v", err)
	}

	return t, now, nil
}

func (s *Service) GetAllRoomsReadMessageTime(userID string) []*repository.GetAllRoomsReadMessageTimeResult {
	return s.Repo.GetAllRoomsReadMessageTime(userID)
}

type GetRoomsUnreadMessageCountResult struct {
	RoomID string `json:"room_id"`
	Count  int    `json:"count"`
}

func (s *Service) GetRoomsUnreadMessageCount(times []*repository.GetAllRoomsReadMessageTimeResult) ([]*GetRoomsUnreadMessageCountResult, error) {
	var result []*GetRoomsUnreadMessageCountResult
	for _, t := range times {
		cnt, err := s.Repo.GetMessageCountByRoomIDAndSent(t.RoomID, t.ReadTime)
		if err != nil {
			return nil, fmt.Errorf("GetMessageCountByRoomIDAndSent error from GetRoomsUnreadMessageCount. %v", err)
		}

		result = append(result, &GetRoomsUnreadMessageCountResult{
			RoomID: t.RoomID,
			Count:  cnt,
		})
	}

	return result, nil
}

func (s *Service) GetRecentMessageByRoomID(roomID string) (*repository.GetRecentMessageByRoomIDResult, error) {
	return s.Repo.GetRecentMessageByRoomID(roomID)
}

func (s *Service) GetMessageReadTime(roomID, userID string) (time.Time, error) {
	return s.Repo.GetMessageReadTime(roomID, userID)
}

func (s *Service) GetUnreadMessages(roomID string, start time.Time, end time.Time) []*repository.GetMessagesByRoomIDAndTimeResult {
	return s.Repo.GetMessagesByRoomIDAndTime(roomID, start, end)
}

func (s *Service) GetUnreadMessageCount(roomID string, t time.Time) (*int, error) {
	return s.Repo.GetUnreadMessageCount(roomID, t)
}

/*func (s *Service) GetMessageByRoomIDAndSent(roomID string, sent time.Time) ([]string, error) {
	return s.Repo.GetMessageByRoomIDAndSent(roomID, sent)
}*/

func (s *Service) UpdateMessageReadTime(roomID, userID string, t time.Time) error {
	return s.Repo.UpdateMessageReadTime(roomID, userID, t)
}
