package domain

import "time"

type Team struct {
	ID        int       `json:"-"`
	Name      string    `json:"team_name"`
	Members   []*User   `json:"members"`
	CreatedAt time.Time `json:"-"`
}
