package gapi

import (
	"context"
	accountv1 "gogrpcapi/pb/account/v1"
	userv1 "gogrpcapi/pb/user/v1"
	"gogrpcapi/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
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
	return &userv1.DeleteUserResponse{
		UserId: "vdvd",
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
