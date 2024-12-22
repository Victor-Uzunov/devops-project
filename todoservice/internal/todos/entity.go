package todos

import (
	"database/sql"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"time"
)

type Entity struct {
	ID          string                  `db:"id"`
	ListID      string                  `db:"list_id"`
	Title       string                  `db:"title"`
	Description string                  `db:"description"`
	Tags        sql.NullString          `db:"tags"`
	Completed   bool                    `db:"completed"`
	DueDate     sql.NullTime            `db:"due_date"`
	StartDate   sql.NullTime            `db:"start_date"`
	Priority    constants.PriorityLevel `db:"priority"`
	CreatedAt   time.Time               `db:"created_at"`
	UpdatedAt   time.Time               `db:"updated_at"`
	AssignedTo  *string                 `db:"assigned_to"`
}
