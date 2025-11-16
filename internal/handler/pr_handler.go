package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/chilly266futon/reviewer-assignment-service/internal/dto"
	"github.com/chilly266futon/reviewer-assignment-service/internal/service"
)

// PRHandler обрабатывает запросы для Pull Requests
type PRHandler struct {
	prService *service.PRService
	logger    *zap.Logger
}

// NewPRHandler создаёт новый handler для PR
func NewPRHandler(prService *service.PRService, logger *zap.Logger) *PRHandler {
	return &PRHandler{
		prService: prService,
		logger:    logger,
	}
}

// Create создаёт PR с автоматическим назначением ревьюеров
func (h *PRHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePRRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode request", zap.Error(err))
		respondError(w, "INVALID_REQUEST", "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, "INVALID_REQUEST", err.Error(), http.StatusBadRequest)
		return
	}

	pr, err := h.prService.CreatePR(r.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		handleServiceError(w, err, h.logger)
		return
	}

	respondJSON(w, dto.PRResponse{PR: pr}, http.StatusCreated)
}

// Merge помечает PR как смерженный
func (h *PRHandler) Merge(w http.ResponseWriter, r *http.Request) {
	var req dto.MergePRRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode request", zap.Error(err))
		respondError(w, "INVALID_REQUEST", "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, "INVALID_REQUEST", err.Error(), http.StatusBadRequest)
		return
	}

	pr, err := h.prService.MergePR(r.Context(), req.PullRequestID)
	if err != nil {
		handleServiceError(w, err, h.logger)
		return
	}

	respondJSON(w, dto.PRResponse{PR: pr}, http.StatusOK)
}

// Reassign переназначает ревьюера
func (h *PRHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	var req dto.ReassignReviewerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode request", zap.Error(err))
		respondError(w, "INVALID_REQUEST", "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, "INVALID_REQUEST", err.Error(), http.StatusBadRequest)
		return
	}

	newReviewerID, pr, err := h.prService.ReassignReviewer(r.Context(), req.PullRequestID, req.OldReviewerID)
	if err != nil {
		handleServiceError(w, err, h.logger)
		return
	}

	respondJSON(w, dto.ReassignResponse{
		PR:         pr,
		ReplacedBy: newReviewerID,
	}, http.StatusOK)
}
