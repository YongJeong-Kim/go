package worker

import (
	"context"
	"github.com/YongJeong-Kim/go/goasynq/tasks"
	"github.com/hibiken/asynq"
	"log"
)

const (
	TaskTest = "task:test"
)

func NewTaskServer(t Distributor) {
	//logger, _ := zap.NewProduction()
	//defer logger.Sync()

	//logger := NewTaskLog()
	//redis.SetLogger(logger)
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
			ErrorHandler: asynq.ErrorHandlerFunc(HandleErrorFunc),
			//Logger: logger,
			// See the godoc for other configuration options
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskTest, t.CreateUserTask)
	mux.Handle(tasks.TypeImageResize, tasks.NewImageProcessor())
	// ...register other handlers...

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func HandleErrorFunc(ctx context.Context, task *asynq.Task, err error) {
	log.Fatal("asynq handle error: ", err)
}
