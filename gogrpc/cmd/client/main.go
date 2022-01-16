package main

import (
	"context"
	"flag"
	"gogrpc/pb"
	"gogrpc/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

func main() {
	serverAddress := flag.String("address", "", "server address")
	flag.Parse()
	log.Printf("dial server: %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot dial server", err)
	}

	personClient := pb.NewPersonServiceClient(conn)
	person := sample.NewPerson()
	req := &pb.CreatePersonRequest{
		Person: person,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := personClient.CreatePerson(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Print("already exists person")
		} else {
			log.Fatalf("cannot create person %s", err)
		}
		return
	}
	log.Print(res.PersonId)
}
