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

func msg() {
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

func subs() {
	port := flag.Int("port", 0, "server  port")
	flag.Parse()

	subsServer := service.NewSubsServer()
	grpcServer := grpc.NewServer()
	pb.RegisterSubscribeServiceServer(grpcServer, subsServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server")
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server")
	}
	log.Print("start subs server")
}

func main() {
	subs()
}
