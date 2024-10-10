package models

import "time"

// UserProject represents a GitHub project added by a user
type UserProject struct {
	ID          int       `json:"id"`
	UserID      string    `json:"user_id"`
	ProjectType string    `json:"project_type"` // 'personal' or 'org'
	Username    string    `json:"username"`
	PAT         string    `json:"pat"`        // Store encrypted PAT
	RepoNames   string    `json:"repo_names"` // Comma-separated repo names
	CreatedAt   time.Time `json:"created_at"`
}
