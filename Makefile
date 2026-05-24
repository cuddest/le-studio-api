run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./... -v -cover

lint:
	golangci-lint run

migrate:
	migrate -path migrations -database "$${DATABASE_URL}" up

seed:
	go run ./cmd/server

docker-up:
	docker compose up -d
