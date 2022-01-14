package service

import (
	"context"
	"gogrpc/pb"
	"log"
)

type PersonServer struct {
}

func CreatePerson(ctx context.Context, req *pb.CreatePersonRequest) (*pb.CreatePersonResponse, error) {
	person := req.GetPerson()
	log.Printf("receive a create person: shirt brand: %s", person.Shirt.Brand)

	res := &pb.CreatePersonResponse{
		PersonId: "afa",
	}

	return res, nil
}