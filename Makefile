SERVICE_DIR=users
BINARY=users-service

.PHONY: build run docker-build up down

build:
	cd $(SERVICE_DIR) && go build -o $(BINARY) ./cmd/users

run:
	cd $(SERVICE_DIR) && go run ./cmd/users

docker-build:
	docker build -t users-service ./users

up:
	docker-compose up --build

down:
	docker-compose down

