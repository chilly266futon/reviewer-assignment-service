FROM golang:1.24.2-alpine AS builder

# Установка зависимостей
RUN apk add --no-cache git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка бинарного файла
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/api ./cmd/api

FROM alpine:3.19

WORKDIR /app

# Копируем бинарник из стадии сборки
COPY --from=builder /bin/api /app/api

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["/app/api"]