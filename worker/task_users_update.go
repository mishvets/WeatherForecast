package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	db "github.com/mishvets/WeatherForecast/db/sqlc"
	"github.com/mishvets/WeatherForecast/internal/errs"
)

const (
	TaskNotifyUsers = "task:notify_users"
)

type PayloadNotifyUsers struct {
	Frequency db.FrequencyEnum `json:"frequency"`
}

func (processor *RedisTaskProcessor) ProcessTaskNotifyUsers(ctx context.Context, task *asynq.Task) error {
	var payload PayloadNotifyUsers
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	// TODO: select those updated more than 10 minutes ago
	cities, err := processor.store.GetCitiesForUpdate(ctx, payload.Frequency)
	if err != nil {
		return fmt.Errorf("failed to get cities: %w", err)
	}

	for _, city := range cities {
		weatherData, err := processor.weatherService.GetWeatherForCity(ctx, city)
		if err != nil && !errors.Is(err, errs.CityNotFound) {
			log.Printf("fail to get weather data for %s: %v", city, err)
			continue
		}

		arg := db.UpdateWeatherParams{
			City:        city,
			Temperature: weatherData.Temperature,
			Humidity:    int16(weatherData.Humidity),
			Description: weatherData.Description,
		}

		_, err = processor.store.UpdateWeather(ctx, arg)
		if err != nil {
			log.Printf("failed to update weather data for %s: %v", city, err)
			continue
		}
		taskPayload := &PayloadSendNotifyEmails{
			City:      city,
			Frequency: db.FrequencyEnum(payload.Frequency),
		}
		opts := []asynq.Option{
			asynq.MaxRetry(10),
			asynq.ProcessIn(5 * time.Second),
		}
		processor.distributor.DistributeTaskSendNotifyEmails(ctx, taskPayload, opts...)
	}

	log.Printf(
		"ProcessTaskNotifyUsers: type - %v, payload - %s",
		task.Type(), task.Payload(),
	)
	return nil
}
