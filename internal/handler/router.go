package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/chilly266futon/reviewer-assignment-service/internal/service"
)

// NewRouter создаёт и настраивает HTTP router
func NewRouter(
	teamService *service.TeamService,
	userService *service.UserService,
	prService *service.PRService,
	pool *pgxpool.Pool,
	logger *zap.Logger,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Логирование запросов
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			logger.Info("request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.Status()),
				zap.Duration("duration", time.Since(start)),
			)
		})
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			logger.Error("database health check failed", zap.Error(err))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"error","database":"disconnected"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","database":"connected"}`))
	})

	// Инициализируем handlers
	teamHandler := NewTeamHandler(teamService, logger)
	userHandler := NewUserHandler(userService, logger)
	prHandler := NewPRHandler(prService, logger)

	// API routes
	r.Route("/team", func(r chi.Router) {
		r.Post("/add", teamHandler.Add)
		r.Get("/get", teamHandler.Get)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", userHandler.SetIsActive)
		r.Get("/getReview", userHandler.GetReview)
	})

	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", prHandler.Create)
		r.Post("/merge", prHandler.Merge)
		r.Post("/reassign", prHandler.Reassign)
	})

	return r
}
