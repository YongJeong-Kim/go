package main

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"sort"
)

func (s *Server) Connect(c *gin.Context) {

}

func (s *Server) GetRoomsByUserID(c *gin.Context) {
	var req struct {
		UserID string `uri:"user_id"`
	}
	if err := c.ShouldBindUri(&req); err != nil {

	}

	rooms := s.Service.GetRoomsByUserID(req.UserID)

	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].Time.After(rooms[j].Time)
	})
	c.JSON(http.StatusOK, gin.H{
		"rooms": rooms,
	})
}

type Server struct {
	Service Servicer
	Redis   *redis.Client
	Router  *gin.Engine
}

func NewServer(svr Servicer, redis *redis.Client) *Server {
	return &Server{
		Service: svr,
		Redis:   redis,
	}
}

func (s *Server) SetupRouter() {
	r := gin.New()
	r.POST("/users/:user_id/connect", s.Connect)
	r.GET("/users/:user_id", s.GetRoomsByUserID)

	r.POST("/rooms/:room_id", func(c *gin.Context) {
		var reqURI struct {
			RoomID string `uri:"room_id"`
		}
		if err := c.ShouldBindUri(&reqURI); err != nil {

		}

		var reqMsg struct {
			Message string `json:"message"`
		}
		if err := c.ShouldBindJSON(&reqMsg); err != nil {

		}

		userID := c.Request.Header.Get("user")
		err := s.Service.SendMessage(&SendMessageParam{
			RoomID:  reqURI.RoomID,
			Sender:  userID,
			Message: reqMsg.Message,
			After: func() error {
				err := s.Service.SetRecentMessage(reqURI.RoomID, reqMsg.Message)
				if err != nil {
					return err
				}

				err = s.Service.AddMessageUnread(reqURI.RoomID, userID)
				if err != nil {
					return err
				}
				return nil
			},
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.Status(http.StatusNoContent)
	})

	s.Router = r
}
