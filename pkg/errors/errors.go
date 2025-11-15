package errors

import "errors"

// Domain errors - используются в сервисах

var (
	ErrNotFound     = errors.New("not found")
	ErrTeamExists   = errors.New("team already exists")
	ErrPRExists     = errors.New("pull request already exists")
	ErrPRMerged     = errors.New("cannot modify merged PR")
	ErrNotAssigned  = errors.New("reviewer not assigned to PR")
	ErrNoCandidate  = errors.New("no candidate available for reassignment")
	ErrInvalidInput = errors.New("invalid input")
)

// Коды ошибок для API (из OpenAPI)

const (
	CodeTeamExists  = "TEAM_EXISTS"
	CodePRExists    = "PR_EXISTS"
	CodePRMerged    = "PR_MERGED"
	CodeNotAssigned = "NOT_ASSIGNED"
	CodeNoCandidate = "NO_CANDIDATE"
	CodeNotFound    = "NOT_FOUND"
)

// MapErrorToCode мапит доменную ошибку в API код ошибки
func MapErrorToCode(err error) string {
	switch {
	case errors.Is(err, ErrTeamExists):
		return CodeTeamExists
	case errors.Is(err, ErrPRExists):
		return CodePRExists
	case errors.Is(err, ErrPRMerged):
		return CodePRMerged
	case errors.Is(err, ErrNotAssigned):
		return CodeNotAssigned
	case errors.Is(err, ErrNoCandidate):
		return CodeNoCandidate
	case errors.Is(err, ErrNotFound):
		return CodeNotFound
	default:
		return "INTERNAL_ERROR"
	}
}
