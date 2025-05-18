package main

import (
	"database/sql"
	"log"

	"github.com/hibiken/asynq"
	_ "github.com/lib/pq" // Driver for PostgreSQL
	"github.com/mishvets/WeatherForecast/api"
	db "github.com/mishvets/WeatherForecast/db/sqlc"
	"github.com/mishvets/WeatherForecast/mailer"
	"github.com/mishvets/WeatherForecast/service"
	"github.com/mishvets/WeatherForecast/util"
	"github.com/mishvets/WeatherForecast/worker"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	weatherService := service.NewServiceWeather(config.WeatherApiUrl, config.WeatherApiKey)
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(redisOpt, store, config, weatherService)
	go runTaskScheduler(redisOpt)
	
	server := api.NewServer(store, taskDistributor, weatherService)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, config util.Config, weatherService service.Service) {
	mailer := mailer.NewGmailSender(config.EmailSenderName, config.EmailSenderAdress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer, weatherService)
	log.Print("start task processor")
	if err := taskProcessor.Start(); err != nil {
		log.Fatal("cannot start task processor: ", err)
	}
}

func runTaskScheduler(redisOpt asynq.RedisClientOpt) {
	taskScheduler := worker.NewRedisScheduler(redisOpt)
	err := taskScheduler.Start()
	if err != nil {
		log.Fatal("cannot start taskScheduler: ", err)
	}
}
