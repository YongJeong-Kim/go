package main

import (
	"context"
	accountv1 "github.com/yongjeong-kim/go/goapigw/accountpb/pb/account/v1"
	"github.com/yongjeong-kim/go/goapigw/shop/gapi"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const accountAddr = "localhost:19090"

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(accountAddr, opts...)
	if err != nil {
		log.Fatal("new client error", err)
	}

	accountClient := accountv1.NewAccountServiceClient(conn)

	server := gapi.NewShopServer(accountClient)

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)

	server.RunGRPCServer(ctx, group)
	server.RunGatewayServer(ctx, group)
	if err := group.Wait(); err != nil {
		log.Fatal("errgroup wait error", err)
	}
}
