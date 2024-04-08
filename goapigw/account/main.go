package main

import (
	"context"
	"golang.org/x/sync/errgroup"
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
}
