package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/nats-io/nats.go"
	"github.com/scylladb/gocqlx/v2"
	"gounread/service"
	"log"
	"net/http"
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
	r.GET("/users/:user_id/rooms", s.GetRoomsByUserID)
	r.POST("/connect", s.ConnectClient)

	roomRouter := r.Group("/rooms")
	{
		roomRouter.POST("/:room_id/send", s.SendMessage)
		roomRouter.PUT("/:room_id/read", s.ReadMessage)
		roomRouter.GET("/:room_id", s.GetRoomStatusInLobby)
	}
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
