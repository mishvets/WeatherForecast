package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	// Id int64 `json:"id"`
	Email string `json:"email"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	taskInfo, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf(
		"DistributeTaskSendVerifyEmail: type  - %v, payload - %s, queue - %v, max_retry - %v",
		task.Type(), task.Payload(), taskInfo.Queue, taskInfo.MaxRetry,
	)
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	subscription, err := processor.store.GetSubscription(ctx, payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("subscription doesn't exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get subsctiption: %w", err)
	}

	log.Printf(
		"ProcessTaskSendVerifyEmail: type - %v, payload - %s, email - %v",
		task.Type(), task.Payload(), subscription.Email,
	)
	return nil
}
