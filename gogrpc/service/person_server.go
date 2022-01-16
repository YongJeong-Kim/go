package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gogrpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type PersonServer struct {
	pb.UnimplementedPersonServiceServer
	Store PersonStore
}

func NewPersonServer(store PersonStore) *PersonServer {
	return &PersonServer{
		Store: store,
	}
}

func (server *PersonServer) CreatePerson(ctx context.Context, req *pb.CreatePersonRequest) (*pb.CreatePersonResponse, error) {
	//time.Sleep(6 * time.Second)

	if err := contextError(ctx); err != nil {
		return nil, err
	}

	person := req.GetPerson()
	log.Printf("receive a create person id: %s", person.Id)

	if len(person.Id) > 0 {
		_, err := uuid.Parse(person.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "person id invalid. %v", err)
		}
	} else {
		uuid, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "generate uuid failed.")
		}
		person.Id = uuid.String()
	}

	err := server.Store.Save(person)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "person save failed. %v", err)
	}

	res := &pb.CreatePersonResponse{
		PersonId: person.Id,
	}

	return res, nil
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is canceled."))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded."))
	default:
		return nil
	}
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
