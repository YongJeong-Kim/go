package main

import (
	"flag"
	"fmt"
	"gogrpc/pb"
	"gogrpc/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	port := flag.Int("port", 0, "server port")
	flag.Parse()

	personServer := service.NewPersonServer(service.NewInMemoryPersonStore())
	grpcServer := grpc.NewServer()
	pb.RegisterPersonServiceServer(grpcServer, personServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server")
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server")
	}
}