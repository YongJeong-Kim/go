package service_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"gogrpc/pb"
	"gogrpc/sample"
	"gogrpc/serializer"
	"gogrpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
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
