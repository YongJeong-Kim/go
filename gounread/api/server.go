package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/nats-io/nats.go"
	"github.com/scylladb/gocqlx/v2"
	"gounread/service"
	"log"
	"net/http"
	"sort"
)

func (s *Server) ConnectClient(c *gin.Context) {
	userID := c.Request.Header.Get("user")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user",
		})
		return
	}

	rooms := s.Service.GetRoomsByUserID(userID)
	c.JSON(http.StatusOK, rooms)
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
	Nats    *nats.Conn
	Service service.Servicer
	Router  *gin.Engine
}

func NewServer(svr service.Servicer) *Server {
	return &Server{
		Service: svr,
	}
}

func (s *Server) SetupRouter() {
	r := gin.New()
	r.GET("/users/:user_id", s.GetRoomsByUserID)
	r.POST("/rooms/:room_id/send", s.SendMessage)
	r.PUT("/rooms/:room_id/read", s.ReadMessage)
	r.POST("/connect", s.ConnectClient)
	r.GET("/unread/:room_id", s.GetUnreadCount)
	s.Router = r
}

func NewSession() gocqlx.Session {
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Keyspace = "keyspace_name"
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "scylla",
		Password: "1234",
	}

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatal("create cluster error", err)
	}

	return session
}
