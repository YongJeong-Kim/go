package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
	"gorelationship/service"
	mockfriendsvc "gorelationship/service/mock/friend"
	mockusersvc "gorelationship/service/mock/user"
	"testing"
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

		mutualRouter := friendRouter.Group("/mutual")
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

func newTestServer(t *testing.T) (*Server, *mockfriendsvc.MockFriender, *mockusersvc.MockUserManager) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mf := mockfriendsvc.NewMockFriender(ctrl)
	mu := mockusersvc.NewMockUserManager(ctrl)

	svc := service.NewService(mf, mu)
	svr := NewServer(svc)
	svr.SetupRouter()

	return svr, mf, mu
}
