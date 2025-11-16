.PHONY: build test lint docker-build up down clean coverage

# Сборка бинарного файла
build:
	@echo "Building application..."
	@go build -o bin/api cmd/api/main.go

# Запуск тестов
test:
	@echo "Running tests..."
	@go test -v -race -short ./...

# Запуск тестов с покрытием
coverage:
	@echo "Running tests with coverage..."
	@go test -v -short -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out | tail -1

# Линтер
lint:
	@echo "Running linter..."
	@golangci-lint run --timeout=5m

# Сборка Docker образа
docker-build:
	@echo "Building Docker image..."
	@docker build -t reviewer-assignment-service:latest .

# Запуск всех сервисов
up:
	@echo "Starting services..."
	@docker-compose up --build

# Остановка контейнеров
down:
	@echo "Stopping services..."
	@docker-compose down

# Очистка
clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -f coverage.out
	@go clean

# Установка зависимостей
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Форматирование кода
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Запуск миграций
migrate-up:
	@echo "Running migrations..."
	@migrate -path ./migrations -database "postgres://reviewer_user:reviewer_password@localhost:5432/reviewer_db?sslmode=disable" up

migrate-down:
	@echo "Rolling back migrations..."
	@migrate -path ./migrations -database "postgres://reviewer_user:reviewer_password@localhost:5432/reviewer_db?sslmode=disable" down

# Помощь
help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make test          - Run tests"
	@echo "  make coverage      - Run tests with coverage"
	@echo "  make lint          - Run linter"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make up            - Start all services with docker-compose"
	@echo "  make down          - Stop all services"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make deps          - Download dependencies"
	@echo "  make fmt           - Format code"
	@echo "  make migrate-up    - Run database migrations"
	@echo "  make migrate-down  - Rollback database migrations"

