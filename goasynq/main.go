package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"goasynq/api"
	"goasynq/service"
	"log"
)

const redisAddr = "127.0.0.1:16379"

func main() {
	//srv := asynq.NewServer(
	//	asynq.RedisClientOpt{Addr: redisAddr},
	//	asynq.Config{
	//		// Specify how many concurrent workers to use
	//		Concurrency: 10,
	//		// Optionally specify multiple queues with different priority.
	//		Queues: map[string]int{
	//			"critical": 6,
	//			"default":  3,
	//			"low":      1,
	//		},
	//		// See the godoc for other configuration options
	//	},
	//)
	//
	//// mux maps a type to a handler
	//mux := asynq.NewServeMux()
	//mux.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)
	//mux.Handle(tasks.TypeImageResize, tasks.NewImageProcessor())
	//// ...register other handlers...
	//
	//if err := srv.Run(mux); err != nil {
	//	log.Fatalf("could not run server: %v", err)
	//}

	go api.NewTaskServer()

	conn, err := sqlx.Connect("mysql", "root:1234@tcp(localhost:13306)/test?parseTime=true")
	if err != nil {
		log.Fatal("connect db failed: ", err)
	}

	service := service.NewService(conn)
	asynqclient := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:16379"})
	server := api.NewServer(service, asynqclient)

	r := gin.New()
	r.GET("/test", server.CreateUser)
	err = r.Run(":8080")
	if err != nil {
		log.Fatal("run server failed:", err)
	}
}
