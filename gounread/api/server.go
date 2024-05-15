package api

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"gounread/service"
	"log"
	"net/http"
	"sort"
)

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

func (s *Server) SetupSubscription(ctx context.Context, channels []string) {
	sub := s.Redis.Subscribe(ctx, channels...)
	defer sub.Close()
	log.Println("subscribe: ", channels)

	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			log.Fatal("receive message error. ", err)
		}

		// unmarshal msg
		log.Printf("receive message: %v\n", msg)
		var payload *Payload
		_ = json.Unmarshal([]byte(msg.Payload), &payload)

		// select users from room
		var exceptMe []string
		users := s.Service.GetUsersByRoomID(payload.RoomID)
		for _, u := range users {
			if payload.Sender != u {
				exceptMe = append(exceptMe, u)
			}
		}

		// all users add unread except sender
		for _, e := range exceptMe {
			err := s.Service.IncrementUnreadMessage(payload.RoomID, e, 1)
			if err != nil {
				log.Println("add unread error. ", err)
			}
		}
	}
}

type Server struct {
	Service service.Servicer
	Socket  websocket.Upgrader
	Redis   *redis.Client
	Router  *gin.Engine
}

func NewServer(svr service.Servicer, redis *redis.Client) *Server {
	return &Server{
		Service: svr,
		Redis:   redis,
	}
}

func (s *Server) SetupRouter() {
	r := gin.New()
	r.GET("/users/:user_id", s.GetRoomsByUserID)
	r.POST("/rooms/:room_id/send", s.SendMessage)
	r.PUT("/rooms/:room_id/read", s.ReadMessage)
	/*r.POST("/rooms/:room_id", func(c *gin.Context) {
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
		err := s.Service.SendMessage(&service.SendMessageParam{
			RoomID:  reqURI.RoomID,
			Sender:  userID,
			Message: reqMsg.Message,
			After: func() error {
				err := s.Service.SetRecentMessage(reqURI.RoomID, reqMsg.Message)
				if err != nil {
					return err
				}
			},
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.Status(http.StatusNoContent)
	})*/

	s.Router = r
}

func NewRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		//Addr: "redis:16379",
		Addr:     "localhost:16379",
		Password: "",
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("redis connect failed. ", err)
	}

	return client
}
