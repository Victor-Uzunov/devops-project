package models

import "github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"

type Access struct {
	ListID string         `json:"list_id"`
	UserID string         `json:"user_id"`
	Role   constants.Role `json:"role"`
	Status string         `json:"status"`
}
