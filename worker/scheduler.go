package worker

import (
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
)

type Scheduler interface {
	Start() error
}

type RedisScheduler struct {
	scheduler *asynq.Scheduler
}

func NewRedisScheduler(redisOpt asynq.RedisConnOpt) Scheduler {
	scheduler := asynq.NewScheduler(
		redisOpt,
		nil,
	)

	return &RedisScheduler{
		scheduler: scheduler,
	}
}

func (userNotify *RedisScheduler) Start() error {
	registerHourlyUpdates(userNotify)
	registerDailyUpdates(userNotify)

	return userNotify.scheduler.Start()
}

func registerHourlyUpdates(userNotify *RedisScheduler) {
	payload, err := json.Marshal(PayloadNotifyUsers{
		Frequency: "hourly",
	})
	if err != nil {
		log.Fatalf("failed to marshal hourly task payload: %v", err)
	}
	// note: for testing
	_, err = userNotify.scheduler.Register("*/1 * * * *", asynq.NewTask(TaskNotifyUsers, payload))
	_, err = userNotify.scheduler.Register("0 * * * *", asynq.NewTask(TaskNotifyUsers, payload))
	if err != nil {
		log.Fatalf("scheduler task registration error: %v", err)
	}
}

func registerDailyUpdates(userNotify *RedisScheduler) {
	payload, err := json.Marshal(PayloadNotifyUsers{
		Frequency: "daily",
	})
	if err != nil {
		log.Fatalf("failed to marshal hourly task payload: %v", err)
	}
	// note: for testing
	_, err = userNotify.scheduler.Register("*/2 * * * *", asynq.NewTask(TaskNotifyUsers, payload))
	_, err = userNotify.scheduler.Register("0 8 * * *", asynq.NewTask(TaskNotifyUsers, payload))
	if err != nil {
		log.Fatalf("scheduler task registration error: %v", err)
	}
}

// TODO: update everyone once a day, and then additionally those who are hourly