# Weather Forecast API

A Go-based weather API application with email subscription functionality. Users can subscribe to receive weather updates (daily or hourly) for a chosen city via email.

## Tech Stack

* **[PostgreSQL](https://www.postgresql.org/)** for data persistence
* **[sqlc](https://github.com/sqlc-dev/sqlc)** for type-safe database queries
* **[Viper](https://github.com/spf13/viper)** for handle configuration
* **[Asynq](https://github.com/hibiken/asynq)** queueing tasks and processing them asynchronously with workers
* **[Redis](https://redis.io/)** for Asynq queue backend and scheduling
* **[email](https://github.com/jordan-wright/email)** for email sending
* **[Gin](https://gin-gonic.com/)** for handling and routing HTTP API requests
* **[Testify](https://github.com/stretchr/testify)** testing toolkit
* **[WeatherAPI](https://www.weatherapi.com/)** for weather data
* **[golang-migrate](https://github.com/golang-migrate/migrate)** for managing database migrations
* **[gomock](https://github.com/golang/mock)** for mocking

## Local Setup Guide

### Option 1: Full Local Installation

### Prerequisites

Before running the application locally, ensure you have installed the following:

- [Go](https://go.dev/doc/install) – Required for building and running the application.
- [Docker](https://docs.docker.com/engine/install/) – Used for managing PostgreSQL and Redis containers.
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) – CLI tool for database migrations.

### Installation Steps

1. Clone the repository:
   ```sh
   git clone https://github.com/mishvets/WeatherForecast.git
   ```

2. Navigate to the project directory:
   ```sh
   cd WeatherForecast
   ```

3. Start PostgreSQL:
   ```sh
   make postgres
   ```

4. Start Redis:
   ```sh
   make redis
   ```

5. Create the database:
   ```sh
   make createdb
   ```

6. Apply migrations:
   ```sh
   make migrateup
   ```

#### Running the Server

Start the application:
```sh
make server
```

---

### Option 2: Run with Docker Compose

Alternatively, you can start the application and dependencies using Docker Compose:

```sh
docker compose up
```

This will build and run the application along with PostgreSQL and Redis containers inside Docker.

## API Endpoints

### `GET /weather`
- **Description:** Get current weather for a specified city.
- **Parameters:** `city` (string) – City name for weather forecast.
- **Response:** JSON object with `temperature`, `humidity`, and `description`.

### `POST /subscribe`
- **Description:** Subscribe an email to receive weather updates.
- **Parameters:**
  - `email` (string) – Email address.
  - `city` (string) – City name.
  - `frequency` (string) – Update frequency (`hourly` or `daily`).
- **Response:** Confirmation of subscription.

### `GET /confirm/{token}`
- **Description:** Confirm email subscription using a token.
- **Parameters:** `token` (string) – Confirmation token.
- **Response:** Subscription confirmation status.

### `GET /unsubscribe/{token}`
- **Description:** Unsubscribe from weather updates using a token.
- **Parameters:** `token` (string) – Unsubscribe token.
- **Response:** Unsubscription confirmation.

## Database Schema

### `subscriptions`

| Column     | Type      | Description                           |
| ---------- | --------- | ------------------------------------- |
| id         | BIGSERIAL | Primary key                           |
| email      | VARCHAR   | User's email address                  |
| city       | VARCHAR   | Subscribed city                       |
| frequency  | ENUM      | `hourly` or `daily`                   |
| confirmed  | BOOLEAN   | Confirmation status                   |
| token      | UUID      | Token for confirmation/unsubscription |
| created_at | TIMESTAMP | Record creation timestamp             |

### `weather_data`

| Column      | Type      | Description         |
| ----------- | --------- | ------------------- |
| id          | BIGSERIAL | Primary key         |
| city        | VARCHAR   | Unique city name    |
| temperature | REAL      | Current temperature |
| humidity    | INT       | Current humidity    |
| description | TEXT      | Weather condition   |
| updated\_at | TIMESTAMP | Last update time    |

## Service Logic Overview

### Subscription Process
The service allows users to subscribe to weather updates, either **hourly** or **daily**.
1. When a user sends a subscription request, their data is stored in the `subscribers` table.
2. They receive a confirmation email with a personal token link for email verification.
3. Clicking the confirmation link triggers a `GET` request to confirm the subscription.
4. The user's status changes to `confirmed`, and their city is added to the `weather_data` table (if the city isn't present).

User requests are handled using **Gin Web Framework**, which provides **fast** and **efficient** request processing.

### Weather Data Retrieval
The service uses **asynq.Scheduler** to fetch weather updates from [WeatherAPI.com](https://www.weatherapi.com/) every hour or day.
1. When the scheduled timer event triggers, a task for updating is added to TaskProcessor. Unique list of cities is generated from confirmed subscribers matching the update frequency.
2. A request is sent to WeatherAPI for each city, and the retrieved data is stored in the database.
3. After storing the weather data, a separate task is created in **TaskDistributor** to send notifications to users.
4. If weather data for a city is unavailable, users receive empty results, with the `Description` field indicating that the city's weather information is missing.

## TODO / Improvements

- Send requests to [WeatherAPI.com](https://www.weatherapi.com/) for multiple cities in one query instead of separate requests (requires a paid subscription).

- Optimize City Name Handling. Different spellings of the same city generate separate requests to WeatherAPI. Using their location dictionary could help reduce redundant API calls.

- Send Re-subscription Email if no data for requested city.

- Check how long ago the weather was updated and if less than `delta` min, return data from the database rather than sending a new request to the API

- When the task is triggered for `daily` updates, all confirmed subscribers must be updated (regardless of the subscription type). For `hourly` updates, only those with the appropriate **frequency** must be selected. In the current implementation, they are collected separately by subscription type.
