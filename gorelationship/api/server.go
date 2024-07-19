package api

import (
	"github.com/gin-gonic/gin"
	"gorelationship/service"
)

func (s *Server) SetupRouter() {
	r := gin.New()
	r.SetTrustedProxies(nil)

	friendRouter := r.Group("/friends")
	{
		friendRouter.GET("", s.listFriends)
		friendRouter.GET("/count", s.friendCount)
	}

	userRouter := r.Group("/users")
	{
		userRouter.POST("/users", s.createUser)
	}

	s.Router = r
}

type Server struct {
	Service *service.Service
	Router  *gin.Engine
}

func NewServer(svc *service.Service) *Server {
	return &Server{
		Service: svc,
	}
}
