package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	db "github.com/mishvets/WeatherForecast/db/sqlc"
)

const TaskSendNotifyEmails = "task:send_notify_emails"

type PayloadSendNotifyEmails struct {
	City      string           `json:"city"`
	Frequency db.FrequencyEnum `json:"frequency"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendNotifyEmails(
	ctx context.Context,
	payload *PayloadSendNotifyEmails,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskSendNotifyEmails, jsonPayload, opts...)
	taskInfo, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf(
		"DistributeTaskSendNotifyEmails: type  - %v, payload - %s, queue - %v, max_retry - %v",
		task.Type(), task.Payload(), taskInfo.Queue, taskInfo.MaxRetry,
	)
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendNotifyEmails(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendNotifyEmails
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	arg := db.GetEmailsForUpdateParams{
		City:      payload.City,
		Frequency: payload.Frequency,
	}
	emails, err := processor.store.GetEmailsForUpdate(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to get emails for city(%s): %w", payload.City, err)
	}
	weather, err := processor.store.GetWeather(ctx, payload.City)
	if err != nil {
		return fmt.Errorf("failed to get weather data for city(%s): %w", payload.City, err)
	}

	subject := "Update from WeatherForecast"
	content := fmt.Sprintf(
		`<h1>Hello, this is the weather data for %s</h1>
	<p>Temperature - %.1f&deg;C,<br>
	Humidity - %d%%,<br>
	Description - %s</p>`,
		payload.City,
		weather.Temperature,
		weather.Humidity,
		weather.Description,
	)
	err = processor.emailSender.SendEmail(subject, content, []string{}, emails)
	if err != nil {
		log.Printf("failed to send notify email: %v", err)
	}
	log.Printf(
		"ProcessTaskSendNotifyEmails: type - %v, payload - %s",
		task.Type(), task.Payload(),
	)
	return nil
}
