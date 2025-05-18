package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	db "github.com/mishvets/WeatherForecast/db/sqlc"
)

const (
	TaskNotifyUsers = "task:notify_users"
	MaxCitiesPerRequest = 50
)

type PayloadNotifyUsers struct {
	Frequency db.FrequencyEnum `json:"frequency"`
}

func (processor *RedisTaskProcessor) ProcessTaskNotifyUsers(ctx context.Context, task *asynq.Task) error {
	var payload PayloadNotifyUsers
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	// cities, err := processor.store.GetCitiesForUpdate(ctx, payload.Frequency)
	// if err != nil {
	// 	return fmt.Errorf("failed to get cities: %w", err)
	// }

	// for i := 0; i < len(cities); i += MaxCitiesPerRequest {
	// 	end := i + MaxCitiesPerRequest
	// 	if end > len(cities) {
	// 		end = len(cities)
	// 	}
	// 	weatherData, err := processor.weatherService.GetWeatherForCities(ctx, cities[i:end])
	// 	if err != nil {
	// 		return fmt.Errorf("fail to get weather slice: %w", err)
	// 	}

	// 	// TODO: process create new weather record in separate func
	// 	arg := db.CreateNewWeatherTxParams{
	// 		CreateWeatherParams: db.CreateWeatherParams{
	// 			City:        payload.City,
	// 			Temperature: weatherData.Temperature,
	// 			Humidity:    int32(weatherData.Humidity), // TODO: check int8
	// 			Description: weatherData.Description,
	// 		},
	// 		ID: payload.ID,
	// 	}
	// 	err = processor.store.CreateNewWeatherTx(ctx, arg)
	// 	if err != nil {
	// 		if errors.Is(err, sql.ErrNoRows) {
	// 			return fmt.Errorf("the requested subscription(%v) not exist anymore: %w", arg.ID, asynq.SkipRetry)
	// 		}
	// 		if pqErr, ok := err.(*pq.Error); ok {
	// 			switch pqErr.Code.Name() {
	// 			case "unique_violation":
	// 				return fmt.Errorf("city(%s) already present: %w", payload.City, asynq.SkipRetry)
	// 			}
	// 		}
	// 		return fmt.Errorf("failed to add weather data: %w", err)
	// 	}
	// }

	log.Printf(
		"ProcessTaskNotifyUsers: type - %v, payload - %s",
		task.Type(), task.Payload(),
	)
	return nil
}
