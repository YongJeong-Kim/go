package gapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/yongjeong-kim/go/goapigw/account/service"
	"github.com/yongjeong-kim/go/goapigw/account/token"
	accountv1 "github.com/yongjeong-kim/go/goapigw/accountpb/pb/account/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

type AccountServer struct {
	accountv1.UnimplementedAccountServiceServer
	servicer service.AccountServicer
	maker    token.TokenMaker
}

func NewAccountServer(maker token.TokenMaker, svc service.AccountServicer) *AccountServer {
	return &AccountServer{
		maker:    maker,
		servicer: svc,
	}
}

func (s *AccountServer) CreateAccount(ctx context.Context, req *accountv1.CreateAccountRequest) (*accountv1.CreateAccountResponse, error) {
	return &accountv1.CreateAccountResponse{
		Account: &accountv1.Account{
			AccountId: "1234",
		},
	}, nil
}

func (s *AccountServer) Login(ctx context.Context, req *accountv1.LoginRequest) (*accountv1.LoginResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()
	if username != "aaa" || password != "1234" {
		_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", strconv.Itoa(http.StatusBadRequest)))
		return nil, &APIError{
			Inner:      fmt.Errorf("invalid username or password"),
			StatusCode: http.StatusBadRequest,
			Message:    "invalid username or password",
		}
	}

	tk, err := s.servicer.Login(username, password, 5*time.Minute)
	if err != nil {
		return nil, err
	}

	return &accountv1.LoginResponse{
		Token: tk,
	}, nil
}

func (s *AccountServer) UploadImage(ctx context.Context, req *accountv1.UploadImageRequest) (*emptypb.Empty, error) {
	log.Println(req.GetImage().GetFilename())
	log.Println(req.GetImage().GetContent())
	log.Println(req.GetImage().GetContentType())
	return nil, nil
}

func (s *AccountServer) setupGRPCServer(ctx context.Context) *grpc.Server {
	grpcServer := grpc.NewServer(
	//grpc.UnaryInterceptor(auth.UnaryServerInterceptor(exampleAuthFunc)),
	)
	RegisterService(grpcServer, s)
	reflection.Register(grpcServer)

	return grpcServer
}
func (s *AccountServer) RunGRPCServer(ctx context.Context, group *errgroup.Group) {
	grpcServer := s.setupGRPCServer(ctx)
	listener, err := net.Listen("tcp", "0.0.0.0:19090")
	if err != nil {
		log.Fatal("create listener failed.", err)
	}

	group.Go(func() error {
		log.Println("start grpc server")
		err = grpcServer.Serve(listener)
		if err != nil {
			log.Println("serve failed.", err)
			return err
		}

		return nil
	})

	group.Go(func() error {
		<-ctx.Done()
		log.Println("graceful shutdown grpc server...")
		grpcServer.GracefulStop()
		log.Println("graceful shutdown grpc server")
		return nil
	})
}

func (s *AccountServer) setupGatewayServer(ctx context.Context) *runtime.ServeMux {
	marshaler := &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, marshaler)

	grpcMux := runtime.NewServeMux(
		GWMultipartForm(marshaler),
		jsonOption,
		runtime.WithForwardResponseOption(httpResponseModifier),
		runtime.WithIncomingHeaderMatcher(incomingHeaderMatcher),
		runtime.WithMetadata(metadataMatcher),
		runtime.WithOutgoingHeaderMatcher(headerMatcher),
		runtime.WithRoutingErrorHandler(routingErrorHandler),
		runtime.WithErrorHandler(errorHandler),
	)

	err := RegisterHandlerServer(ctx, grpcMux, s)
	if err != nil {
		log.Fatal("cannot register handler server", err)
	}

	return grpcMux
}

func (s *AccountServer) RunGatewayServer(ctx context.Context, group *errgroup.Group) {
	grpcMux := s.setupGatewayServer(ctx)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.GetRouter(grpcMux),
	}

	group.Go(func() error {
		log.Println("start gateway server")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println("listen and serve err: ", err)
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-ctx.Done()
		log.Println("graceful shutdown gateway server...")
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Println("graceful shutdown gateway server failed.", err)
			return err
		}
		log.Println("stopped gateway server")
		return nil
	})
}

func (s *AccountServer) GetRouter(wrapHandler http.Handler) *gin.Engine {
	r := gin.New()
	r.Group("/v1/account/*{grpc_gateway}").Any("", gin.WrapH(wrapHandler))

	return r
}
