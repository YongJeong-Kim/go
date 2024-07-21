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
		friendRouter.GET("/accept", s.friendAccept)
		friendRouter.GET("/count", s.friendCount)

		mutualRouter := r.Group("/mutual")
		{
			mutualRouter.GET("/:user_id", s.mutualFriends)
			mutualRouter.GET("/:user_id/count", s.mutualFriendCount)
		}

		requestRouter := friendRouter.Group("request")
		{
			requestRouter.GET("", s.listFriendRequests)
			requestRouter.POST("", s.friendRequest)
			requestRouter.GET("/count", s.friendRequestCount)
		}
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
