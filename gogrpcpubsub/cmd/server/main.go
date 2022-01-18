package main

import (
	"flag"
	"fmt"
	"gogrpcpubsub/pb"
	"gogrpcpubsub/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	port := flag.Int("port", 0, "server  port")
	flag.Parse()

	msgServer := service.NewMsgServer()
	grpcServer := grpc.NewServer()
	pb.RegisterMsgServiceServer(grpcServer, msgServer)

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
