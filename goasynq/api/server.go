package api

import (
	"github.com/YongJeong-Kim/go/goasynq/service"
	"github.com/YongJeong-Kim/go/goasynq/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

type Server struct {
	Service         service.Servicer
	AsynqClient     *asynq.Client
	Router          *gin.Engine
	TaskDistributor worker.Distributor
	//TaskLog     worker.TaskLogger
}

func NewServer(
	service service.Servicer,
	asynqClient *asynq.Client,
	taskDistributor worker.Distributor,
	// taskLog worker.TaskLogger,
) *Server {
	return &Server{
		Service:         service,
		AsynqClient:     asynqClient,
		TaskDistributor: taskDistributor,
		//TaskLog:     taskLog,
	}
}

func (server *Server) SetupRouter() {
	r := gin.New()
	r.POST("/users", server.CreateUser)

	server.Router = r
}
