package worker

import (
	"context"

	"github.com/hibiken/asynq"
)


type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistrubutor(redistOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redistOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
