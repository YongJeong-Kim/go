package main

import (
	"context"
	"github.com/nats-io/nats.go"
	"golang.org/x/sync/errgroup"
	"gounread/api"
	"gounread/embedded"
	"gounread/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var interruptSignal = []os.Signal{syscall.SIGTERM, syscall.SIGINT, os.Interrupt}

func main() {
	var svc service.Servicer = service.NewService(api.NewSession())
	server := api.NewServer(svc)

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignal...)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)

	server.SetupRouter()
	svr := &http.Server{
		Addr:    ":8080",
		Handler: server.Router,
	}

	s1 := embedded.NewServer(embedded.Opts1)
	s2 := embedded.NewServer(embedded.Opts2)
	s3 := embedded.NewServer(embedded.Opts3)

	group.Go(func() error {
		/*	server.SetupSubscription(ctx, []string{
			"cd8bb8a2-f947-4777-92b9-fbdb839c67ac",
			"01f84cfa-e487-494c-82e5-e75f95ef0573",
		})*/
		log.Println("start s1")
		s1.Start()
		log.Println("start s2")
		s2.Start()
		log.Println("start s3")
		s3.Start()

		ns, err := nats.Connect(
			strings.Join(embedded.Servers, ","),
			nats.ConnectHandler(func(conn *nats.Conn) {
				log.Println("client connection")
			}),
			nats.UserInfo("aaa", "1234"),
		)
		if err != nil {
			log.Fatal(err)
		}
		server.Nats = ns

		return nil
	})
	group.Go(func() error {
		log.Println("start http server")
		if err := svr.ListenAndServe(); err != nil {
			log.Fatal("listen failed.", err)
		}
		return nil
	})
	group.Go(func() error {
		<-ctx.Done()
		err := svr.Shutdown(ctx)
		if err != nil {
			log.Println("http server shutdown failed.", err)
			return nil
		}

		log.Println("drain connection")
		server.Nats.Drain()
		log.Println("s1 wait for shutdown")
		s1.WaitForShutdown()
		log.Println("s2 wait for shutdown")
		s2.WaitForShutdown()
		log.Println("s3 wait for shutdown")
		s3.WaitForShutdown()

		return nil
	})

	err := group.Wait()
	if err != nil {
		log.Fatal("wait group failed", err)
	}
}