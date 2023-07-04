package main

import (
	"gogrpcapi/gapi"
	"gogrpcapi/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	server := gapi.NewServer()
	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start grpc server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server")
	}
}
