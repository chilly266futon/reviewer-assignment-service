package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	pkgErrors "github.com/chilly266futon/reviewer-assignment-service/internal/dto"
	serviceErrors "github.com/chilly266futon/reviewer-assignment-service/pkg/errors"
)

// respondJSON отправляет JSON ответ
func respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Если не можем закодировать ответ - логируем, но уже поздно что-то отправлять
		return
	}
}

// respondError отправляет стандартизированную ошибку
func respondError(w http.ResponseWriter, code, message string, status int) {
	respondJSON(w, pkgErrors.ErrorResponse{
		Error: pkgErrors.ErrorDetail{
			Code:    code,
			Message: message,
		},
	}, status)
}

// handleServiceError обрабатывает ошибки из service layer
func handleServiceError(w http.ResponseWriter, err error, logger *zap.Logger) {
	switch {
	case errors.Is(err, serviceErrors.ErrNotFound):
		respondError(w, serviceErrors.CodeNotFound, err.Error(), http.StatusNotFound)
	case errors.Is(err, serviceErrors.ErrTeamExists):
		respondError(w, serviceErrors.CodeTeamExists, err.Error(), http.StatusBadRequest)
	case errors.Is(err, serviceErrors.ErrPRExists):
		respondError(w, serviceErrors.CodePRExists, err.Error(), http.StatusConflict)
	case errors.Is(err, serviceErrors.ErrPRMerged):
		respondError(w, serviceErrors.CodePRMerged, err.Error(), http.StatusConflict)
	case errors.Is(err, serviceErrors.ErrNotAssigned):
		respondError(w, serviceErrors.CodeNotAssigned, err.Error(), http.StatusConflict)
	case errors.Is(err, serviceErrors.ErrNoCandidate):
		respondError(w, serviceErrors.CodeNoCandidate, err.Error(), http.StatusConflict)
	case errors.Is(err, serviceErrors.ErrInvalidInput):
		respondError(w, "INVALID_REQUEST", err.Error(), http.StatusBadRequest)
	default:
		logger.Error("internal server error", zap.Error(err))
		respondError(w, "INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)
	}
}
