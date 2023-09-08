package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"goasynq/service"
	"goasynq/worker"
	"log"
	"net/http"
)

type Server struct {
	Service     service.Servicer
	AsynqClient *asynq.Client
	Router      *gin.Engine
	//TaskLog     worker.TaskLogger
}

func NewServer(
	service service.Servicer,
	asynqClient *asynq.Client,
	// taskLog worker.TaskLogger,
) *Server {
	return &Server{
		Service:     service,
		AsynqClient: asynqClient,
		//TaskLog:     taskLog,
	}
}

func (server *Server) SetupRouter() {
	r := gin.New()
	r.GET("/test", server.CreateUser)

	server.Router = r
}

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

			task := asynq.NewTask(worker.TaskTest, payload, nil)
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
