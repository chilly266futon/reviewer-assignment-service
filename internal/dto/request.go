package dto

// CreateTeamRequest - запрос на создание команды
type CreateTeamRequest struct {
	TeamName string              `json:"team_name"`
	Members  []TeamMemberRequest `json:"members"`
}

// TeamMemberRequest - информация о члене команды
type TeamMemberRequest struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// SetIsActiveRequest - запрос на изменение статуса активности пользователя
type SetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

// CreatePRRequest - запрос на создание PR
type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

// MergePRRequest - запрос на merge PR
type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

// ReassignReviewerRequest - запрос на переназначение ревьюера
type ReassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

// Методы для валидации запросов

func (r *CreateTeamRequest) Validate() error {
	if r.TeamName == "" {
		return ErrMissingField("team_name")
	}
	return nil
}

func (r *SetIsActiveRequest) Validate() error {
	if r.UserID == "" {
		return ErrMissingField("user_id")
	}
	return nil
}

func (r *CreatePRRequest) Validate() error {
	if r.PullRequestID == "" {
		return ErrMissingField("pull_request_id")
	}
	if r.PullRequestName == "" {
		return ErrMissingField("pull_request_name")
	}
	if r.AuthorID == "" {
		return ErrMissingField("author_id")
	}
	return nil
}

func (r *MergePRRequest) Validate() error {
	if r.PullRequestID == "" {
		return ErrMissingField("pull_request_id")
	}
	return nil
}

func (r *ReassignReviewerRequest) Validate() error {
	if r.PullRequestID == "" {
		return ErrMissingField("pull_request_id")
	}
	if r.OldUserID == "" {
		return ErrMissingField("old_user_id")
	}
	return nil
}
