package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type GetRoomsByUserIDResponse struct {
	RoomID        string    `json:"room_id"`
	RecentMessage string    `json:"recent_message"`
	Time          time.Time `json:"time"`
	UnreadCount   string    `json:"unread_count"`
}

func (s *Server) GetRoomsByUserID(c *gin.Context) {
	var req struct {
		UserID string `uri:"user_id"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	rooms, err := s.Service.GetRoomsByUserID(req.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].Time.After(rooms[j].Time)
	})

	times := s.Service.GetAllRoomsReadMessageTime(req.UserID)
	counts, err := s.Service.GetRoomsUnreadMessageCount(times)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var resp []*GetRoomsByUserIDResponse
	for _, r := range rooms {
		for _, c := range counts {
			if r.RoomID == c.RoomID {
				resp = append(resp, &GetRoomsByUserIDResponse{
					RoomID:        r.RoomID,
					Time:          r.Time,
					RecentMessage: r.RecentMessage,
					UnreadCount:   strconv.Itoa(c.Count),
				})
			}
		}
	}

	c.JSON(http.StatusOK, resp)
}
