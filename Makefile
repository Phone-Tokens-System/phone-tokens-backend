SERVICE_DIR=.
BINARY=app

.PHONY: build run docker-build up down

build:
	cd $(SERVICE_DIR) && go build -o $(BINARY) ./cmd

run:
	cd $(SERVICE_DIR) && go run ./cmd

docker-build:
	docker build -t app .

up:
	docker-compose up --build

down:
	docker-compose down
