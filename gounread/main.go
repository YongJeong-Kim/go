package main

import (
	"context"
	"github.com/nats-io/nats.go"
	"golang.org/x/sync/errgroup"
	"gounread/api"
	"gounread/embedded"
	"gounread/repository"
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
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var repo repository.Repositorier = repository.NewRepository(api.NewSession())
	var svc service.Servicer = service.NewService(repo)
	server := api.NewServer(svc, nil)

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
		log.Println("start s1")
		s1.Start()
		log.Println("start s2")
		s2.Start()
		log.Println("start s3")
		s3.Start()

		nc, err := nats.Connect(
			strings.Join(embedded.Servers, ","),
			nats.ConnectHandler(func(conn *nats.Conn) {
				log.Println("client connection")
			}),
			nats.UserInfo("aaa", "1234"),
		)
		if err != nil {
			log.Fatal(err)
		}
		//server.Nats = nc
		var n embedded.Notifier = embedded.NewNotify(nc)
		server.Notify = n

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
		server.Notify.Drain()
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
