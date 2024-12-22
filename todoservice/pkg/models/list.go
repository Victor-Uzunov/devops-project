package models

import (
	"encoding/json"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"time"
)

type List struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	OwnerID     string               `json:"owner_id"`
	SharedWith  []string             `json:"shared_with"`
	Tags        json.RawMessage      `json:"tags"`
	CreatedAt   time.Time            `json:"creation_date"`
	UpdatedAt   time.Time            `json:"last_update_date"`
	Visibility  constants.Visibility `json:"visibility"`
}
