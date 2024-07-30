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

		acceptRouter := friendRouter.Group("/accept")
		{
			acceptRouter.PUT("", s.friendAccept)
		}

		mutualRouter := r.Group("/mutual")
		{
			mutualRouter.GET("/:user_id", s.mutualFriends)
			mutualRouter.GET("/:user_id/count", s.mutualFriendCount)
		}

		requestRouter := friendRouter.Group("/request")
		{
			requestRouter.GET("", s.listFromRequests)
			requestRouter.POST("", s.friendRequest)
			requestRouter.GET("/count", s.fromRequestCount)
		}
	}

	userRouter := r.Group("/users")
	{
		userRouter.POST("", s.createUser)
		userRouter.GET("/:user_id", s.getUser)
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
