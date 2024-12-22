package todos

import (
	"database/sql"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"time"
)

type Converter struct{}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) ConvertTodoToModel(entity Entity) models.Todo {
	return models.Todo{
		ID:          entity.ID,
		ListID:      entity.ListID,
		Title:       entity.Title,
		Description: entity.Description,
		Tags:        pkg.JSONRawMessageFromNullableString(entity.Tags),
		Completed:   entity.Completed,
		DueDate:     convertNullTimeToTime(entity.DueDate),
		StartDate:   convertNullTimeToTime(entity.StartDate),
		Priority:    entity.Priority,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		AssignedTo:  entity.AssignedTo,
	}
}

func (c *Converter) ConvertTodoToEntity(todo models.Todo) Entity {
	return Entity{
		ID:          todo.ID,
		ListID:      todo.ListID,
		Title:       todo.Title,
		Description: todo.Description,
		Tags:        pkg.NewNullableStringFromJSONRawMessage(todo.Tags),
		Completed:   todo.Completed,
		DueDate:     convertTimeToNullTime(*todo.DueDate),
		StartDate:   convertTimeToNullTime(*todo.StartDate),
		Priority:    todo.Priority,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
		AssignedTo:  todo.AssignedTo,
	}
}

func convertNullTimeToTime(nullTime sql.NullTime) *time.Time {
	if nullTime.Valid {
		return &nullTime.Time
	} else {
		return nil
	}
}

func convertTimeToNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}
