package main

import (
	"context"
	"flag"
	"gogrpcpubsub/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	address := flag.String("address", "0.0.0.0:8080", "server address")
	to := flag.String("to", "", "user id")
	from := flag.String("from", "", "user id")
	content := flag.String("content", "", "user content")
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("connection error", err)
	}

	client := pb.NewSubscribeServiceClient(conn)
	stream, err := client.SubscribeBidi(context.Background())
	if err != nil {
		log.Fatal("connection error", err)
	}

	waitc := make(chan struct{})
	req := &pb.SubscribeRequest{
		To:      *to,
		From:    *from,
		Event:   "sendToUser",
		Content: *content,
	}
	go func() {
		err = stream.Send(req)
		if err != nil {
			log.Fatal("send failed.", err)
		}
		//waitc <- struct{}{}
	}()
	<-waitc
}
