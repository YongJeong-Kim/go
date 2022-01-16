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
	"io"
	"log"
	"time"
)

func createPerson(personClient pb.PersonServiceClient) {
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
	log.Print("create person:", res.PersonId)
}

func searchPerson(personClient pb.PersonServiceClient, filter *pb.Filter) {
	log.Print("search person:", filter)

	req := &pb.SearchPersonRequest{
		Filter: filter,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := personClient.SearchPerson(ctx, req)
	if err != nil {
		log.Fatalf("cannot search person", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("cannot receive response:", err)
		}

		p := res.GetPerson()
		log.Print("found:", p)
	}
}

func main() {
	serverAddress := flag.String("address", "", "server address")
	flag.Parse()
	log.Printf("dial server: %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot dial server", err)
	}

	personClient := pb.NewPersonServiceClient(conn)
	//createPerson(personClient)

	for i := 0; i < 10; i++ {
		createPerson(personClient)
	}
	filter := &pb.Filter{
		Brand: "LACOSTE",
	}
	searchPerson(personClient, filter)
}
