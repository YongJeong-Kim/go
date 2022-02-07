package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/google/uuid"
	"gogrpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
)

const maxImageSize = 1 << 50

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

func (server *PersonServer) SearchPerson(
	req *pb.SearchPersonRequest,
	stream pb.PersonService_SearchPersonServer,
) error {
	filter := req.GetFilter()
	log.Printf("receive filter: %v", filter)

	err := server.Store.Search(
		stream.Context(),
		filter,
		func(person *pb.Person) error {
			res := &pb.SearchPersonResponse{
				Person: &pb.Person{
					Id: "",
				},
			}
			err := stream.Send(res)
			if err != nil {
				return err
			}
			log.Printf("sent person id: %s", person.GetId())
			return nil
		},
	)

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil
}

func (server *PersonServer) UploadImage(stream pb.PersonService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		log.Print("recv error", err)
	}
	log.Print("req from client: ", req)

	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return logError(err)
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		imageSize += size
		if imageSize > maxImageSize {
			return logError(status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, maxImageSize))
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
		}
	}
	//imageID, err := server.imageStore.Save(laptopID, imageType, imageData)
	//if err != nil {
	//	return logError(status.Errorf(codes.Internal, "cannot save image to the store: %v", err))
	//}

	// save file
	//savePath := "C:\\Users\\admin\\Desktop\\Docker Desktop Installer.exe"
	//err = ioutil.WriteFile(savePath, imageData.Bytes(), 0644)
	//if err != nil {
	//	return logError(status.Errorf(codes.Internal, "cannot save image file %v", err))
	//}

	res := &pb.UploadImageResponse{
		//Id:   imageID,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}

	log.Printf("saved image with size: %d", imageSize)
	return nil
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
