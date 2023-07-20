package gapi

import (
	"context"
	"github.com/stretchr/testify/require"
	userv1 "gogrpcapi/pb/user/v1"
	"log"
	"testing"
)

func TestAA(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//
	//server := NewServer()
	//grpcServer := grpc.NewServer()
	//
	//RegisterService(grpcServer, server)
	//reflection.Register(grpcServer)
	//
	//listener, err := net.Listen("tcp", "0.0.0.0:9090")
	//if err != nil {
	//	log.Fatal("cannot create listener")
	//}
	//
	//go RunGatewayServer()
	//
	//log.Printf("start grpc server at %s", listener.Addr().String())
	//
	//rec := httptest.NewRecorder()
	//req, err := http.NewRequest(http.MethodPost, "/v1/users", nil)
	//require.NoError(t, err)
	//
	//m.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(1)
	//
	//err = grpcServer.Serve(listener)
	//if err != nil {
	//	log.Fatal("cannot start grpc server")
	//}
	//m := accountv1mock.NewMockAccountServiceServer(ctrl)

	s := NewServer()

	res, err := s.CreateUser(context.Background(), &userv1.CreateUserRequest{
		User: &userv1.User{
			Name: "fff",
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, res.UserId)

	log.Println(res.UserId)
}
