package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/chilly266futon/reviewer-assignment-service/internal/dto"
	"github.com/chilly266futon/reviewer-assignment-service/internal/service"
)

// TeamHandler обрабатывает запросы для команд
type TeamHandler struct {
	teamService *service.TeamService
	logger      *zap.Logger
}

// NewTeamHandler создаёт новый handler для команд
func NewTeamHandler(teamService *service.TeamService, logger *zap.Logger) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
		logger:      logger,
	}
}

// Add создаёт команду с участниками
func (h *TeamHandler) Add(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTeamRequest

	// Парсинг JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode request", zap.Error(err))
		respondError(w, "INVALID_REQUEST", "invalid JSON", http.StatusBadRequest)
		return
	}

	// HTTP-уровень валидация
	if err := req.Validate(); err != nil {
		respondError(w, "INVALID_REQUEST", err.Error(), http.StatusBadRequest)
		return
	}

	// Маппинг DTO → Service Input
	input := &service.CreateTeamInput{
		TeamName: req.TeamName,
		Members:  make([]service.TeamMemberInput, len(req.Members)),
	}

	for i, m := range req.Members {
		input.Members[i] = service.TeamMemberInput{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}

	// Вызов service
	team, err := h.teamService.CreateTeam(r.Context(), input)
	if err != nil {
		handleServiceError(w, err, h.logger)
		return
	}

	// Успешный ответ
	respondJSON(w, dto.TeamResponse{Team: team}, http.StatusCreated)
}

// Get возвращает команду по имени
func (h *TeamHandler) Get(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		respondError(w, "INVALID_REQUEST", "team_name parameter is required", http.StatusBadRequest)
		return
	}

	team, err := h.teamService.GetTeam(r.Context(), teamName)
	if err != nil {
		handleServiceError(w, err, h.logger)
		return
	}

	respondJSON(w, team, http.StatusOK)
}
