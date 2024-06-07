package service

import "time"

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

func (s *Service) GetUsersByRoomID(roomID string) []string {
	q := `SELECT users FROM room WHERE id = ?`
	users := s.Session.Query(q, []string{"users"}).Bind(roomID).Iter()

	var result []string
	for {
		if !users.Scan(&result) {
			break
		}
	}

	return result
}

func (s *Service) GetRoomUsersByRoomID(roomID string) ([]string, error) {
	var result []string
	q := `SELECT users FROM room WHERE id = ?`
	err := s.Session.Query(q, []string{"users"}).Bind(roomID).Get(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
