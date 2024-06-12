package service

import (
	"github.com/gocql/gocql"
	"time"
)

type GetRoomsByUserIDResult struct {
	ID            string    `db:"id" json:"id"`
	RecentMessage string    `db:"recent_message" json:"recent_message"`
	Time          time.Time `db:"time" json:"time"`
}

func (s *Service) GetRoomsByUserID(userID string) []*GetRoomsByUserIDResult {
	q := `SELECT id, recent_message, time FROM room WHERE users CONTAINS ?`
	rooms := s.Session.Query(q, []string{"id", "recent_message", "time"}).Bind(userID).Iter()
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

func (s *Service) CreateRoomMembers() {
	batch := s.Session.NewBatch(gocql.LoggedBatch)
	batch.Query(`INSERT INTO room(id, time, users) VALUES ()`)
	s.Session.ExecuteBatch()
}
