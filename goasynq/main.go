package main

import (
	"context"
	"errors"
	"github.com/YongJeong-Kim/go/goasynq/api"
	"github.com/YongJeong-Kim/go/goasynq/service"
	"github.com/YongJeong-Kim/go/goasynq/worker"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
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
	conn, err := GetConnection()
	if err != nil {
		log.Fatal("connect db failed: ", err)
	}

	service := service.NewService(conn)
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
		DB:   3,
	})
	taskDistributor := worker.NewTaskDistributor()
	//taskLog := worker.NewTaskLog()
	server := api.NewServer(service, asynqClient, taskDistributor)
	server.SetupRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: server.Router,
	}

	go worker.NewTaskServer(server.TaskDistributor)
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

func GetConnection() (*sqlx.DB, error) {
	return sqlx.Connect("mysql", mysqlSource)
}
