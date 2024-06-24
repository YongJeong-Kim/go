package repository

import (
	"fmt"
	"time"
)

type GetRoomsByUserIDResult struct {
	RoomID        string    `db:"room_id" json:"room_id"`
	RecentMessage string    `db:"recent_message" json:"recent_message"`
	Time          time.Time `db:"time" json:"time"`
}

func (r *Repository) GetRoomsByUserID(userID string) ([]*GetRoomsByUserIDResult, error) {
	q := `SELECT room_id, recent_message, time FROM room WHERE users CONTAINS ?`
	scanner := r.Session.Query(q, nil).Bind(userID).Iter().Scanner()

	var result []*GetRoomsByUserIDResult
	for scanner.Next() {
		var roomID, recentMessage string
		var t time.Time

		err := scanner.Scan(&roomID, &recentMessage, &t)
		if err != nil {
			return nil, fmt.Errorf("GetRoomsByUserID scan failed. %v", err)
		}
		result = append(result, &GetRoomsByUserIDResult{
			RoomID:        roomID,
			RecentMessage: recentMessage,
			Time:          t,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("GetRoomsByUserID next failed. %v", err)
	}

	if result == nil {
		return nil, fmt.Errorf("user not found. %s", userID)
	}
	return result, nil
}

func (r *Repository) GetUsersByRoomID(roomID string) ([]string, error) {
	var result []string
	q := `SELECT users FROM room WHERE room_id = ?`
	err := r.Session.Query(q, nil).Bind(roomID).Get(&result)
	if err != nil {
		return nil, fmt.Errorf("get users by room id error. %v", err)
	}

	return result, nil
}

func (r *Repository) CreateRoom(roomID string, users []string) error {
	q := `INSERT INTO room(room_id, time, users) VALUES (?, toTimestamp(now()), ?)`
	err := r.Session.Query(q, nil).Bind(roomID, users).Exec()
	if err != nil {
		return fmt.Errorf("create room error. %v", err)
	}
	return nil
}
