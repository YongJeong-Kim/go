package gapi

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	accountv1 "github.com/yongjeong-kim/go/goapigw/accountpb/pb/account/v1"
	shopv1 "github.com/yongjeong-kim/go/goapigw/shoppb/pb/shop/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"net/http"
	"strconv"
)

const shopAddr = "localhost:29090"

type ShopServer struct {
	shopv1.UnimplementedShopServiceServer
	AccountClient accountv1.AccountServiceClient
}

func NewShopServer(accountClient accountv1.AccountServiceClient) *ShopServer {
	return &ShopServer{
		AccountClient: accountClient,
	}
}

func (s *ShopServer) RunGRPCServer(ctx context.Context, group *errgroup.Group) {
	grpcServer := grpc.NewServer(
	//grpc.UnaryInterceptor(auth.UnaryServerInterceptor(exampleAuthFunc)),
	)
	RegisterService(grpcServer, s)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", shopAddr)
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

func (s *ShopServer) RunGatewayServer(ctx context.Context, group *errgroup.Group) {
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
		runtime.WithRoutingErrorHandler(routingHandlerError),
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

	/*conn, err := grpc.NewClient(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	defer conn.Close()

	shopv1.RegisterShopServiceHandler(ctx, grpcMux, conn)

	err = http.ListenAndServe("0.0.0.0:8081", grpcMux)
	if err != nil {

	}*/

	srv := &http.Server{
		Addr: ":8081",
		//Handler: s.GetRouter(grpcMux),
		Handler: grpcMux,
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
	return "", true
}

func metadataMatcher(ctx context.Context, req *http.Request) metadata.MD {
	bearer := req.Header.Get("Authorization")
	log.Println("token:", bearer)
	return nil
}

func headerMatcher(header string) (string, bool) {
	return "", true
}

func (s *ShopServer) GetRouter(wrapHandler http.Handler) *gin.Engine {
	r := gin.New()
	//r.Group("/v1/shop/*{grpc_gateway}").Any("", gin.WrapH(wrapHandler))
	r.Group("").Any("", gin.WrapH(wrapHandler))

	return r
}

func (s *ShopServer) CreateShop(ctx context.Context, req *shopv1.CreateShopRequest) (*shopv1.CreateShopResponse, error) {
	return &shopv1.CreateShopResponse{
		Shop: &shopv1.Shop{
			ShopId: "qweqwe",
		},
	}, nil
}

func routingHandlerError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, statusCode int) {

}
