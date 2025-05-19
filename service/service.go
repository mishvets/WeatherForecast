package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mishvets/WeatherForecast/internal/errs"
	"github.com/mishvets/WeatherForecast/util"
)

type Service interface {
	GetWeatherForCity(ctx context.Context, city string) (GetWeatherForCityResult, error)
}

type ServiceWeather struct {
	url    string
	apiKey string
}

func NewServiceWeather(url string, apiKey string) Service {
	return &ServiceWeather{
		url:    url,
		apiKey: apiKey,
	}
}

type weatherServiceRes struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`

	Current struct {
		TempC     float32 `json:"temp_c"`
		Humidity  int8    `json:"humidity"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`

	Error *struct {
		Code    uint16 `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type GetWeatherForCityResult struct {
	Temperature float32 `json:"temperature"`
	Humidity    int8    `json:"humidity"`
	Description string  `json:"description"`
}

func (service *ServiceWeather) GetWeatherForCity(ctx context.Context, city string) (GetWeatherForCityResult, error) {
	var weather weatherServiceRes
	var result GetWeatherForCityResult

	serviceUrl := fmt.Sprintf("%s?key=%s&q=%s", service.url, service.apiKey, city)

	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	body, err := util.GetRequest(reqCtx, serviceUrl)
	if err != nil {
		log.Printf("Error making API call. URL: %s, Error: %v", serviceUrl, err)
		return result, err
	}

	err = json.Unmarshal(body, &weather)
	if err != nil {
		log.Printf("Error unmarshaling API response. Body: %s, Error: %v", body, err)
		return result, err
	}

	if weather.Error != nil {
		if weather.Error.Code == 1006 {
			result.Description = weather.Error.Message
			err = errs.CityNotFound
		} else {
			err = fmt.Errorf("%s", weather.Error.Message)
		}
		return result, err
	}

	result = createGetWeatherForCityResult(weather)
	return result, nil
}

func createGetWeatherForCityResult(response weatherServiceRes) GetWeatherForCityResult {
	return GetWeatherForCityResult{
		Temperature: response.Current.TempC,
		Humidity:    response.Current.Humidity,
		Description: response.Current.Condition.Text,
	}
}
