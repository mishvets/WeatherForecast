CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TYPE frequency_enum AS ENUM ('daily', 'hourly');

CREATE TABLE subscriptions (
    -- Check is id really needed
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    frequency frequency_enum NOT NULL,
    confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    token UUID DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE weather_data (
    id BIGSERIAL PRIMARY KEY,
    city VARCHAR(100) NOT NULL UNIQUE,
    temperature REAL NOT NULL,
    humidity INT NOT NULL,
    description TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CREATE TABLE weather_updates (
--     id SERIAL PRIMARY KEY,
--     subscription_id INT REFERENCES subscriptions(id) ON DELETE CASCADE,
--     sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     weather_snapshot JSONB NOT NULL
-- );

CREATE INDEX ON weather_data(city);
CREATE INDEX ON subscriptions(email);
CREATE INDEX ON subscriptions(token);

ALTER TABLE subscriptions ALTER COLUMN token SET DEFAULT gen_random_uuid();