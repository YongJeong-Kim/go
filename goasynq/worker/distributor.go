package worker

import (
	"context"
	"github.com/hibiken/asynq"
)

type Distributor interface {
	CreateUserTask(ctx context.Context, t *asynq.Task) error
}

type TaskDistributor struct {
}

func NewTaskDistributor() Distributor {
	return &TaskDistributor{}
}
