package main

import (
	"context"
	"github.com/yongjeong-kim/go/goapigw/account/gapi"
	"github.com/yongjeong-kim/go/goapigw/account/service"
	"github.com/yongjeong-kim/go/goapigw/account/token"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	/*	var maker token.TokenMaker = token.NewPasetoMaker()
		var servicer service.AccountServicer = service.NewAccountService(maker)
		accountServer := api.NewAccountServer(servicer)
		accountServer.SetupRouter()
		accountServer.Router.Run(":8080")*/

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)
	var maker token.TokenMaker = token.NewPasetoMaker()
	var servicer service.AccountServicer = service.NewAccountService(maker)
	server := gapi.NewAccountServer(maker, servicer)

	server.RunGatewayServer(ctx, group)
	server.RunGRPCServer(ctx, group)
	if err := group.Wait(); err != nil {
		log.Fatal("group wait failed", err)
	}
}
