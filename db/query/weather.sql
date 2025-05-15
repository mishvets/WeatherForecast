-- name: CreateWeather :one
INSERT INTO weather_data (
  city,
  temperature,
  humidity,
  description
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetWeather :one
SELECT * FROM weather_data
WHERE city = $1 LIMIT 1;

-- name: GetWeatherForUpdate :one
SELECT * FROM weather_data
WHERE city = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: UpdateWeather :one
UPDATE weather_data
SET temperature = $2,
    humidity = $3,
    description = $4
WHERE city = $1
RETURNING *;

-- name: DeleteWeather :exec
DELETE FROM weather_data WHERE city = $1;