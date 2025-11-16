package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/chilly266futon/reviewer-assignment-service/internal/dto"
	"github.com/chilly266futon/reviewer-assignment-service/internal/service"
)

type UserHandler struct {
	userService *service.UserService
	logger      *zap.Logger
}

func NewUserHandler(userService *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// SetIsActive изменяет статус активности пользователя
func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req dto.SetIsActiveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode request", zap.Error(err))
		respondError(w, "INVALID_REQUEST", "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, "INVALID_REQUEST", err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.SetIsActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		handleServiceError(w, err, h.logger)
		return
	}

	respondJSON(w, dto.UserResponse{User: user}, http.StatusOK)
}

// GetReview возвращает PR'ы где пользователь назначен ревьюером
func (h *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		respondError(w, "INVALID_REQUEST", "user_id parameter is required", http.StatusBadRequest)
		return
	}

	prs, err := h.userService.GetUserReviews(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err, h.logger)
		return
	}

	// Маппинг к краткому формату
	prShorts := make([]*dto.PRShort, len(prs))
	for i, pr := range prs {
		prShorts[i] = dto.ToPRShort(pr)
	}

	respondJSON(w, dto.UserReviewsResponse{
		UserID:       userID,
		PullRequests: prShorts,
	}, http.StatusOK)
}
