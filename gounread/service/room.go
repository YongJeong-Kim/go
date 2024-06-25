package service

import (
	"fmt"
	"github.com/google/uuid"
	"gounread/repository"
	"sort"
	"strconv"
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

type ListRoomsByUserIDResult struct {
	RoomID        string    `json:"room_id"`
	Time          time.Time `json:"time"`
	RecentMessage string    `json:"recent_message"`
	UnreadCount   string    `json:"unread_count"`
}

func (s *Service) ListRoomsByUserID(userID string) ([]*ListRoomsByUserIDResult, error) {
	rooms, err := s.Repo.GetRoomsByUserID(userID)
	if err != nil {
		return nil, err
	}

	allReadTimes, err := s.Repo.GetAllRoomsMessageReadTime(userID)
	if err != nil {
		return nil, err
	}

	type getRoomsUnreadMessageCountResult struct {
		RoomID string `json:"room_id"`
		Count  int    `json:"count"`
	}

	var unreadMessageCounts []*getRoomsUnreadMessageCountResult
	for _, t := range allReadTimes {
		cnt, err := s.Repo.GetMessageCountByRoomIDAndSent(t.RoomID, t.ReadTime)
		if err != nil {
			return nil, fmt.Errorf("GetMessageCountByRoomIDAndSent error from GetRoomsUnreadMessageCount. %v", err)
		}

		unreadMessageCounts = append(unreadMessageCounts, &getRoomsUnreadMessageCountResult{
			RoomID: t.RoomID,
			Count:  cnt,
		})
	}

	var result []*ListRoomsByUserIDResult
	for _, r := range rooms {
		for _, c := range unreadMessageCounts {
			if r.RoomID == c.RoomID {
				result = append(result, &ListRoomsByUserIDResult{
					RoomID:        r.RoomID,
					Time:          r.Time,
					RecentMessage: r.RecentMessage,
					UnreadCount:   strconv.Itoa(c.Count),
				})
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.After(result[j].Time)
	})

	return result, nil
}
