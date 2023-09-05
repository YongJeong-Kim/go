package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gogrpcapi/gapi"
	accountv1 "gogrpcapi/pb/account/v1"
	userv1 "gogrpcapi/pb/user/v1"
	"gogrpcapi/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"net/http"
	"strconv"
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
	//authFn := func(ctx context.Context) (context.Context, error) {
	//	token, err := auth.AuthFromMD(ctx, "bearer")
	//	if err != nil {
	//		return nil, err
	//	}
	//	// TODO: This is example only, perform proper Oauth/OIDC verification!
	//	if token != "yolo" {
	//		return nil, status.Error(codes.Unauthenticated, "invalid auth token")
	//	}
	//	// NOTE: You can also pass the token in the context for further interceptors or gRPC service code.
	//	return ctx, nil
	//}
	//
	//allButHealthZ := func(ctx context.Context, callMeta interceptors.CallMeta) bool {
	//	return userv1.SimpleServer_ServiceDesc.ServiceName != callMeta.Service
	//	//return healthpb.Health_ServiceDesc.ServiceName != callMeta.Service
	//}
	//
	//tokenMaker, err := token.NewJWTMaker("qidfjfiekwhigfvcdsxzaqwer45t6y7u8ijhnbvfcdx")
	//if err != nil {
	//	log.Fatal("create token maker failed", err)
	//}
	//server := gapi.NewServer(tokenMaker)
	//grpcServer := grpc.NewServer(
	//	grpc.ChainUnaryInterceptor(
	//		selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFn), selector.MatchFunc(allButHealthZ)),
	//	),
	//)
	//
	//RegisterService(grpcServer, server)
	//reflection.Register(grpcServer)
	//
	//listener, err := net.Listen("tcp", "0.0.0.0:19090")
	//if err != nil {
	//	log.Fatal("cannot create listener", err)
	//}
	//
	//go RunGatewayServer()
	//
	//log.Printf("start grpc server at %s", listener.Addr().String())
	//err = grpcServer.Serve(listener)
	//if err != nil {
	//	log.Fatal("cannot start grpc server")
	//}

	go RunGatewayServer()
	runGRPCServer()
}

func parseToken(token string) (struct{}, error) {
	return struct{}{}, nil
}
func userClaimFromToken(struct{}) string {
	return "foobar"
}

var tokenInfoKey struct{}

func exampleAuthFunc(ctx context.Context) (context.Context, error) {
	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	ctx = logging.InjectFields(ctx, logging.Fields{"auth.sub", userClaimFromToken(tokenInfo)})

	// WARNING: In production define your own type to avoid context collisions.
	return context.WithValue(ctx, tokenInfoKey, tokenInfo), nil
}

func runGRPCServer() {
	tokenMaker, err := token.NewJWTMaker("vDeow21qdQdnO4hPf82xFCd183DbtUos8v4EgY910Uh")
	if err != nil {
		log.Fatal("create token maker failed", err)
	}
	server := gapi.NewServer(tokenMaker)
	grpcServer := grpc.NewServer(

		grpc.UnaryInterceptor(auth.UnaryServerInterceptor(exampleAuthFunc)),
	)
	RegisterService(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", "0.0.0.0:19090")
	if err != nil {
		log.Fatal("create listener failed.", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("serve failed.", err)
	}
}

func RunGatewayServer() {
	tokenMaker, err := token.NewJWTMaker("vDeow21qdQdnO4hPf82xFCd183DbtUos8v4EgY910Uh")
	if err != nil {
		log.Fatal("create token maker failed", err)
	}
	server := gapi.NewServer(tokenMaker)

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = RegisterHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server", err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: server.GetRouter(grpcMux, tokenMaker),
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}

	//mux := http.NewServeMux()
	//mux.Handle("/", grpcMux)
	//
	//listener, err := net.Listen("tcp", "0.0.0.0:8080")
	//if err != nil {
	//	log.Fatal("cannot create listener", err)
	//}
	//
	////handler := gapi.HTTPLogger(mux)
	////log.Println("start http gateway server at", listener.Addr().String())
	//
	//err = http.Serve(listener, mux)
	//if err != nil {
	//	log.Fatal("cannot start HTTP gateway server")
	//}
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
