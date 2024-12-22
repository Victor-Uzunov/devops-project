package todos_test

import (
	"database/sql"
	"encoding/json"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/todos"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg"
	constants2 "github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func TestConvertTodoToModel(t *testing.T) {
	converter := todos.NewConverter()
	creationTime := time.Now()
	updateTime := time.Now()
	dueTime := time.Now()
	nullTime := sql.NullTime{
		Time:  dueTime,
		Valid: true,
	}
	assignedTo := "someone"

	tags, err := json.Marshal([]constants2.TagType{constants2.TagWork, constants2.TagFinance})
	require.NoError(t, err)

	entity := todos.Entity{
		ID:          "1",
		Title:       "Test Todo",
		Description: "This is a test todo",
		ListID:      "listID",
		Tags:        pkg.NewValidNullableString(string(tags)),
		Completed:   false,
		DueDate:     nullTime,
		Priority:    constants2.PriorityHigh,
		CreatedAt:   creationTime,
		UpdatedAt:   updateTime,
		AssignedTo:  &assignedTo,
	}

	expectedModel := models.Todo{
		ID:          "1",
		Title:       "Test Todo",
		Description: "This is a test todo",
		ListID:      "listID",
		Tags:        json.RawMessage(string(tags)),
		Completed:   false,
		DueDate:     &dueTime,
		Priority:    constants2.PriorityHigh,
		CreatedAt:   creationTime,
		UpdatedAt:   updateTime,
		AssignedTo:  &assignedTo,
	}

	got := converter.ConvertTodoToModel(entity)

	if !reflect.DeepEqual(got, expectedModel) {
		t.Errorf("ConvertTodoToModel() = %v, want %v", got, expectedModel)
	}
}

func TestConvertTodoToEntity(t *testing.T) {
	converter := todos.NewConverter()
	creationTime := time.Now()
	updateTime := time.Now()
	dueTime := time.Now()
	nullTime := sql.NullTime{
		Time:  dueTime,
		Valid: true,
	}
	assignedTo := "someone"

	tags, err := json.Marshal([]constants2.TagType{constants2.TagWork, constants2.TagFinance})
	require.NoError(t, err)

	model := models.Todo{
		ID:          "1",
		Title:       "Test Todo",
		Description: "This is a test todo",
		ListID:      "listID",
		Tags:        json.RawMessage(string(tags)),
		Completed:   false,
		DueDate:     &dueTime,
		StartDate:   &dueTime,
		Priority:    constants2.PriorityHigh,
		CreatedAt:   creationTime,
		UpdatedAt:   updateTime,
		AssignedTo:  &assignedTo,
	}

	expectedEntity := todos.Entity{
		ID:          "1",
		Title:       "Test Todo",
		Description: "This is a test todo",
		ListID:      "listID",
		Tags:        pkg.NewValidNullableString(string(tags)),
		Completed:   false,
		DueDate:     nullTime,
		StartDate:   nullTime,
		Priority:    constants2.PriorityHigh,
		CreatedAt:   creationTime,
		UpdatedAt:   updateTime,
		AssignedTo:  &assignedTo,
	}

	got := converter.ConvertTodoToEntity(model)

	if !reflect.DeepEqual(got, expectedEntity) {
		t.Errorf("ConvertTodoToEntity() = %v, want %v", got, expectedEntity)
	}
}
