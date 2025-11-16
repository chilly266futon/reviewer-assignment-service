package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/chilly266futon/reviewer-assignment-service/internal/config"
	"github.com/chilly266futon/reviewer-assignment-service/internal/repository/postgres"
	"github.com/chilly266futon/reviewer-assignment-service/pkg/logger"
)

func main() {
	// Загрузка .env только для локальной разработки, игнорируем ошибку в Docker
	_ = godotenv.Load()

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Инициализация логгера
	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting reviewer-assignment-service...",
		zap.String("port", cfg.ServerPort),
		zap.String("log_level", cfg.LogLevel),
	)

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Подключаемся к БД
	log.Info("Connecting to database",
		zap.String("host", cfg.DBHost),
		zap.String("database", cfg.DBName),
	)

	pool, err := postgres.NewPool(ctx, cfg, log)
	if err != nil {
		log.Fatal("failed to create database pool", zap.Error(err))
	}
	defer postgres.Close(pool)

	log.Info("database connection established")

	// Создаем transaction manager
	txManager := postgres.NewTxManager(pool)

	// Инициализация репозиториев
	userRepo := postgres.NewUserRepository(pool, log)
	teamRepo := postgres.NewTeamRepository(pool, log)
	prRepo := postgres.NewPRRepository(pool, txManager, log)

	log.Info("repositories initialized")

	// TODO: инициализация сервисов (next task)
	_ = userRepo
	_ = teamRepo
	_ = prRepo

	router := chi.NewRouter()

	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// Проверяем подключение к БД
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"error","database":"disconnected"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","database":"connected"}`))
	})

	// TODO: routes

	// Создаем HTTP сервер
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("Server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server stopped gracefully")
}
