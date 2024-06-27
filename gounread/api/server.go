package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"gounread/embedded"
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
	//Nats    *nats.Conn
	Notify  embedded.Notifier
	Service service.Servicer
	Router  *gin.Engine
}

func NewServer(svr service.Servicer, n embedded.Notifier) *Server {
	return &Server{
		Service: svr,
		Notify:  n,
	}
}

func (s *Server) SetupRouter() {
	r := gin.New()
	r.GET("/users/:user_id/rooms", s.ListRoomsByUserID)
	r.POST("/connect", s.ConnectClient)

	roomRouter := r.Group("/rooms")
	{
		roomRouter.POST("/:room_id/send", s.SendMessage)
		roomRouter.PUT("/:room_id/read", s.ReadMessage)
		roomRouter.POST("", s.CreateRoom)
	}
	s.Router = r
}

func NewSession() gocqlx.Session {
	log.Println("new session start")
	//hosts := []string{"172.29.0.2", "172.29.0.3", "172.29.0.4"}
	hosts := []string{"localhost:19042", "localhost:29042", "localhost:39042"}
	//hosts := []string{"localhost:9142", "localhost:9242", "localhost:9342"}
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = "chat"
	cluster.ProtoVersion = 4
	// Cassandra uses QUORUM when querying system auth for the default Cassandra user "cassandra"
	// Cassandra uses LOCAL_ONE when querying system_auth for other users
	cluster.Consistency = gocql.LocalQuorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "scylla",
		Password: "1234",
	}
	/*cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{
		Min:        time.Second,
		Max:        10 * time.Second,
		NumRetries: 5,
	}*/
	/*	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{
		NumRetries: 3,
	}*/
	// Client Timeout should always be higher than the request timeouts configured on the server,
	// e.g.) ConnectTimeout = 5s, Timeout = 10s (ConnectTimeout <= Timeout)
	// if couldn't recover the scylla by the timeout -> gocql: no hosts available in the pool
	cluster.Timeout = 15 * time.Second        // client response timeout. default 11s
	cluster.ConnectTimeout = 15 * time.Second // default 11s
	// if DisableInitialHostLookup true, fast startup. but not recommended
	/*	cluster.DisableInitialHostLookup = false
		cluster.PoolConfig = gocql.PoolConfig{
			//HostSelectionPolicy: gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy()),
			//HostSelectionPolicy: gocql.TokenAwareHostPolicy(gocql.DCAwareRoundRobinPolicy("datacenter1")),
			HostSelectionPolicy: gocql.TokenAwareHostPolicy(gocql.DCAwareRoundRobinPolicy("datacenter1"), gocql.ShuffleReplicas()),
			//HostSelectionPolicy: gocql.TokenAwareHostPolicy(gocql.DCAwareRoundRobinPolicy("datacenter1"), gocql.NonLocalReplicasFallback()),
			//HostSelectionPolicy: gocql.SingleHostReadyPolicy(gocql.RoundRobinHostPolicy()),
		}*/
	// when 2 nodes+
	// reconnection policy enabled, 300s+ startup
	// reconnection policy disabled, 36s startup
	// startup impact based on ReconnectionPolicy duration
	// cluster nodes create session has expensive calculate
	// long time == waiting for startup
	// if DisableInitialHostLookup true, fast startup
	//	cluster.ReconnectInterval = time.Second
	/*	cluster.ReconnectionPolicy = &gocql.ExponentialReconnectionPolicy{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     5 * time.Second,
		MaxRetries:      3,
	}*/
	/*	cluster.ReconnectionPolicy = &gocql.ConstantReconnectionPolicy{
		Interval:   5 * time.Second,
		MaxRetries: 3,
	}*/
	// don't need a pool. Create a global Session
	// See https://pkg.go.dev/github.com/gocql/gocql?utm_source=godoc#Session
	// It's safe for concurrent use by multiple goroutines and a typical usage scenario is to have one global session object to interact with the whole Cassandra cluster.
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatal("create cluster error", err)
	}
	log.Println("create session success")
	// Wait for schema agreement before proceeding
	if err := session.AwaitSchemaAgreement(context.Background()); err != nil {
		log.Fatal("Failed to reach schema agreement:", err)
	}

	return session
}
