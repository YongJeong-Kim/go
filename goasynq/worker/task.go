package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"time"
)

type CreateUserTaskPayload struct {
	Name string
}

func MaxRetryReachedHandler(ctx context.Context) error {
	retried, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)
	log.Printf("retried: %d, max retry: %d\n", retried, maxRetry)

	if retried >= maxRetry {
		return fmt.Errorf("max retry reached. %d, %d", retried, maxRetry)
	}

	return nil
}

func (t *TaskDistributor) CreateUserTask(ctx context.Context, task *asynq.Task) error {
	//err := MaxRetryReachedHandler(ctx)
	//if err != nil {
	//	return err
	//}

	var p CreateUserTaskPayload
	err := json.Unmarshal(task.Payload(), &p)
	if err != nil {
		return err
	}

	log.Println("5 secs...")
	time.Sleep(5 * time.Second)
	log.Println("create user task done")

	return fmt.Errorf("create user task error occur: ", err)
}
