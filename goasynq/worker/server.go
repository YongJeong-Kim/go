package worker

import (
	"context"
	"github.com/YongJeong-Kim/go/goasynq/tasks"
	"github.com/hibiken/asynq"
	"log"
)

const (
	TaskUser = "task:user"
)

type TaskServer struct {
	Server      *asynq.Server
	Mux         *asynq.ServeMux
	Distributor Distributor
}

func NewTaskServer(distributor Distributor) *TaskServer {
	return &TaskServer{
		Distributor: distributor,
	}
}

func (t *TaskServer) SetupTaskServer() {
	//logger, _ := zap.NewProduction()
	//defer logger.Sync()

	//logger := NewTaskLog()
	//redis.SetLogger(logger)
	t.Server = asynq.NewServer(
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
			ErrorHandler: asynq.ErrorHandlerFunc(HandleErrorFunc),
			//Logger: logger,
			// See the godoc for other configuration options
		},
	)
}

func (t *TaskServer) RunTaskServer() {
	log.Println("run task server")
	if err := t.Server.Run(t.Mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func (t *TaskServer) SetupServeMux() {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskUser, t.Distributor.CreateUserTask)
	mux.Handle(tasks.TypeImageResize, tasks.NewImageProcessor())

	t.Mux = mux
}

func HandleErrorFunc(ctx context.Context, task *asynq.Task, err error) {
	reachedErr := MaxRetryReachedHandler(ctx)
	if reachedErr != nil {
		log.Println(reachedErr)
	}

	log.Println("asynq handle error: ", err)
}
