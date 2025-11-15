package dto

import "github.com/chilly266futon/reviewer-assignment-service/internal/domain"

// ErrorResponse - стандартный формат ошибки (из OpenAPI)
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// TeamResponse - ответ с информацией о команде
type TeamResponse struct {
	Team *domain.Team `json:"team"`
}

// UserResponse - ответ с информацией о пользователе
type UserResponse struct {
	User *domain.User `json:"user"`
}

// PRResponse - ответ с информацией о PR
type PRResponse struct {
	PR *domain.PullRequest `json:"pull_request"`
}

// ReassignResponse - ответ с информацией о переназначенном ревьюере
type ReassignResponse struct {
	PR         *domain.PullRequest `json:"pull_request"`
	ReplacedBy string              `json:"replaced_by"`
}

// UserReviewsResponse - ответ с информацией о PR, ожидающих ревью от пользователя
type UserReviewsResponse struct {
	UserID       string     `json:"user_id"`
	PullRequests []*PRShort `json:"pull_requests"`
}

// PRShort - краткая информация о PR
type PRShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

// ToPRShort конвертирует domain.PullRequest в краткий формат
func ToPRShort(pr *domain.PullRequest) *PRShort {
	return &PRShort{
		PullRequestID:   pr.ID,
		PullRequestName: pr.Name,
		AuthorID:        pr.AuthorID,
		Status:          pr.Status,
	}
}
