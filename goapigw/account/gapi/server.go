package main

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
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"syscall"
)

type AccountServer struct {
	accountv1.UnimplementedAccountServiceServer
	servicer service.AccountServicer
}

func NewAccountServer(svc service.AccountServicer) *AccountServer {
	return &AccountServer{
		servicer: svc,
	}
}

func (s *AccountServer) CreateAccount(ctx context.Context, req *accountv1.CreateAccountRequest) (*accountv1.CreateAccountResponse, error) {
	return &accountv1.CreateAccountResponse{
		Account: &accountv1.Account{
			AccountId: "",
		},
	}, nil
}

func (s *AccountServer) RunGRPCServer(ctx context.Context, group *errgroup.Group) {
	var maker token.TokenMaker = token.NewPasetoMaker()
	accountService := service.NewAccountService(maker)

	server := NewAccountServer(accountService)
	grpcServer := grpc.NewServer(
	//grpc.UnaryInterceptor(auth.UnaryServerInterceptor(exampleAuthFunc)),
	)
	RegisterService(grpcServer, server)
	reflection.Register(grpcServer)

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

func (s *AccountServer) RunGatewayServer(ctx context.Context, group *errgroup.Group) {
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(
		jsonOption,
		runtime.WithForwardResponseOption(httpResponseModifier),
		runtime.WithIncomingHeaderMatcher(incomingHeaderMatcher),
		runtime.WithMetadata(metadataMatcher),
		runtime.WithOutgoingHeaderMatcher(headerMatcher),
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
			//creating a new HTTTPStatusError with a custom status, and passing error
			newError := runtime.HTTPStatusError{
				HTTPStatus: 400,
				Err:        err,
			}
			// using default handler to do the rest of heavy lifting of marshaling error and adding headers
			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, writer, request, &newError)
		}),
	)

	err := RegisterHandlerServer(ctx, grpcMux, s)
	if err != nil {
		log.Fatal("cannot register handler server", err)
	}

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

func httpResponseModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	// set http status code
	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		// delete the headers to not expose any grpc-metadata in http response
		delete(md.HeaderMD, "x-http-code")
		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
		w.WriteHeader(code)
	}

	return nil
}

func incomingHeaderMatcher(header string) (string, bool) {
	log.Println("header:", header)
	return "", true
}

func metadataMatcher(ctx context.Context, req *http.Request) metadata.MD {
	bearer := req.Header.Get("Authorization")
	log.Println("token:", bearer)
	return nil
}

func headerMatcher(header string) (string, bool) {
	log.Println("header:", header)
	return "", true
}

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func RegisterService(grpcServer *grpc.Server, server *AccountServer) {
	accountv1.RegisterAccountServiceServer(grpcServer, server)
}

func RegisterHandlerServer(ctx context.Context, grpcMux *runtime.ServeMux, server *AccountServer) error {
	register := func(errs ...error) error {
		for _, err := range errs {
			if err != nil {
				return fmt.Errorf("register handler server failed.", err)
			}
		}

		return nil
	}

	return register(
		accountv1.RegisterAccountServiceHandlerServer(ctx, grpcMux, server),
	)
}

func (s *AccountServer) GetRouter(wrapHandler http.Handler) *gin.Engine {
	r := gin.New()
	r.Group("/v1/*{grpc_gateway}").Any("", gin.WrapH(wrapHandler))
	r.GET("/another", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return r
}
