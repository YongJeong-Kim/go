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
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("connection failed", err)
	}

	client := pb.NewSubscribeServiceClient(conn)
	stream, err := client.SubscribeBidi(context.Background())
	if err != nil {
		log.Fatal("connection failed. ", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Print("EOF")
			}
			if err != nil {
				log.Fatal("response error. ", err)
			}
			log.Print("unsubscribe successful. ", res)
			close(waitc)
		}
	}()

	go func() {
		req := &pb.SubscribeRequest{
			Id:    "ccc",
			Event: "unsubscribe",
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("client send error. ", err)
		}
	}()
	<-waitc
}
