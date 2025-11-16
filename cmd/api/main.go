// cmd/api/main.go
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

	"go.uber.org/zap"

	"github.com/chilly266futon/reviewer-assignment-service/internal/config"
	"github.com/chilly266futon/reviewer-assignment-service/internal/handler"
	"github.com/chilly266futon/reviewer-assignment-service/internal/repository/postgres"
	"github.com/chilly266futon/reviewer-assignment-service/internal/service"
	"github.com/chilly266futon/reviewer-assignment-service/pkg/logger"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Инициализация логгера
	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logger: %v\n", err)
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

	// Создаем transaction manager и репозитории
	txManager := postgres.NewTxManager(pool)
	userRepo := postgres.NewUserRepository(pool, log)
	teamRepo := postgres.NewTeamRepository(pool, log)
	prRepo := postgres.NewPRRepository(pool, txManager, log)

	log.Info("repositories initialized")

	// Инициализируем сервисы
	teamService := service.NewTeamService(teamRepo, userRepo, log)
	userService := service.NewUserService(userRepo, prRepo, log)
	prService := service.NewPRService(prRepo, userRepo, log)

	log.Info("services initialized")

	// Создаем router
	router := handler.NewRouter(teamService, userService, prService, pool, log)

	log.Info("router configured")

	// Создаем HTTP сервер
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер в горутине
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
