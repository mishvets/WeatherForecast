package db

import (
	"context"
	"fmt"
	"strings"
)

func (q *Queries) UpdateWeatherBatch(ctx context.Context, updates []UpdateWeatherParams) ([]WeatherDatum, error) {
	if len(updates) == 0 {
		return nil, nil
	}

	var args []interface{}
	var values []string

	for i, u := range updates {
		n := i*4 + 1
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d)", n, n+1, n+2, n+3))
		args = append(args, u.City, u.Temperature, u.Humidity, u.Description)
	}

	query := fmt.Sprintf(`
        UPDATE weather_data AS wd
        SET
            temperature = v.temperature,
            humidity = v.humidity,
            description = v.description,
            updated_at = NOW()
        FROM (
            VALUES %s
        ) AS v(city, temperature, humidity, description)
        WHERE wd.city = v.city
        RETURNING wd.id, wd.city, wd.temperature, wd.humidity, wd.description, wd.updated_at;
    `, strings.Join(values, ","))

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []WeatherDatum
	for rows.Next() {
		var wd WeatherDatum
		if err := rows.Scan(
			&wd.ID, &wd.City, &wd.Temperature,
			&wd.Humidity, &wd.Description, &wd.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, wd)
	}

	return result, nil
}
