package service

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type SendMessageParam struct {
	RoomID  string `json:"room_id"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
	After   func() error
}

func (s *Service) SendMessage(param *SendMessageParam) error {
	q := `INSERT INTO message (room_id, sender, msg, sent) VALUES (?, ?, ?, toTimestamp(now()))`
	err := s.Session.Query(q, []string{"asd"}).Bind(param.RoomID, param.Sender, param.Message).ExecRelease()
	if err != nil {
		return fmt.Errorf("send message failed. %v", err)
	}
	return nil
}

func (s *Service) SetRecentMessage(roomID, recentMessage string) error {
	q := `UPDATE room SET recent_message = ?, time = toTimestamp(now()) WHERE id = ?`
	err := s.Session.Query(q, []string{}).Bind(recentMessage, roomID).ExecRelease()
	if err != nil {
		return fmt.Errorf("add message unread failed. %v", err)
	}
	return nil
}

/*func (s *Service) IncrementUnreadMessage(roomID, userID string, inc int) error {
	q := `UPDATE message_unread SET unread = unread + ? WHERE room_id = ? AND user_id = ?`
	err := s.Session.Query(q, []string{}).Bind(inc, roomID, userID).ExecRelease()
	if err != nil {
		return fmt.Errorf("add message unread failed. %v", err)
	}
	return nil
}*/

/*func (s *Service) GetUnreadCount(roomID, userID string) (int, error) {
	var result int
	q := `SELECT unread FROM message_unread WHERE room_id = ? AND user_id = ? LIMIT 1`
	err := s.Session.Query(q, []string{}).Bind(roomID, userID).GetRelease(&result)
	if err != nil {
		return 0, fmt.Errorf("reset unread failed. %v", err)
	}

	return result, nil
}*/

type GetRecentMessagesResult struct {
	Sent   time.Time `db:"sent" json:"sent"`
	Msg    string    `db:"msg" json:"msg"`
	Sender string    `db:"sender" json:"sender"`
}

func (s *Service) GetRecentMessages(roomID string, limit int) []*GetRecentMessagesResult {
	q := `SELECT sent, msg, sender FROM message WHERE room_id = ? LIMIT ?`
	messages := s.Session.Query(q, []string{}).Bind(roomID, limit).Iter()

	var result []*GetRecentMessagesResult
	for {
		var m GetRecentMessagesResult
		if !messages.StructScan(&m) {
			break
		}
		result = append(result, &m)
	}

	return result
}

func (s *Service) ReadMessage(roomID, userID string) error {
	q := `INSERT INTO message_read(room_id, user_id, read_time) VALUES (?, ?, toTimestamp(now()))`
	err := s.Session.Query(q, []string{}).Bind(roomID, userID).ExecRelease()
	if err != nil {
		return fmt.Errorf("read message failed. %v", err)
	}
	return nil
}

type GetAllRoomsReadMessageTimeResult struct {
	RoomID   string    `db:"room_id" json:"room_id"`
	ReadTime time.Time `db:"read_time" json:"read_time"`
}

func (s *Service) GetAllRoomsReadMessageTime(userID string) []*GetAllRoomsReadMessageTimeResult {
	q := `SELECT room_id, read_time FROM message_read_by_user WHERE user_id = ?`
	counts := s.Session.Query(q, []string{}).Bind(userID).Iter()

	var result []*GetAllRoomsReadMessageTimeResult
	for {
		var r GetAllRoomsReadMessageTimeResult
		if !counts.StructScan(&r) {
			break
		}
		result = append(result, &r)
	}

	return result
}

type GetRoomsUnreadMessageCountResult struct {
	RoomID string `json:"room_id"`
	Count  int    `json:"count"`
}

func (s *Service) GetRoomsUnreadMessageCount(times []*GetAllRoomsReadMessageTimeResult) []*GetRoomsUnreadMessageCountResult {
	var result []*GetRoomsUnreadMessageCountResult
	for _, t := range times {
		var cnt int
		q := `SELECT COUNT(room_id) AS cnt FROM message WHERE room_id = ? AND sent >= ? AND sent <= toTimestamp(now())`
		err := s.Session.Query(q, []string{}).Bind(t.RoomID, t.ReadTime).GetRelease(&cnt)
		if err != nil {
			log.Println("get message count error. ", err)
			return nil
		}
		result = append(result, &GetRoomsUnreadMessageCountResult{
			RoomID: t.RoomID,
			Count:  cnt,
		})
	}

	return result
}

type GetRoomStatusInLobbyResult struct {
	RoomID        string `json:"room_id"`
	RecentMessage string `json:"recent_message"`
	UnreadCount   string `json:"unread_count"`
}

func (s *Service) GetRoomStatusInLobby(roomID, userID string) (*GetRoomStatusInLobbyResult, error) {
	var r struct {
		RoomID        string `db:"id" json:"room_id"`
		RecentMessage string `db:"recent_message" json:"recent_message"`
	}
	q := `SELECT id, recent_message FROM room WHERE id = ?`
	err := s.Session.Query(q, []string{}).Bind(roomID).Get(&r)
	if err != nil {
		return nil, fmt.Errorf("get room status recent message failed. %v", err)
	}

	// get read message time
	var t time.Time
	q = `SELECT read_time FROM message_read WHERE room_id = ? AND user_id = ?`
	err = s.Session.Query(q, []string{}).Bind(roomID, userID).GetRelease(&t)
	if err != nil {
		return nil, fmt.Errorf("get room status read time failed. %v", err)
	}

	// calc message time between read message time
	var cnt int
	q = `SELECT COUNT(room_id) AS cnt FROM message WHERE room_id = ? AND sent >= ? AND sent <= toTimestamp(now())`
	err = s.Session.Query(q, []string{}).Bind(roomID, t).Get(&cnt)
	if err != nil {
		return nil, fmt.Errorf("get message status unread count failed. %v", err)
	}

	return &GetRoomStatusInLobbyResult{
		RoomID:        r.RoomID,
		RecentMessage: r.RecentMessage,
		UnreadCount:   strconv.Itoa(cnt),
	}, nil
}
