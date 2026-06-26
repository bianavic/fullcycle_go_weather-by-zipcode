.PHONY: run build test test-coverage tidy fmt docker-up docker-down

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test -race -cover ./...

test-coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

tidy:
	go mod tidy

fmt:
	gofmt -w .
	goimports -w .

docker-up:
	docker compose up --build

docker-down:
	docker compose down
