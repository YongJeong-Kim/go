package main

import (
	"context"
	"flag"
	"gogrpcpubsub/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

func main() {
	address := flag.String("address", "0.0.0.0:8080", "server address")
	id := flag.String("id", "", "user id")
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot connect server")
	}
	client := pb.NewSubscribeServiceClient(conn)
	if err != nil {
		log.Fatal("cannot create service client")
	}

	stream, err := client.SubscribeBidi(context.Background())
	if err != nil {
		log.Fatal("client error")
	}

	req := &pb.SubscribeRequest{
		Id:    *id,
		Event: "subscribe",
	}

	waitc := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Print("no more data")
				break
			}
			if err != nil {
				log.Print("recv error", err)
				close(waitc)
			}
			log.Print("recv :", res)
		}
	}()
	go func() {
		err := stream.Send(req)
		if err != nil {
			log.Print("send error ", err)
		}
	}()
	<-waitc
}
