package service

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/chilly266futon/reviewer-assignment-service/internal/domain"
)

// Unit tests для приватных методов

func TestSelectRandomReviewers_ReturnsUpToMaxCount(t *testing.T) {
	svc := &PRService{
		rng:    rand.New(rand.NewSource(42)),
		logger: zap.NewNop(),
	}

	candidates := []*domain.User{
		{ID: "user1"},
		{ID: "user2"},
		{ID: "user3"},
		{ID: "user4"},
		{ID: "user5"},
	}

	t.Run("returns maxCount reviewers when enough candidates", func(t *testing.T) {
		result := svc.selectRandomReviewers(candidates, 2)
		assert.Len(t, result, 2)
		assert.NotEqual(t, result[0], candidates[1])
	})

	t.Run("returns all candidates when maxCount exceeds available", func(t *testing.T) {
		result := svc.selectRandomReviewers(candidates[:2], 5)
		assert.Len(t, result, 2)
	})

	t.Run("returns empty slice for empty candidates", func(t *testing.T) {
		result := svc.selectRandomReviewers([]*domain.User{}, 2)
		assert.Empty(t, result)
	})

}

func TestSelectRandomReviewer_SingleSelection(t *testing.T) {
	svc := &PRService{
		rng:    rand.New(rand.NewSource(42)),
		logger: zap.NewNop(),
	}

	candidates := []*domain.User{
		{ID: "user1"},
		{ID: "user2"},
		{ID: "user3"},
	}

	t.Run("selects one reviewer", func(t *testing.T) {
		result := svc.selectRandomReviewer(candidates)
		assert.NotEmpty(t, result)
		assert.Contains(t, []string{"user1", "user2", "user3"}, result)
	})

	t.Run("returns empty string for empty candidates", func(t *testing.T) {
		result := svc.selectRandomReviewer([]*domain.User{})
		assert.Empty(t, result)
	})
}

// Статистический тест (может быть flaky)
func TestSelectRandomReviewers_Distribution(t *testing.T) {
	// Проверяем что выбор равномерный
	candidates := []*domain.User{
		{ID: "user1"},
		{ID: "user2"},
		{ID: "user3"},
		{ID: "user4"},
		{ID: "user5"},
	}

	counts := make(map[string]int)
	iterations := 10000

	svc := NewPRService(nil, nil, zap.NewNop())

	for i := 0; i < iterations; i++ {
		selected := svc.selectRandomReviewers(candidates, 2)
		for _, id := range selected {
			counts[id]++
		}
	}

	expected := iterations * 2 / len(candidates) // Ожидаемое количество выборов на пользователя
	tolerance := float64(expected) * 0.1         // Допустимое отклонение 10%

	for id, count := range counts {
		diff := float64(count - expected)
		if diff < 0 {
			diff = -diff
		}
		assert.Less(t, diff, tolerance,
			"user %s selected %d times, expected around %d", id, count, expected)
	}
}
