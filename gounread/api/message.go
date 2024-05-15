package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gounread/service"
	"net/http"
)

type Payload struct {
	RoomID  string `json:"room_id"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
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

	err := s.Service.SendMessage(&service.SendMessageParam{
		RoomID:  reqURI.RoomID,
		Sender:  c.Request.Header.Get("user"),
		Message: reqJSON.Message,
		After:   nil,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	payload := &Payload{
		RoomID:  reqURI.RoomID,
		Sender:  c.Request.Header.Get("user"),
		Message: reqJSON.Message,
	}
	b, _ := json.Marshal(payload)
	err = s.Redis.Publish(c.Request.Context(), reqURI.RoomID, b).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
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
	count, err := s.Service.GetUnreadCount(reqURI.RoomID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = s.Service.IncrementUnreadMessage(reqURI.RoomID, userID, -count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	messages := s.Service.GetRecentMessages(reqURI.RoomID, 10)
	c.JSON(http.StatusOK, messages)
}
