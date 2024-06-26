package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gounread/service"
	"net/http"
	"time"
)

type Payload struct {
	RoomID  string    `json:"room_id"`
	Sender  string    `json:"sender"`
	Sent    time.Time `json:"sent"`
	Message string    `json:"message"`
	Unread  []string  `json:"unread"`
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
	now := time.Now().UTC()
	payload, err := s.Service.SendMessage(&service.SendMessageParam{
		RoomID:  reqURI.RoomID,
		Sender:  userID,
		Sent:    now,
		Message: reqJSON.Message,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	b, _ := json.Marshal(payload)
	err = s.Notify.Publish("room."+payload.RoomID, b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s%v", "room publish error. ", err),
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

	messages, err := s.Service.ReadMessage(reqURI.RoomID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(messages) == 0 {
		recentMessages, err := s.Service.GetRecentMessages(reqURI.RoomID, 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, recentMessages)
		return
	}

	c.JSON(http.StatusOK, messages)
}
