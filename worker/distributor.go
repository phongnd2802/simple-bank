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


func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) *RedisTaskDistributor {
	client := asynq.NewClient(redisOpt)

	return &RedisTaskDistributor{client: client}
}

var _ TaskDistributor = (*RedisTaskDistributor)(nil)
