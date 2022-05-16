package main

import (
	"context"
	"flag"
	"gogrpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

func main() {
	serverAddress := flag.String("address", "", "server address")
	flag.Parse()
	log.Printf("dial server: %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot dial server", err)
	}

	shirtClient := pb.NewShirtServiceClient(conn)
	stream, err := shirtClient.Broadcast(context.Background())
	if err != nil {
		return
	}
	req := &pb.ShirtRequest{
		Shirt: &pb.Shirt{
			//Brand: "Subscribe",
			Brand: "Publish",
			//Brand: "Only",
		},
	}

	go func() {
		err = stream.Send(req)
	}()

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Println("end")
			break
		}
		if err != nil {
			log.Fatal("fatal err ", err)
		}
		log.Println(res.GetShirt())
	}
}
