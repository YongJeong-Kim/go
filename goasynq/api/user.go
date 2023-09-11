package api

import (
	"encoding/json"
	"github.com/YongJeong-Kim/go/goasynq/service"
	"github.com/YongJeong-Kim/go/goasynq/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"log"
	"net/http"
	"time"
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

			opts := []asynq.Option{
				asynq.MaxRetry(3),
				asynq.ProcessIn(2 * time.Second),
				asynq.Queue("critical"),
			}

			task := asynq.NewTask(worker.TaskTest, payload, opts...)
			info, err := server.AsynqClient.EnqueueContext(c, task)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(info)
			return nil
		},
	})
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusCreated)
}
