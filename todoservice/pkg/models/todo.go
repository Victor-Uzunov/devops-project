package models

import (
	"encoding/json"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"time"
)

type Todo struct {
	ID          string                  `json:"id"`
	ListID      string                  `json:"list_id"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Tags        json.RawMessage         `json:"tags"`
	Completed   bool                    `json:"completed"`
	DueDate     *time.Time              `json:"due_date"`
	StartDate   *time.Time              `json:"start_date"`
	Priority    constants.PriorityLevel `json:"priority"`
	CreatedAt   time.Time               `json:"creation_date"`
	UpdatedAt   time.Time               `json:"last_update_date"`
	AssignedTo  *string                 `json:"assigned_to"`
}
