package main

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/redis/go-redis/v9"
	"github.com/scylladb/gocqlx/v2"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var interruptSignal = []os.Signal{syscall.SIGTERM, syscall.SIGINT, os.Interrupt}

func main() {
	var svc Servicer = NewService(NewSession(), NewRedisClient())
	server := NewServer(svc, NewRedisClient())

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignal...)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)

	server.SetupRouter()
	svr := &http.Server{
		Addr:    ":8080",
		Handler: server.Router,
	}
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

		err = server.Redis.ShutdownSave(ctx).Err()
		if err != nil {
			log.Println("redis client shutdown failed.", err)
			return nil
		}
		log.Println("redis client shutdown")
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

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatal("redis client ping failed.", err)
	}

	return client
}
