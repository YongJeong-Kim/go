package repository

import (
	"fmt"
	"github.com/gocql/gocql"
	"time"
)

type CreateMessageParam struct {
	RoomID      string
	Sender      string
	Sent        time.Time
	Message     string
	UnreadUsers []string
}

func (r *Repository) CreateMessage(param *CreateMessageParam) error {
	q := `INSERT INTO message (room_id, sender, msg, sent, unread) VALUES (?, ?, ?, ?, ?)`
	err := r.Session.Query(q, nil).Bind(param.RoomID, param.Sender, param.Message, param.Sent, param.UnreadUsers).ExecRelease()
	if err != nil {
		return fmt.Errorf("send message failed. %v", err)
	}
	return nil
}

func (r *Repository) UpdateRecentMessage(roomID, recentMessage string) error {
	q := `UPDATE room SET recent_message = ?, time = toTimestamp(now()) WHERE room_id = ?`
	err := r.Session.Query(q, nil).Bind(recentMessage, roomID).ExecRelease()
	if err != nil {
		return fmt.Errorf("update recent message failed. %v", err)
	}
	return nil
}

func (r *Repository) GetUnreadMessageCount(roomID string, t time.Time) (int, error) {
	var cnt int
	q := `SELECT COUNT(room_id) AS cnt FROM message WHERE room_id = ? AND sent >= ? AND sent <= toTimestamp(now())`
	err := r.Session.Query(q, nil).Bind(roomID, t).Get(&cnt)
	if err != nil {
		return 0, fmt.Errorf("get message status unread count failed. %v", err)
	}
	return cnt, nil
}

type GetMessagesByRoomIDAndTimeResult struct {
	RoomID string    `db:"room_id" json:"room_id"`
	Sent   time.Time `db:"sent" json:"sent"`
	Msg    string    `db:"msg" json:"msg"`
	Sender string    `db:"sender" json:"sender"`
	Unread []string  `db:"unread" json:"unread"`
}

func (r *Repository) GetMessagesByRoomIDAndTime(roomID string, start time.Time, end time.Time) ([]*GetMessagesByRoomIDAndTimeResult, error) {
	q := `SELECT room_id, sent, msg, sender, unread FROM message WHERE room_id = ? AND sent >= ? AND sent <= ?`
	scanner := r.Session.Query(q, nil).Bind(roomID, start, end).Iter().Scanner()

	var result []*GetMessagesByRoomIDAndTimeResult
	for scanner.Next() {
		var (
			rID    string
			sent   time.Time
			msg    string
			sender string
			unread []string
		)

		err := scanner.Scan(&rID, &sent, &msg, &sender, &unread)
		if err != nil {
			return nil, fmt.Errorf("get messages by room id and time scan failed. %v", err)
		}

		result = append(result, &GetMessagesByRoomIDAndTimeResult{
			RoomID: roomID,
			Msg:    msg,
			Sent:   sent,
			Sender: sender,
			Unread: unread,
		})
	}
	return result, nil
}

func (r *Repository) GetMessageReadTime(roomID, userID string) (time.Time, error) {
	var t time.Time
	q := `SELECT read_time FROM message_read WHERE room_id = ? AND user_id = ?`
	err := r.Session.Query(q, nil).Bind(roomID, userID).GetRelease(&t)
	if err != nil {
		return time.Time{}, fmt.Errorf("get message read time failed. %v", err)
	}
	return t, nil
}

type GetRecentMessageByRoomIDResult struct {
	RoomID        string `db:"room_id" json:"room_id"`
	RecentMessage string `db:"recent_message" json:"recent_message"`
}

func (r *Repository) GetRecentMessageByRoomID(roomID string) (*GetRecentMessageByRoomIDResult, error) {
	var rr GetRecentMessageByRoomIDResult
	q := `SELECT room_id, recent_message FROM room WHERE room_id = ? LIMIT 1`
	err := r.Session.Query(q, nil).Bind(roomID).Get(&rr)
	if err != nil {
		return nil, fmt.Errorf("get room recent message failed. %v", err)
	}
	return &rr, nil
}

type GetAllRoomsReadMessageTimeResult struct {
	RoomID   string    `db:"room_id" json:"room_id"`
	ReadTime time.Time `db:"read_time" json:"read_time"`
}

func (r *Repository) GetAllRoomsReadMessageTime(userID string) ([]*GetAllRoomsReadMessageTimeResult, error) {
	q := `SELECT room_id, read_time FROM message_read_by_user WHERE user_id = ?`
	scanner := r.Session.Query(q, []string{}).Bind(userID).Iter().Scanner()

	var result []*GetAllRoomsReadMessageTimeResult
	for scanner.Next() {
		var roomID string
		var readTime time.Time

		err := scanner.Scan(&roomID, &readTime)
		if err != nil {
			return nil, fmt.Errorf("get all rooms read message time error. %v", err)
		}

		result = append(result, &GetAllRoomsReadMessageTimeResult{
			RoomID:   roomID,
			ReadTime: readTime,
		})
	}

	return result, nil
}

type UpdateUnreadMessageBatchParam struct {
	UserID   string
	Messages []*GetMessagesByRoomIDAndTimeResult
}

func (r *Repository) UpdateUnreadMessageBatch(param *UpdateUnreadMessageBatchParam) error {
	// not working range update
	q := `UPDATE message SET unread = unread - ? WHERE room_id = ? AND sent = ?`
	// Be careful. batching queries with different partition keys, anti-pattern -> poor performance
	// https://docs.datastax.com/en/cql/3.3/cql/cql_using/useBatchBadExample.html
	// In this case, same partition key

	batch := r.Session.NewBatch(gocql.LoggedBatch)
	for _, m := range param.Messages {
		batch.Query(q, []string{param.UserID}, m.RoomID, m.Sent)
	}

	err := r.Session.ExecuteBatch(batch)
	if err != nil {
		return fmt.Errorf("read message batch failed. %v", err)
	}
	return nil
}

func (r *Repository) UpdateMessageReadTime(roomID string, userID string, now time.Time) error {
	q := `UPDATE message_read SET read_time = ? WHERE room_id = ? AND user_id = ?`
	err := r.Session.Query(q, nil).Bind(now, roomID, userID).ExecRelease()
	if err != nil {
		return fmt.Errorf("update message read time failed. %v", err)
	}
	return nil
}

type GetRecentMessagesResult struct {
	RoomID string    `db:"room_id" json:"room_id"`
	Sent   time.Time `db:"sent" json:"sent"`
	Msg    string    `db:"msg" json:"msg"`
	Sender string    `db:"sender" json:"sender"`
}

func (r *Repository) GetRecentMessages(roomID string, limit int) ([]*GetRecentMessagesResult, error) {
	q := `SELECT room_id, sent, msg, sender FROM message WHERE room_id = ? LIMIT ?`
	scanner := r.Session.Query(q, []string{}).Bind(roomID, limit).Iter().Scanner()

	var result []*GetRecentMessagesResult
	for scanner.Next() {
		var (
			rID    string
			sent   time.Time
			msg    string
			sender string
		)

		err := scanner.Scan(&rID, &sent, &msg, &sender)
		if err != nil {
			return nil, fmt.Errorf("get recent messages error. %v", err)
		}

		result = append(result, &GetRecentMessagesResult{
			RoomID: rID,
			Sender: sender,
			Sent:   sent,
			Msg:    msg,
		})
	}

	return result, nil
}

func (r *Repository) GetMessageCountByRoomIDAndSent(roomID string, readTime time.Time) (int, error) {
	var cnt int
	q := `SELECT COUNT(room_id) AS cnt FROM message WHERE room_id = ? AND sent >= ? AND sent <= toTimestamp(now())`
	err := r.Session.Query(q, []string{}).Bind(roomID, readTime).GetRelease(&cnt)
	if err != nil {
		return 0, fmt.Errorf("get message count error. %v", err)
	}

	return cnt, nil
}
