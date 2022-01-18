package main

import (
	"context"
	"flag"
	"gogrpcpubsub/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

func main() {
	address := flag.String("address", "0.0.0.0:8080", "server port")
	to := flag.String("to", "qqq", "input to user")
	from := flag.String("from", "www", "input from user")
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot connect to server", err)
	}

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	msgClient := pb.NewMsgServiceClient(conn)
	stream, err := msgClient.SendToUser(context.Background())
	md, _ := stream.Header()
	log.Print(md)
	//metadata.AppendToOutgoingContext(stream.Context(), "myidkey", "myidvalue")
	err = stream.Send(&pb.SendUserRequest{
		Msg: &pb.Msg{
			Id:      *from,
			To:      *to,
			From:    *from,
			Content: "hello, I'm " + *from,
		},
	})
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
			log.Println("err io EOF in goroutine")
			return
		}
		if err != nil {
			//waitResponse <- fmt.Errorf("cannot receive stream response: %v", err)
			log.Println("err not nil in goroutine")
			return
		}
		log.Printf("recv: %s, from: %s, id: %s", res.Content, res.From)
		time.Sleep(1 * time.Second)
	}
	//}()

	//select {}
}
