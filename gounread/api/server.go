package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/nats-io/nats.go"
	"github.com/scylladb/gocqlx/v2"
	"gounread/service"
	"log"
	"net/http"
	"time"
)

func (s *Server) ConnectClient(c *gin.Context) {
	userID := c.Request.Header.Get("user")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user",
		})
		return
	}

	rooms, err := s.Service.GetRoomsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
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

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})
	})
	roomRouter := r.Group("/rooms")
	{
		roomRouter.POST("/:room_id/send", s.SendMessage)
		roomRouter.PUT("/:room_id/read", s.ReadMessage)
		roomRouter.GET("/:room_id", s.GetRoomStatusInLobby)
		roomRouter.POST("", s.CreateRoom)
	}
	s.Router = r
}

func NewSession() gocqlx.Session {
	//hosts := []string{"localhost:19042", "localhost:29042", "localhost:39042"}
	//hosts := []string{"localhost:19042"}
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Keyspace = "keyspace_name_2"
	cluster.Consistency = gocql.LocalOne
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "scylla",
		Password: "1234",
	}
	cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{
		Min:        time.Second,
		Max:        10 * time.Second,
		NumRetries: 5,
	}
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	/*cluster.PoolConfig = gocql.PoolConfig{
		HostSelectionPolicy: gocql.HostPoolHostPolicy(hostpool.New([]string{"localhost:9042"})),
	}*/

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatal("create cluster error", err)
	}

	return session
}
