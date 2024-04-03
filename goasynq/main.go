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
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	mysqlSource = "root:1234@tcp(localhost:13306)/test?parseTime=true"
	redisAddr   = "127.0.0.1:16379"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	conn, err := GetConnection()
	if err != nil {
		log.Fatal("connect db failed: ", err)
	}

	service := service.NewService(conn)
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
		// 0 only working
		//DB:        0,
	})
	taskDistributor := worker.NewTaskDistributor()
	//taskLog := worker.NewTaskLog()
	server := api.NewServer(service, asynqClient, taskDistributor)
	server.SetupRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: server.Router,
	}

	taskServer := worker.NewTaskServer(server.TaskDistributor)
	taskServer.SetupTaskServer()
	taskServer.SetupServeMux()

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)

	taskServer.RunTaskServer(ctx, group)

	group.Go(func() error {
		log.Println("run http server")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-ctx.Done()
		log.Println("shutdown http server")
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Println("shutdown http server failed", err)
			return err
		}

		return nil
	})

	err = group.Wait()
	if err != nil {
		log.Fatal("wait group failed.", err)
	}

	//go func() {
	//	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
	//		log.Fatalf("listen: %s\n", err)
	//	}
	//}()
	//
	//quit := make(chan os.Signal)
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//<-quit
	//log.Println("Shutdown Server ...")
	//
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//if err := srv.Shutdown(ctx); err != nil {
	//	log.Fatal("Server Shutdown:", err)
	//}
	//// catching ctx.Done(). timeout of 5 seconds.
	//select {
	//case <-ctx.Done():
	//	log.Println("timeout of 5 seconds.")
	//}
	//log.Println("Server exiting")
}

func GetConnection() (*sqlx.DB, error) {
	return sqlx.Connect("mysql", mysqlSource)
}
