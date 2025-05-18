package worker

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	db "github.com/mishvets/WeatherForecast/db/sqlc"
	"github.com/mishvets/WeatherForecast/mailer"
	"github.com/mishvets/WeatherForecast/service"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	ProcessTaskGetWeatherData(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server         *asynq.Server
	store          db.Store
	emailSender    mailer.EmailSender
	weatherService service.Service
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, emailSender mailer.EmailSender, weatherService service.Service) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Printf(
					"Process task failed: type  - %v, payload - %s, error - %v",
					task.Type(), task.Payload(), err)
			}),
		},
	)

	return &RedisTaskProcessor{
		server:         server,
		store:          store,
		emailSender:    emailSender,
		weatherService: weatherService,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	mux.HandleFunc(TaskGetWeatherData, processor.ProcessTaskGetWeatherData)

	return processor.server.Start(mux)
}
