package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	db "github.com/mishvets/WeatherForecast/db/sqlc"
	"github.com/mishvets/WeatherForecast/internal/errs"
)

// TODO: rename to createWeatherData
const (
	TaskGetWeatherData = "task:get_weather_data"
)

type PayloadGetWeatherData struct {
	ID    int64     `json:"id"`
	Token uuid.UUID `json:"token"`
	Email string    `json:"email"`
	City  string    `json:"city"`
}

func (distributor *RedisTaskDistributor) DistributeTaskGetWeatherData(
	ctx context.Context,
	payload *PayloadGetWeatherData,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskGetWeatherData, jsonPayload, opts...)
	taskInfo, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf(
		"DistributeTaskGetWeatherData: type  - %v, payload - %s, queue - %v, max_retry - %v",
		task.Type(), task.Payload(), taskInfo.Queue, taskInfo.MaxRetry,
	)
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskGetWeatherData(ctx context.Context, task *asynq.Task) error {
	var payload PayloadGetWeatherData
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	weatherData, err := processor.weatherService.GetWeatherForCity(ctx, payload.City)
	if err != nil {
		if errors.Is(err, errs.CityNotFound) {
			// TODO: fix and send cancelation email
			// handleErrorCityNotFound(ctx, payload.Token, processor.emailSender, processor.store)
			// return fmt.Errorf("%s: %w", errs.CityNotFound, asynq.SkipRetry)
		} else {
			return fmt.Errorf("%w", err)
		}
	}

	// TODO: process create new weather record in separate func
	arg := db.CreateNewWeatherTxParams{
		CreateWeatherParams: db.CreateWeatherParams{
			City:        payload.City,
			Temperature: weatherData.Temperature,
			Humidity:    int32(weatherData.Humidity), // TODO: check int8
			Description: weatherData.Description,
		},
		ID: payload.ID,
	}
	err = processor.store.CreateNewWeatherTx(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("the requested subscription(%v) not exist anymore: %w", arg.ID, asynq.SkipRetry)
		}
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return fmt.Errorf("city(%s) already present: %w", payload.City, asynq.SkipRetry)
			}
		}
		return fmt.Errorf("failed to add weather data: %w", err)
	}
	log.Printf(
		"ProcessTaskGetWeatherData: type - %v, payload - %s",
		task.Type(), task.Payload(),
	)
	return nil
}

// func handleErrorCityNotFound(ctx context.Context, token uuid.UUID, emailSender mailer.EmailSender, store db.Store) error {
// 	// send cancelation email
// 	afterDelete := func(city string, email string) error {
// 		subject := "WeatherForecast subscription cancelled"
// 		content := fmt.Sprintf(
// 			`Unfortunately, the city "%s" wasn't present in our data, so your subscription has been cancelled.
// 		Please try again with a different city name.`,
// 			city,
// 		)
// 		to := []string{email}
// 		err := emailSender.SendEmail(subject, content, to)
// 		if err != nil {
// 			return fmt.Errorf("failed to send cancelation email: %w", err)
// 	}

// 	arg := db.DeleteSubscriptionTxParams{
// 		Token:     token,
// 		AfterConfirm: afterDelete,
// 	}

// 	_, err = store.DeleteSubscriptionTx(ctx, arg)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return fmt.Errorf("the requested subscription not exist anymore: %w", asynq.SkipRetry)
// 		}
// 		return fmt.Errorf("failed to delete subscription: %w", err)
// 	}
// }
