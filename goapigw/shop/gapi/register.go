package gapi

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	shopv1 "github.com/yongjeong-kim/go/goapigw/shoppb/pb/shop/v1"
	"google.golang.org/grpc"
)

func RegisterService(grpcServer *grpc.Server, server *ShopServer) {
	shopv1.RegisterShopServiceServer(grpcServer, server)
}

func RegisterHandlerServer(ctx context.Context, grpcMux *runtime.ServeMux, server *ShopServer) error {
	register := func(errs ...error) error {
		for _, err := range errs {
			if err != nil {
				return fmt.Errorf("register handler server failed.", err)
			}
		}

		return nil
	}

	return register(
		shopv1.RegisterShopServiceHandlerServer(ctx, grpcMux, server),
	)
}
