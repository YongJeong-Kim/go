package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gounread/repository"
	"net/http"
	"strconv"
	"time"
)

type Payload struct {
	RoomID      string    `json:"room_id"`
	Sender      string    `json:"sender"`
	Sent        time.Time `json:"sent"`
	Message     string    `json:"message"`
	UnreadCount int       `json:"unread_count"`
}

func (s *Server) SendMessage(c *gin.Context) {
	var reqURI struct {
		RoomID string `uri:"room_id" binding:"required"`
	}
	if err := c.ShouldBindUri(&reqURI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var reqJSON struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.Request.Header.Get("user")

	users, err := s.Service.GetUsersByRoomID(reqURI.RoomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	for i := range users {
		if userID == users[i] {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}

	now := time.Now().UTC()
	err = s.Service.CreateMessage(&repository.CreateMessageParam{
		RoomID:      reqURI.RoomID,
		Sender:      userID,
		Sent:        now,
		Message:     reqJSON.Message,
		UnreadUsers: users,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = s.Service.UpdateRecentMessage(reqURI.RoomID, reqJSON.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = s.Service.UpdateMessageReadTime(reqURI.RoomID, userID, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	payload := &Payload{
		RoomID:      reqURI.RoomID,
		Sender:      c.Request.Header.Get("user"),
		Sent:        now,
		Message:     reqJSON.Message,
		UnreadCount: len(users),
	}
	b, _ := json.Marshal(payload)

	err = s.Nats.Publish("room."+reqURI.RoomID, b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s%v", "room publish error. ", err),
		})
		return
	}

	err = s.Nats.Publish("lobby."+reqURI.RoomID, b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s%v", "lobby publish error. ", err),
		})
		return
	}

	err = s.Nats.Publish("focus.lobby."+reqURI.RoomID, b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s%v", "focus lobby publish error. ", err),
		})
		return
	}

	c.Status(http.StatusCreated)
}

func (s *Server) ReadMessage(c *gin.Context) {
	// select recent message limit 10 order by time desc
	var reqURI struct {
		RoomID string `uri:"room_id" binding:"required"`
	}
	if err := c.ShouldBindUri(&reqURI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.Request.Header.Get("user")

	start, end, err := s.Service.ReadMessage(reqURI.RoomID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	messages := s.Service.GetUnreadMessages(reqURI.RoomID, start, end)
	if len(messages) == 0 {
		recentMessages := s.Service.GetRecentMessages(reqURI.RoomID, 10)
		c.JSON(http.StatusOK, recentMessages)
		return
	}

	//s.Nats.Publish("", []byte("1"))
	c.JSON(http.StatusOK, messages)
}

type GetRoomStatusInLobbyResponse struct {
	RoomID        string `json:"room_id"`
	RecentMessage string `json:"recent_message"`
	UnreadCount   string `json:"unread_count"`
}

func (s *Server) GetRoomStatusInLobby(c *gin.Context) {
	var reqURI struct {
		RoomID string `uri:"room_id" binding:"required"`
	}
	if err := c.ShouldBindUri(&reqURI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.Request.Header.Get("user")

	recentMessage, err := s.Service.GetRecentMessageByRoomID(reqURI.RoomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	t, err := s.Service.GetMessageReadTime(reqURI.RoomID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	count, err := s.Service.GetUnreadMessageCount(reqURI.RoomID, t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := &GetRoomStatusInLobbyResponse{
		RoomID:        reqURI.RoomID,
		RecentMessage: recentMessage.RecentMessage,
		UnreadCount:   strconv.Itoa(*count),
	}

	c.JSON(http.StatusOK, response)
}
