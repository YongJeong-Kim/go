package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gogrpcapi/gapi"
	accountv1 "gogrpcapi/pb/account/v1"
	userv1 "gogrpcapi/pb/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
)

func RegisterService(grpcServer *grpc.Server, server *gapi.Server) {
	accountv1.RegisterAccountServiceServer(grpcServer, server)
	userv1.RegisterSimpleServerServer(grpcServer, server)
}

func RegisterHandlerServer(ctx context.Context, grpcMux *runtime.ServeMux, server *gapi.Server) error {
	register := func(errs ...error) error {
		for _, err := range errs {
			if err != nil {
				return fmt.Errorf("register handler server failed.", err)
			}
		}

		return nil
	}

	return register(
		userv1.RegisterSimpleServerHandlerServer(ctx, grpcMux, server),
		accountv1.RegisterAccountServiceHandlerServer(ctx, grpcMux, server))
}

func main() {
	server := gapi.NewServer()
	grpcServer := grpc.NewServer()

	RegisterService(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		log.Fatal("cannot create listener")
	}

	go RunGatewayServer()

	log.Printf("start grpc server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server")
	}
}

func RunGatewayServer() {
	server := gapi.NewServer()

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := RegisterHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal("cannot create listener", err)
	}

	handler := gapi.HTTPLogger(mux)
	log.Println("start http gateway server at", listener.Addr().String())

	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server")
	}
}
