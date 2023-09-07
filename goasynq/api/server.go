package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"goasynq/service"
	"goasynq/tasks"
	"goasynq/worker"
	"log"
	"net/http"
)

type Server struct {
	Service     service.Servicer
	AsynqClient *asynq.Client
}

func NewServer(service service.Servicer, asynqClient *asynq.Client) *Server {
	return &Server{
		Service:     service,
		AsynqClient: asynqClient,
	}
}

const (
	TaskTest = "task:test"
)

func (server *Server) CreateUser(c *gin.Context) {
	name := "asdf"
	err := server.Service.CreateUser(&service.CreateUserParam{
		Name: name,
		After: func(name string) error {
			log.Println("create user after")

			payload, err := json.Marshal(worker.CreateUserTaskPayload{Name: name})
			if err != nil {
				log.Fatal(err)
			}

			task := asynq.NewTask(TaskTest, payload, nil)
			info, err := server.AsynqClient.EnqueueContext(c, task)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(info)
			return fmt.Errorf("afafaf")
		},
	})
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

func NewTaskServer() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:16379"},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskTest, worker.CreateUserTask)
	mux.Handle(tasks.TypeImageResize, tasks.NewImageProcessor())
	// ...register other handlers...

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
