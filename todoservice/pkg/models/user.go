package models

import (
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"time"
)

type User struct {
	ID        string         `json:"id"`
	Email     string         `json:"email"`
	GithubID  string         `json:"github_id"`
	Role      constants.Role `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
