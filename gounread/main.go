package main

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"golang.org/x/sync/errgroup"
	"gounread/api"
	"gounread/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var interruptSignal = []os.Signal{syscall.SIGTERM, syscall.SIGINT, os.Interrupt}

func main() {
	var svc service.Servicer = service.NewService(NewSession())
	server := api.NewServer(svc, api.NewRedis())

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignal...)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)

	server.SetupRouter()
	svr := &http.Server{
		Addr:    ":8080",
		Handler: server.Router,
	}
	group.Go(func() error {
		server.SetupSubscription(ctx, []string{
			"cd8bb8a2-f947-4777-92b9-fbdb839c67ac",
			"01f84cfa-e487-494c-82e5-e75f95ef0573",
		})
		return nil
	})
	group.Go(func() error {
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

		return nil
	})

	err := group.Wait()
	if err != nil {
		log.Fatal("wait group failed", err)
	}
}

func NewSession() gocqlx.Session {
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Keyspace = "keyspace_name"
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "scylla",
		Password: "1234",
	}

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatal("create cluster error", err)
	}

	return session
}
