package api

import (
	"log"
	"slices"

	"github.com/go-playground/validator/v10"
	db "github.com/mishvets/WeatherForecast/db/sqlc"
)

var validFrequency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if frequency, ok := fieldLevel.Field().Interface().(db.FrequencyEnum); ok {
		log.Println("###", slices.Contains(db.AllFrequencyEnumValues(), frequency)) // TODO: delete
		return slices.Contains(db.AllFrequencyEnumValues(), frequency)
	}
	return false
}
