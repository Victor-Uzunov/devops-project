package converters

import (
	"encoding/json"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/generated/graphql"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	format "github.com/Victor-Uzunov/devops-project/todoservice/pkg/time"
	"time"
)

//go:generate mockery --name=TodoConverter --output=automock --with-expecter=true --outpkg=automock --case=underscore --disable-version-string
type TodoConverter interface {
	ConvertTodoToGraphQL(todo models.Todo) (*graphql.Todo, error)
	ConvertCreateTodoInput(input graphql.CreateTodoInput) (models.Todo, error)
	ConvertUpdateTodoInput(input graphql.UpdateTodoInput) (models.Todo, error)
	ConvertMultipleTodoToGraphQL(todos []*models.Todo) ([]*graphql.Todo, error)
}

type ConverterTodoGraphQL struct{}

func NewConverterTodoGraphQL() TodoConverter {
	return &ConverterTodoGraphQL{}
}

func (c *ConverterTodoGraphQL) ConvertTodoToGraphQL(todo models.Todo) (*graphql.Todo, error) {
	priority, err := convertPriorityToGraphQL(todo.Priority)
	if err != nil {
		return &graphql.Todo{}, fmt.Errorf("convertPriorityToGraphQL: %w", err)
	}
	var tags []string

	if todo.Tags != nil {
		err := json.Unmarshal(todo.Tags, &tags)
		if err != nil {
			fmt.Println("Error unmarshalling tags:", err)
			tags = make([]string, 0)
		}
	}

	return &graphql.Todo{
		ID:          todo.ID,
		List:        nil,
		Title:       todo.Title,
		Completed:   todo.Completed,
		Description: &todo.Description,
		Tags:        tags,
		Priority:    &priority,
		DueDate:     format.TimeToString(todo.DueDate),
		StartDate:   format.TimeToString(todo.StartDate),
		CreatedAt:   todo.CreatedAt.Format(constants.DateFormat),
		UpdatedAt:   todo.UpdatedAt.Format(constants.DateFormat),
		AssignedTo:  nil,
	}, nil
}

func convertPriorityToGraphQL(priority constants.PriorityLevel) (graphql.Priority, error) {
	switch priority {
	case constants.PriorityLow:
		return graphql.PriorityLow, nil
	case constants.PriorityMedium:
		return graphql.PriorityMedium, nil
	case constants.PriorityHigh:
		return graphql.PriorityHigh, nil
	default:
		return graphql.Priority(constants.PriorityLow), fmt.Errorf("invalid priority level: %v", priority)
	}
}

func (c *ConverterTodoGraphQL) ConvertCreateTodoInput(input graphql.CreateTodoInput) (models.Todo, error) {
	priority, err := ConvertPriorityFromGraphQL(*input.Priority)
	if err != nil {
		return models.Todo{}, fmt.Errorf("convertCreateTodoInput: %w", err)
	}
	var dueDate time.Time
	var startDate time.Time

	if *input.DueDate != "null" {
		dueDate, err = time.Parse(time.RFC3339, *input.DueDate)
		if err != nil {
			return models.Todo{}, fmt.Errorf("convertCreateTodoInput: %w", err)
		}
	}

	if *input.StartDate != "null" {
		startDate, err = time.Parse(time.RFC3339, *input.StartDate)
		if err != nil {
			return models.Todo{}, fmt.Errorf("convertCreateTodoInput: %w", err)
		}
	}

	jsonBytes, err := json.Marshal(input.Tags)
	if err != nil {
		return models.Todo{}, fmt.Errorf("convertCreateTodoInput: %w", err)
	}
	jsonRawMessage := json.RawMessage(jsonBytes)
	return models.Todo{
		Title:       input.Title,
		ListID:      input.ListID,
		Description: *input.Description,
		DueDate:     &dueDate,
		StartDate:   &startDate,
		Priority:    priority,
		Tags:        jsonRawMessage,
		Completed:   *input.Completed,
		AssignedTo:  input.AssignedTo,
	}, nil
}

func (c *ConverterTodoGraphQL) ConvertUpdateTodoInput(input graphql.UpdateTodoInput) (models.Todo, error) {
	priority, err := ConvertPriorityFromGraphQL(*input.Priority)
	if err != nil {
		return models.Todo{}, fmt.Errorf("convertCreateTodoInput: %w", err)
	}
	var dueDate time.Time
	var startDate time.Time

	if *input.DueDate != "null" {
		dueDate, err = time.Parse(time.RFC3339, *input.DueDate)
		if err != nil {
			return models.Todo{}, fmt.Errorf("cannot convert dueDate in graphql converter: %w", err)
		}
	}

	if *input.StartDate != "null" {
		startDate, err = time.Parse(time.RFC3339, *input.StartDate)
		if err != nil {
			return models.Todo{}, fmt.Errorf("cannot convert startDate in graphql converter: %w", err)
		}
	}

	jsonBytes, err := json.Marshal(input.Tags)
	if err != nil {
		return models.Todo{}, fmt.Errorf("convertCreateTodoInput: %w", err)
	}
	jsonRawMessage := json.RawMessage(jsonBytes)
	return models.Todo{
		Title:       *input.Title,
		Description: *input.Description,
		DueDate:     &dueDate,
		StartDate:   &startDate,
		Priority:    priority,
		Tags:        jsonRawMessage,
		AssignedTo:  input.AssignedTo,
	}, nil
}

func (c *ConverterTodoGraphQL) ConvertMultipleTodoToGraphQL(todos []*models.Todo) ([]*graphql.Todo, error) {
	result := make([]*graphql.Todo, 0)
	for _, el := range todos {
		t, err := c.ConvertTodoToGraphQL(*el)
		if err != nil {
			return nil, fmt.Errorf("error converting todo: %w", err)
		}
		result = append(result, t)
	}
	return result, nil
}
