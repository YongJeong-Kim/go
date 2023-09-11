package worker

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"log"
	"time"
)

type CreateUserTaskPayload struct {
	Name string
}

func (t *TaskDistributor) CreateUserTask(ctx context.Context, task *asynq.Task) error {
	var p CreateUserTaskPayload
	err := json.Unmarshal(task.Payload(), &p)
	if err != nil {
		return err
	}

	log.Println("3 sec...")
	time.Sleep(3 * time.Second)
	log.Println("create user task done")

	return nil
}
