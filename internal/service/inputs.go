package service

import "fmt"

// CreateTeamInput входные данные для создания команды
type CreateTeamInput struct {
	TeamName string
	Members  []TeamMemberInput
}

// TeamMemberInput данные участника команды
type TeamMemberInput struct {
	UserID   string
	Username string
	IsActive bool
}

// Validate валидация на уровне service (бизнес-правила)
func (i *CreateTeamInput) Validate() error {
	if i.TeamName == "" {
		return fmt.Errorf("team_name is required")
	}

	if len(i.TeamName) > 100 {
		return fmt.Errorf("team_name too long (max 100 characters)")
	}

	// Уникальность user_id внутри запроса
	seen := make(map[string]bool)
	for _, m := range i.Members {
		if m.UserID == "" {
			return fmt.Errorf("user_id is required for all members")
		}
		if m.Username == "" {
			return fmt.Errorf("username is required for all members")
		}

		if seen[m.UserID] {
			return fmt.Errorf("duplicate user_id in request: %s", m.UserID)
		}
		seen[m.UserID] = true
	}

	return nil
}

// CreatePRInput входные данные для создания PR
type CreatePRInput struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
}

func (i *CreatePRInput) Validate() error {
	if i.PullRequestID == "" {
		return fmt.Errorf("pull_request_id is required")
	}
	if i.PullRequestName == "" {
		return fmt.Errorf("pull_request_name is required")
	}
	if i.AuthorID == "" {
		return fmt.Errorf("author_id is required")
	}

	if len(i.PullRequestID) > 100 {
		return fmt.Errorf("pull_request_id too long")
	}

	return nil
}

// SetIsActiveInput входные данные для изменения статуса активности
type SetIsActiveInput struct {
	UserID   string
	IsActive bool
}

func (i *SetIsActiveInput) Validate() error {
	if i.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	return nil
}

// ReassignReviewerInput входные данные для переназначения ревьюера
type ReassignReviewerInput struct {
	PullRequestID string
	OldReviewerID string
}

func (i *ReassignReviewerInput) Validate() error {
	if i.PullRequestID == "" {
		return fmt.Errorf("pull_request_id is required")
	}
	if i.OldReviewerID == "" {
		return fmt.Errorf("old_reviewer_id is required")
	}
	return nil
}

// MergePRInput входные данные для merge PR
type MergePRInput struct {
	PullRequestID string
}

func (i *MergePRInput) Validate() error {
	if i.PullRequestID == "" {
		return fmt.Errorf("pull_request_id is required")
	}
	return nil
}
