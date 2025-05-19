CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TYPE frequency_enum AS ENUM ('daily', 'hourly');

CREATE TABLE subscriptions (
    -- Check is id really needed
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
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
    humidity SMALLINT NOT NULL,
    description TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX ON weather_data(city);
CREATE INDEX ON subscriptions(email);
CREATE INDEX ON subscriptions(token);

ALTER TABLE subscriptions ALTER COLUMN token SET DEFAULT gen_random_uuid();