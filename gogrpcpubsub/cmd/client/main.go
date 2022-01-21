package main

import (
	"context"
	"flag"
	"github.com/google/uuid"
	"gogrpcpubsub/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"time"
)

func msg() {
	address := flag.String("address", "0.0.0.0:8080", "server port")
	id := flag.String("id", "", "input id")
	to := flag.String("to", "", "input to user")
	from := flag.String("from", "", "input from user")
	content := flag.String("content", "", "input content")
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot connect to server", err)
	}

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	var md metadata.MD
	var ctx context.Context

	switch {
	// subscribe
	case *id != "" && *to == "" && *from == "":
		md = metadata.Pairs("id", *id)
	// send all users
	case *id == "" && *to == "" && *from != "":
		md = metadata.Pairs("from", *from)
	// send specific user
	case *id == "" && *to != "" && *from != "":
		md = metadata.Pairs("to", *to, "from", *from)
	}

	ctx = metadata.NewOutgoingContext(context.Background(), md)
	msgClient := pb.NewMsgServiceClient(conn)
	stream, err := msgClient.SendToUser(ctx)

	req := &pb.SendUserRequest{
		Msg: &pb.Msg{
			Id:      *id,
			To:      *to,
			From:    *from,
			Content: *content,
		},
	}

	err = stream.Send(req)
	if err != nil {
		err = stream.CloseSend()
		if err != nil {
			log.Fatalf("cannot close send: %v", err)
		}
		log.Fatal("send failed.")
	}

	//go func() {
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			//waitResponse <- nil
			log.Println("io EOF in goroutine")
			return
		}
		if err != nil {
			//waitResponse <- fmt.Errorf("cannot receive stream response: %v", err)
			log.Println("err not nil in goroutine", err)
			return
		}

		log.Printf("recv: %s", res.Content)
		time.Sleep(1 * time.Second)
	}
	//}()

	//select {}
}

func main() {
	//subscribe()
	subscribeBidi()
}

func subscribe() {
	address := flag.String("address", "0.0.0.0:8080", "server port")
	id := flag.String("id", "", "input id")
	to := flag.String("to", "", "input to user")
	from := flag.String("from", "", "input from user")
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot connect to server", err)
	}

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	var md metadata.MD
	var ctx context.Context

	switch {
	// subscribe
	case *id != "" && *to == "" && *from == "":
		md = metadata.Pairs("id", *id)
	// send all users
	case *id == "" && *to == "" && *from != "":
		md = metadata.Pairs("from", *from)
	// send specific user
	case *id == "" && *to != "" && *from != "":
		md = metadata.Pairs("to", *to, "from", *from)
	}

	ctx = metadata.NewOutgoingContext(context.Background(), md)
	subscribeClient := pb.NewSubscribeServiceClient(conn)
	req := &pb.SubscribeRequest{
		Id: "test subs id",
	}
	stream, err := subscribeClient.Subscribe(ctx, req)
	if err != nil {
		log.Print("send err not nil")
	}

	waitc := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Print("no more data", err)
				close(waitc)
				return
			}
			if err != nil {
				log.Print("err is not nil", err)
				return
			}
			log.Printf("response id: %s", res.Id)
		}
	}()
	err = stream.CloseSend()
	if err != nil {
		log.Print(err)
		return
	}
	<-waitc
}

func subscribeBidi() {
	address := flag.String("address", "0.0.0.0:8080", "server address")
	_ = flag.String("to", "", "input to user")
	_ = flag.String("from", "", "input from user")
	flag.Parse()

	md := metadata.Pairs("aa", "bb")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("connection failed")
	}
	client := pb.NewSubscribeServiceClient(conn)
	stream, err := client.SubscribeBidi(ctx)
	if err != nil {
		log.Print("client subscribe bidi error", err)
	}

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Print("no more data from server")
				break
			}
			if err != nil {
				log.Print("err is not nil", err)
			}
			log.Printf("recv client id: %s", res.Id)
		}
	}()

	id := uuid.NewString()
	for {
		req := &pb.SubscribeRequest{
			Id: id,
		}
		err := stream.Send(req)
		if err != nil {
			log.Print("send error")
		}
		time.Sleep(1 * time.Second)
	}
}
