package service_test

import (
	"bufio"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"gogrpc/pb"
	"gogrpc/sample"
	"gogrpc/serializer"
	"gogrpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"io/ioutil"
	"net"
	"os"
	"testing"
)

func TestClientCreatePerson(t *testing.T) {
	t.Parallel()

	personServer, serverAddress := startTestPersonServer(t)
	personClient := newTestPersonClient(t, serverAddress)

	person := sample.NewPerson()
	expectedID := person.Id
	req := &pb.CreatePersonRequest{
		Person: person,
	}

	res, err := personClient.CreatePerson(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res.PersonId, expectedID)

	other, err := personServer.Store.Find(res.PersonId)
	require.NoError(t, err)
	require.NotNil(t, other)
	require.Equal(t, other.Id, res.PersonId)

	requireSamePerson(t, person, other)
}

func TestClientUploadImage(t *testing.T) {
	t.Parallel()

	_, address := startTestPersonServer(t)
	personClient := newTestPersonClient(t, address)

	testImageFolder := "../tmp"
	filename := "ggo2.jpg"
	imagePath := fmt.Sprintf("%s/%s", testImageFolder, filename)
	file, err := os.Open(imagePath)
	require.NoError(t, err)
	defer func(file *os.File) {
		err := file.Close()
		require.NoError(t, err)
	}(file)

	stream, err := personClient.UploadImage(context.Background())
	require.NoError(t, err)

	req := &pb.UploadImageRequest{
		ChunkData: nil,
	}
	err = stream.Send(req)
	require.NoError(t, err)

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	size := 0

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)

		size += n
		req = &pb.UploadImageRequest{
			ChunkData: buffer[:n],
		}
		err = stream.Send(req)
		require.NoError(t, err)
	}

	res, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.NotZero(t, res.GetSize())
	require.EqualValues(t, res.GetSize(), size)

	savedImagePath := fmt.Sprintf("%s/%s", testImageFolder, "tmp_" + filename)
	err = ioutil.WriteFile(savedImagePath, buffer, 0644)
	require.NoError(t, err)
	require.FileExists(t, savedImagePath)
	require.NoError(t, os.Remove(savedImagePath))
}

func startTestPersonServer(t *testing.T) (*service.PersonServer, string) {
	personServer := service.NewPersonServer(service.NewInMemoryPersonStore())

	grpcServer := grpc.NewServer()
	pb.RegisterPersonServiceServer(grpcServer, personServer)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go func() {
		err := grpcServer.Serve(listener)
		require.NoError(t, err)
	}()

	return personServer, listener.Addr().String()
}

func newTestPersonClient(t *testing.T, serverAddress string) pb.PersonServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	return pb.NewPersonServiceClient(conn)
}

func requireSamePerson(t *testing.T, p1 *pb.Person, p2 *pb.Person) {
	p1s, err := serializer.ProtobufToJSON(p1)
	require.NoError(t, err)

	p2s, err := serializer.ProtobufToJSON(p2)
	require.NoError(t, err)

	require.Equal(t, p1s, p2s)
}
