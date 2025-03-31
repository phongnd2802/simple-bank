package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/phongnd2802/simple-bank/db/sqlc"
	"github.com/phongnd2802/simple-bank/email"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefaut   = "default"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
	mailer email.EmailSender
}

// Shutdown implements TaskProcessor.
func (processor *RedisTaskProcessor) Shutdown() {
	panic("unimplemented")
}

// Start implements TaskProcessor.
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskVerifyEmail)

	return processor.server.Start(mux)
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer email.EmailSender) *RedisTaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefaut:   5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("process task failed")
			}),
			Logger: NewLogger(),
		},
	)

	return &RedisTaskProcessor{server: server, store: store, mailer: mailer}
}

var _ TaskProcessor = (*RedisTaskProcessor)(nil)
