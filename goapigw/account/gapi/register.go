package gapi

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	accountv1 "github.com/yongjeong-kim/go/goapigw/accountpb/pb/account/v1"
	"google.golang.org/grpc"
)

func RegisterService(grpcServer *grpc.Server, server *AccountServer) {
	accountv1.RegisterAccountServiceServer(grpcServer, server)
}

func RegisterHandlerServer(ctx context.Context, grpcMux *runtime.ServeMux, server *AccountServer) error {
	register := func(errs ...error) error {
		for _, err := range errs {
			if err != nil {
				return fmt.Errorf("register handler server failed. %w", err)
			}
		}

		return nil
	}

	return register(
		accountv1.RegisterAccountServiceHandlerServer(ctx, grpcMux, server),
	)
}
