package lists

import (
	"database/sql"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"time"
)

type Entity struct {
	ID          string               `db:"id"`
	Name        string               `db:"name"`
	Description string               `db:"description"`
	OwnerID     string               `db:"owner_id"`
	SharedWith  []string             `db:"shared_with"`
	Tags        sql.NullString       `db:"tags"`
	CreatedAt   time.Time            `db:"created_at"`
	UpdatedAt   time.Time            `db:"updated_at"`
	Visibility  constants.Visibility `db:"visibility"`
}

type AccessEntity struct {
	ListID string         `db:"list_id"`
	UserID string         `db:"user_id"`
	Role   constants.Role `db:"access_level"`
	Status string         `db:"status"`
}
