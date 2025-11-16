package service

import "fmt"

type CreateTeamInput struct {
	TeamName string
	Members  []TeamMemberInput
}

type TeamMemberInput struct {
	UserID   string
	Username string
	IsActive bool
}

func (i *CreateTeamInput) Validate() error {
	if i.TeamName == "" {
		return fmt.Errorf("team_name is required")
	}
	// Ограничение формата
	if len(i.TeamName) > 100 {
		return fmt.Errorf("team_name too long (max 100 characters)")
	}

	if len(i.Members) == 0 {
		return fmt.Errorf("at least one member required")
	}

	seen := make(map[string]bool)
	for _, m := range i.Members {
		if m.UserID == "" {
			return fmt.Errorf("user_id is required for all members")
		}
		if m.Username == "" {
			return fmt.Errorf("username is required for all members")
		}

		// Проверка дубликатов в запросе
		if seen[m.UserID] {
			return fmt.Errorf("duplicate user_id: %s", m.UserID)
		}
		seen[m.UserID] = true
	}
	return nil
}
