package repository

import (
	"fmt"
	"log"
	"time"
)

type GetRoomsByUserIDResult struct {
	RoomID        string    `db:"room_id" json:"room_id"`
	RecentMessage string    `db:"recent_message" json:"recent_message"`
	Time          time.Time `db:"time" json:"time"`
}

func (r *Repository) GetRoomsByUserID(userID string) ([]*GetRoomsByUserIDResult, error) {
	q := `SELECT room_id, recent_message, time FROM room WHERE users CONTAINS ?`
	rooms := r.Session.Query(q, []string{"room_id", "recent_message", "time"}).Bind(userID).Iter()
	defer func() {
		err := rooms.Close()
		if err != nil {
			log.Println("GetRoomsByUserID close failed. ", err)
		}
	}()

	var result []*GetRoomsByUserIDResult
	for {
		var r GetRoomsByUserIDResult
		if !rooms.StructScan(&r) {
			break
		}
		result = append(result, &r)
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

func (r *Repository) CreateRoom(users []string) error {
	q := `INSERT INTO room(room_id, time, users) VALUES (uuid(), toTimestamp(now()), ?)`
	err := r.Session.Query(q, nil).Bind(users).Exec()
	if err != nil {
		return fmt.Errorf("create room error. %v", err)
	}
	return nil
}
