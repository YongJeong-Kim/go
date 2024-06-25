package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gounread/util"
	"net/http"
	"time"
)

func (s *Server) CreateRoom(c *gin.Context) {
	var req struct {
		UserIDs []string `json:"user_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.Request.Header.Get("user")
	filtered := util.Filter(req.UserIDs, func(u string) bool {
		return userID == u
	})
	if len(filtered) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "must be exclude yourself in invite user list",
		})
		return
	}

	req.UserIDs = append(req.UserIDs, userID)
	err := s.Service.CreateRoom(req.UserIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

type ListRoomsByUserIDResponse struct {
	RoomID        string    `json:"room_id"`
	RecentMessage string    `json:"recent_message"`
	Time          time.Time `json:"time"`
	UnreadCount   string    `json:"unread_count"`
}

func (s *Server) ListRoomsByUserID(c *gin.Context) {
	var req struct {
		UserID string `uri:"user_id"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := uuid.Validate(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rooms, err := s.Service.ListRoomsByUserID(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, rooms)
}
