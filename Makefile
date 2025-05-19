postgres:
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

createdb:
	docker exec -it postgres17 createdb --username=root --owner=root weather_forecast

dropdb:
	docker exec -it postgres17 dropdb weather_forecast

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/weather_forecast?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/weather_forecast?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

retest:
	go test -count=1 -v -cover -short ./...

server:
	go run main.go

redis:
	docker run --name redis -p 6379:6379 -d redis:8-alpine

mock:
	mockgen --package mockdb --destination db/mock/store.go github.com/mishvets/WeatherForecast/db/sqlc Store
	mockgen --package mockworker --destination worker/mock/distributor.go github.com/mishvets/WeatherForecast/worker TaskDistributor
	mockgen --package mockservice --destination service/mock/service.go github.com/mishvets/WeatherForecast/service Service

.PHONY: postgres createdb dropdb migrateup migratedown sqlc retest test server redis mock
