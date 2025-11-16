package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/chilly266futon/reviewer-assignment-service/internal/domain"
	"github.com/chilly266futon/reviewer-assignment-service/internal/repository"
	pkgErrors "github.com/chilly266futon/reviewer-assignment-service/pkg/errors"
)

type TeamService struct {
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
	logger   *zap.Logger
}

func NewTeamService(teamRepo repository.TeamRepository, userRepo repository.UserRepository, logger *zap.Logger) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
		logger:   logger,
	}
}

// CreateTeam создаёт команду с участниками
func (s *TeamService) CreateTeam(ctx context.Context, input *CreateTeamInput) (*domain.Team, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", pkgErrors.ErrInvalidInput, err)
	}

	// Проверяем, существует ли команда с таким именем
	existingTeam, err := s.teamRepo.GetByName(ctx, input.TeamName)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		s.logger.Error("failed to check team existence",
			zap.String("team_name", input.TeamName),
			zap.Error(err),
		)
		return nil, fmt.Errorf("check team: %w", err)
	}

	if existingTeam != nil {
		return nil, pkgErrors.ErrTeamExists
	}

	// Создаем команду
	now := time.Now()
	team := &domain.Team{
		Name:      input.TeamName,
		CreatedAt: now,
	}

	if err := s.teamRepo.Create(ctx, team); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return nil, pkgErrors.ErrTeamExists
		}
		s.logger.Error("failed to create team",
			zap.String("team_name", input.TeamName),
			zap.Error(err),
		)
		return nil, fmt.Errorf("create team: %w", err)
	}

	s.logger.Info("team created",
		zap.String("team_name", team.Name),
		zap.Int("team_id", team.ID),
	)

	// Создаем участников
	var teamMembers []*domain.User
	for _, m := range input.Members {
		user := &domain.User{
			ID:        m.UserID,
			Username:  m.Username,
			TeamID:    team.ID,
			TeamName:  team.Name,
			IsActive:  m.IsActive,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			s.logger.Error("failed to create/update user",
				zap.String("user_id", user.ID),
				zap.Error(err),
			)
			return nil, fmt.Errorf("create/update user: %w", err)
		}

		teamMembers = append(teamMembers, user)
	}

	s.logger.Info("team members created/updated",
		zap.String("team_name", team.Name),
		zap.Int("count", len(teamMembers)),
	)

	team.Members = teamMembers
	return team, nil
}

// GetTeam возвращает команду по названию
func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	if teamName == "" {
		return nil, pkgErrors.ErrInvalidInput
	}

	team, err := s.teamRepo.GetByName(ctx, teamName)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgErrors.ErrNotFound
		}
		s.logger.Error("failed to get team",
			zap.String("team_name", teamName),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get team: %w", err)
	}

	return team, nil
}
