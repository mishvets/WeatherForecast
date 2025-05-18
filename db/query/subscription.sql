-- name: CreateSubscription :one
INSERT INTO subscriptions (
  email,
  city,
  frequency
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetSubscription :one
SELECT * FROM subscriptions
WHERE email = $1 LIMIT 1;

-- name: GetSubscriptionForUpdate :one
SELECT * FROM subscriptions
WHERE token = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ConfirmSubscription :one
UPDATE subscriptions
SET confirmed = $2
WHERE token = $1
RETURNING *;

-- name: DeleteSubscription :one
DELETE FROM subscriptions WHERE token = $1 RETURNING token;

-- name: IsSubscriptionExist :one
SELECT EXISTS (SELECT 1 FROM subscriptions WHERE id = $1);

-- name: GetCitiesForUpdate :many
SELECT DISTINCT city FROM subscriptions WHERE confirmed = true AND frequency = $1;
