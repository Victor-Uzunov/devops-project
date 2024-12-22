package users

import (
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"time"
)

type Entity struct {
	ID        string         `db:"id"`
	Email     string         `db:"email"`
	GithubID  string         `db:"github_id"`
	Role      constants.Role `db:"role"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}
