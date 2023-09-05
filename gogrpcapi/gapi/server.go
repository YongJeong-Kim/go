package gapi

import (
	"context"
	"github.com/gin-gonic/gin"
	accountv1 "gogrpcapi/pb/account/v1"
	userv1 "gogrpcapi/pb/user/v1"
	"gogrpcapi/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

type Server struct {
	UnimplementedServer
	tokenMaker token.Maker
}

type UnimplementedServer struct {
	accountv1.UnimplementedAccountServiceServer
	userv1.UnimplementedSimpleServerServer
}

func NewServer(tokenMaker token.Maker) *Server {
	return &Server{
		tokenMaker: tokenMaker,
	}
}

func (server *Server) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	return &userv1.CreateUserResponse{
		UserId: "vvvv",
	}, nil
}

func (server *Server) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "qqq")
	}
	return &userv1.DeleteUserResponse{
		UserId: authPayload.UserID,
	}, nil
}

func (server *Server) CreateAccount(ctx context.Context, req *accountv1.CreateAccountRequest) (*accountv1.CreateAccountResponse, error) {
	return &accountv1.CreateAccountResponse{
		Account: &accountv1.Account{
			AccountId: "dvdvdv",
		},
	}, nil
}

func (server *Server) UploadUser(ctx context.Context, req *userv1.UploadUserRequest) (*userv1.UploadUserResponse, error) {
	log.Println(req)
	_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "201"))
	return &userv1.UploadUserResponse{
		Id: "2eeeee",
	}, nil
}

func (server *Server) GetRouter(wrapHandler http.Handler, tokenMaker token.Maker) *gin.Engine {
	r := gin.New()
	r.Group("/v1/*{grpc_gateway}").Any("", gin.WrapH(wrapHandler))
	r.GET("/another", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return r
}
