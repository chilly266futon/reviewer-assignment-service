.PHONY: help
help: ## Показать help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Сборка
	go build -o bin/api ./cmd/api

.PHONY: run
run: ## Запустить локально
	go run ./cmd/api

.PHONY: docker-up
docker-up: ## Поднять в docker-compose
	docker-compose up --build

.PHONY: docker-down
docker-down: ## Остановить docker-compose
	docker-compose down -v

.PHONY: test
test: ## Запустить тесты
	go test -v -race -timeout 30s ./...

.PHONY: test-coverage
test-coverage: ## Запустить тесты с покрытием
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: lint
lint: ## Запустить линтер
	golangci-lint run ./...

.PHONY: fmt
fmt: ## Отформатировать код
	gofmt -w .
	goimports -w -local github.com/chilly266futon/reviewer-assignment-service .

.PHONY: clean
clean: ## Очистить артефакты
	rm -rf bin/
	rm -f coverage.out coverage.html

.PHONY: db-up
db-up: ## Запустить только PostgreSQL и миграции
	docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 5
	docker-compose up migrate

.PHONY: db-down
db-down: ## Остановить PostgreSQL
	docker-compose down -v

.PHONY: db-reset
db-reset:  ## Пересоздать БД с нуля
	docker-compose down -v
	docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 5
	docker-compose up migrate

.PHONY: db-shell
db-shell: ## Открыть psql консоль
	docker-compose exec postgres psql -U reviewer_user -d reviewer_db

.PHONY: db-logs
db-logs: ## Показать логи PostgreSQL
	docker-compose logs -f postgres