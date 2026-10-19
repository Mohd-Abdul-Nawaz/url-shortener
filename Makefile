.PHONY: build run test docker-up docker-down

build:
	go build -o url-shortener ./cmd/api

run:
	REDIS_ADDR=localhost:6379 PORT=8080 go run ./cmd/api

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

test:
	go test ./...

tidy:
	go mod tidy