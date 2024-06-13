package service

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type GetRoomsByUserIDResult struct {
	RoomID        string    `db:"room_id" json:"room_id"`
	RecentMessage string    `db:"recent_message" json:"recent_message"`
	Time          time.Time `db:"time" json:"time"`
}

func (s *Service) GetRoomsByUserID(userID string) []*GetRoomsByUserIDResult {
	q := `SELECT room_id, recent_message, time FROM room WHERE users CONTAINS ?`
	rooms := s.Session.Query(q, []string{"room_id", "recent_message", "time"}).Bind(userID).Iter()
	defer rooms.Close()

	var result []*GetRoomsByUserIDResult
	for {
		var r GetRoomsByUserIDResult
		if !rooms.StructScan(&r) {
			break
		}
		result = append(result, &r)
	}
	return result
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

	err := s.Session.Query(`INSERT INTO room(id, time, users) VALUES (uuid(), toTimestamp(now()), ?)`, nil).Bind(users).Exec()
	if err != nil {
		return fmt.Errorf("create room error. %v", err)
	}
	return nil
}
