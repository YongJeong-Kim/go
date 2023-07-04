package gapi

import "gogrpcapi/pb"

type Server struct {
	pb.UnimplementedUserServiceServer
}

func NewServer() *Server {
	return &Server{}
}
