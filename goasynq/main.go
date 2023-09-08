package main

import (
	"context"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"goasynq/api"
	"goasynq/service"
	"goasynq/worker"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	mysqlSource = "root:1234@tcp(localhost:13306)/test?parseTime=true"
	redisAddr   = "127.0.0.1:16379"
)

func main() {
	conn, err := sqlx.Connect("mysql", mysqlSource)
	if err != nil {
		log.Fatal("connect db failed: ", err)
	}

	service := service.NewService(conn)
	asynqclient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	//taskLog := worker.NewTaskLog()
	server := api.NewServer(service, asynqclient)
	server.SetupRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: server.Router,
	}

	go worker.NewTaskServer()
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
