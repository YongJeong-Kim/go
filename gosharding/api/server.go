package api

import (
	"gosharding/service"
)

type Server struct {
	Service *service.Service
}

func NewServer(svc *service.Service) *Server {
	return &Server{
		Service: svc,
	}
}
