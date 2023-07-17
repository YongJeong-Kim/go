package gapi

import (
	"context"
	pb "gogrpcapi/pb/user/v1"
)

type Server struct {
	pb.UnimplementedSimpleServerServer
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return &pb.CreateUserResponse{
		UserId: "vvvv",
	}, nil
}

func (server *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	return &pb.DeleteUserResponse{
		UserId: "vdvd",
	}, nil
}
