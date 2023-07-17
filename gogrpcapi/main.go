package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gogrpcapi/gapi"
	userv1 "gogrpcapi/pb/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
)

func main() {
	server := gapi.NewServer()
	grpcServer := grpc.NewServer()

	userv1.RegisterSimpleServerServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		log.Fatal("cannot create listener")
	}

	go runGatewayServer()

	log.Printf("start grpc server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server")
	}
}

func runGatewayServer() {
	server := gapi.NewServer()
	//if err != nil {
	//	log.Fatal().Err(err).Msg("cannot create server")
	//}

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

	err := userv1.RegisterSimpleServerHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server")
		//log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", grpcMux)

	//statikFS, err := fs.New()
	//if err != nil {
	//	log.Fatal().Err(err).Msg("cannot create statik fs")
	//}

	//swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	//mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal("cannot create listener", err)
		//log.Fatal().Err(err).Msg("cannot create listener")
	}

	//log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())
	handler := gapi.HTTPLogger(mux)
	log.Println("start http gateway server at", listener.Addr().String())

	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server")
		//log.Fatal().Err(err).Msg("cannot start HTTP gateway server")
	}
}
