package domain

import "time"

const (
	StatusOpen   = "OPEN"
	StatusMerged = "MERGED"
)

type PullRequest struct {
	ID                string     `json:"pull_request_id"`
	Name              string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         time.Time  `json:"created_at,omitempty"`
	MergedAt          *time.Time `json:"merged_at,omitempty"` // nullable
}

// IsOpen проверяет, открыт ли PR
func (pr *PullRequest) IsOpen() bool {
	return pr.Status == StatusOpen
}

// IsMerged проверяет, смержен ли PR
func (pr *PullRequest) IsMerged() bool {
	return pr.Status == StatusMerged
}
