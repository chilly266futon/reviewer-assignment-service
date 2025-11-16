package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/zap"

	"github.com/chilly266futon/reviewer-assignment-service/internal/domain"
	"github.com/chilly266futon/reviewer-assignment-service/internal/repository"
	pkgErrors "github.com/chilly266futon/reviewer-assignment-service/pkg/errors"
)

type PRService struct {
	prRepo   repository.PullRequestRepository
	userRepo repository.UserRepository
	logger   *zap.Logger
	rng      *rand.Rand
}

func NewPRService(prRepo repository.PullRequestRepository, userRepo repository.UserRepository, logger *zap.Logger) *PRService {
	// Инициализируем генератор случайных чисел
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &PRService{
		prRepo:   prRepo,
		userRepo: userRepo,
		logger:   logger,
		rng:      rng,
	}
}

// CreatePR создает новый PR с автоматическим назначением ревьюеров
func (s *PRService) CreatePR(ctx context.Context, prID, name, authorID string) (*domain.PullRequest, error) {
	// Валидация
	if prID == "" || name == "" || authorID == "" {
		return nil, pkgErrors.ErrInvalidInput
	}

	// Получаем автора и его команду
	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgErrors.ErrNotFound
		}
		s.logger.Error("failed to get author",
			zap.String("author_id", authorID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get author: %w", err)
	}

	s.logger.Debug("author found",
		zap.String("author_id", author.ID),
		zap.Int("team_id", author.TeamID),
	)

	// Получаем активных кандидатов из команды автора (исключая автора)
	candidates, err := s.userRepo.GetActiveUsersByTeamID(ctx, author.TeamID, []string{authorID})
	if err != nil {
		s.logger.Error("failed to get reviewer candidates",
			zap.Int("team_id", author.TeamID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get candidates: %w", err)
	}

	s.logger.Debug("reviewers candidates found",
		zap.Int("count", len(candidates)),
	)

	// Выбираем до 2 случайных ревьюеров
	reviewerIDs := s.selectRandomReviewers(candidates, 2)

	s.logger.Info("reviewers selected",
		zap.String("pr_id", prID),
		zap.Strings("reviewer_ids", reviewerIDs),
	)

	// Создаем PR
	now := time.Now()
	pr := &domain.PullRequest{
		ID:                prID,
		Name:              name,
		AuthorID:          authorID,
		Status:            domain.StatusOpen,
		AssignedReviewers: reviewerIDs,
		CreatedAt:         now,
	}

	if err := s.prRepo.Create(ctx, pr, reviewerIDs); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return nil, pkgErrors.ErrPRExists
		}
		s.logger.Error("failed to create PR",
			zap.String("pr_id", prID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("create PR: %w", err)
	}

	s.logger.Info("PR created",
		zap.String("pr_id", prID),
		zap.String("author_id", authorID),
		zap.Int("reviewers_count", len(reviewerIDs)),
	)

	return pr, nil
}

// SelectRandomReviewers выбирает случайных ревьюеров из списка кандидатов
func (s *PRService) selectRandomReviewers(candidates []*domain.User, maxCount int) []string {
	n := len(candidates)
	if n == 0 {
		return []string{}
	}

	// Определяем сколько ревьюеров выбрать
	count := n
	if count > maxCount {
		count = maxCount
	}

	perm := s.rng.Perm(n)

	result := make([]string, count)
	for i := 0; i < count; i++ {
		result[i] = candidates[perm[i]].ID
	}

	return result
}

// MergePR идемпотентно мержит PR
func (s *PRService) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	if prID == "" {
		return nil, pkgErrors.ErrInvalidInput
	}

	// Получаем PR
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgErrors.ErrNotFound
		}
		s.logger.Error("failed to get PR",
			zap.String("pr_id", prID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("get PR: %w", err)
	}

	if pr.Status == domain.StatusMerged {
		s.logger.Debug("PR already merged",
			zap.String("pr_id", prID),
		)
		return pr, nil // Уже смержен
	}

	// Обновляем статус
	now := time.Now()
	if err := s.prRepo.UpdateStatus(ctx, prID, domain.StatusMerged, &now); err != nil {
		s.logger.Error("failed to merge PR",
			zap.String("pr_id", prID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("merge PR: %w", err)
	}

	s.logger.Info("PR merged",
		zap.String("pr_id", prID),
	)

	// Получаем обновленный PR
	pr, err = s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR: %w", err)
	}

	return pr, nil
}

func (s *PRService) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (string, *domain.PullRequest, error) {
	if prID == "" || oldReviewerID == "" {
		return "", nil, pkgErrors.ErrInvalidInput
	}

	// Получаем PR
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", nil, pkgErrors.ErrNotFound
		}
		s.logger.Error("failed to get PR",
			zap.String("pr_id", prID),
			zap.Error(err),
		)
		return "", nil, fmt.Errorf("get PR: %w", err)
	}

	// Проверяем статус PR
	if pr.Status == domain.StatusMerged {
		return "", nil, pkgErrors.ErrPRMerged
	}

	// Проверяем, что старый ревьюер назначен на PR
	if !contains(pr.AssignedReviewers, oldReviewerID) {
		return "", nil, pkgErrors.ErrNotAssigned
	}

	// Получаем команду заменяемого ревьюера
	oldReviewer, err := s.userRepo.GetByID(ctx, oldReviewerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", nil, pkgErrors.ErrNotFound
		}
		return "", nil, fmt.Errorf("get old reviewer: %w", err)
	}

	s.logger.Debug("old reviewer found",
		zap.String("reviewer_id", oldReviewer.ID),
		zap.Int("team_id", oldReviewer.TeamID),
	)

	// Список исключаемых пользователей (уже назначенные ревьюеры и автор)
	excludeUserIDs := append([]string{}, pr.AssignedReviewers...)
	excludeUserIDs = append(excludeUserIDs, pr.AuthorID)

	// Получаем активных кандидатов из команды старого ревьюера (исключая уже назначенных и автора)
	candidates, err := s.userRepo.GetActiveUsersByTeamID(ctx, oldReviewer.TeamID, excludeUserIDs)
	if err != nil {
		s.logger.Error("failed to get replacement candidates",
			zap.Int("team_id", oldReviewer.TeamID),
			zap.Error(err),
		)
		return "", nil, fmt.Errorf("get candidates: %w", err)
	}

	// Выбираем случайного нового ревьюера
	newReviewerID := s.selectRandomReviewer(candidates)
	if newReviewerID == "" {
		s.logger.Warn("no replacement candidates available",
			zap.String("pr_id", prID),
			zap.String("old_reviewer_id", oldReviewerID),
			zap.Int("team_id", oldReviewer.TeamID),
		)
		return "", nil, pkgErrors.ErrNoCandidate
	}

	s.logger.Info("new reviewer selected",
		zap.String("pr_id", prID),
		zap.String("old_reviewer_id", oldReviewerID),
		zap.String("new_reviewer_id", newReviewerID),
	)

	// Заменить ревьюера в PR
	if err := s.prRepo.ReplaceReviewer(ctx, prID, oldReviewerID, newReviewerID); err != nil {
		s.logger.Error("failed to replace reviewer",
			zap.String("pr_id", prID),
			zap.Error(err),
		)
		return "", nil, fmt.Errorf("replace reviewer: %w", err)
	}

	// Получаем обновленный PR
	pr, err = s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return "", nil, fmt.Errorf("get updated PR: %w", err)
	}

	return newReviewerID, pr, nil

}

// contains проверяет наличие элемента в срезе
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (s *PRService) selectRandomReviewer(candidates []*domain.User) string {
	if len(candidates) == 0 {
		return ""
	}
	idx := s.rng.Intn(len(candidates))
	return candidates[idx].ID
}
