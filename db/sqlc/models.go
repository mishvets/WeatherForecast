// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type FrequencyEnum string

const (
	FrequencyEnumDaily  FrequencyEnum = "daily"
	FrequencyEnumHourly FrequencyEnum = "hourly"
)

func (e *FrequencyEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = FrequencyEnum(s)
	case string:
		*e = FrequencyEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for FrequencyEnum: %T", src)
	}
	return nil
}

type NullFrequencyEnum struct {
	FrequencyEnum FrequencyEnum `json:"frequency_enum"`
	Valid         bool          `json:"valid"` // Valid is true if FrequencyEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullFrequencyEnum) Scan(value interface{}) error {
	if value == nil {
		ns.FrequencyEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.FrequencyEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullFrequencyEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.FrequencyEnum), nil
}

func AllFrequencyEnumValues() []FrequencyEnum {
	return []FrequencyEnum{
		FrequencyEnumDaily,
		FrequencyEnumHourly,
	}
}

type Subscription struct {
	ID        int64         `json:"id"`
	Email     string        `json:"email"`
	City      string        `json:"city"`
	Frequency FrequencyEnum `json:"frequency"`
	Confirmed bool          `json:"confirmed"`
	Token     uuid.UUID     `json:"token"`
	CreatedAt time.Time     `json:"created_at"`
}

type WeatherDatum struct {
	ID          int64     `json:"id"`
	City        string    `json:"city"`
	Temperature float32   `json:"temperature"`
	Humidity    int16     `json:"humidity"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}
