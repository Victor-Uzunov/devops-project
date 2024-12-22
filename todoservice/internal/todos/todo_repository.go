package todos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/db"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
)

//go:generate mockery --name=TodoRepository --output=automock --with-expecter=true --outpkg=automock --case=underscore --disable-version-string
type TodoRepository interface {
	Update(ctx context.Context, todo models.Todo) error
	Get(ctx context.Context, id string) (models.Todo, error)
	GetAll(ctx context.Context) ([]models.Todo, error)
	GetAllByListID(ctx context.Context, listID string) ([]models.Todo, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, list models.Todo) (string, error)
	CompleteTodo(ctx context.Context, id string) (models.Todo, error)
	UpdateTodoTitle(ctx context.Context, id string, title string) (models.Todo, error)
	UpdateTodoDescription(ctx context.Context, id string, description string) (models.Todo, error)
	UpdateTodoPriority(ctx context.Context, id string, priority constants.PriorityLevel) (models.Todo, error)
	UpdateAssignedTo(ctx context.Context, id string, userID string) (models.Todo, error)
}

type SQLXTodoRepository struct {
	converter *Converter
}

var _ TodoRepository = &SQLXTodoRepository{}

func NewSQLXTodoRepository() TodoRepository {
	return &SQLXTodoRepository{converter: NewConverter()}
}

func (r *SQLXTodoRepository) Create(ctx context.Context, todo models.Todo) (string, error) {
	log.C(ctx).Info("creating todo")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return "", err
	}

	entity := r.converter.ConvertTodoToEntity(todo)
	log.C(ctx).Debugf("converted entity in repo todo: %v", entity)

	if err = pkg.ValidateUUID(entity.ID); err != nil {
		log.C(ctx).Errorf("invalid entity id: %v", err)
		return "", fmt.Errorf("invalid ID format: %w", err)
	}
	if err = pkg.ValidateUUID(entity.ListID); err != nil {
		log.C(ctx).Errorf("invalid entity list id: %v", err)
		return "", fmt.Errorf("invalid list_id format: %w", err)
	}

	insertTodoQuery := `
		INSERT INTO todos (id, title, description, list_id, completed, tags, priority, due_date, start_date, assigned_to, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`

	var id string
	err = tx.QueryRowContext(ctx, insertTodoQuery,
		entity.ID,
		entity.Title,
		entity.Description,
		entity.ListID,
		entity.Completed,
		entity.Tags,
		entity.Priority,
		entity.DueDate,
		entity.StartDate,
		pkg.NullIfEmpty(*entity.AssignedTo),
		entity.CreatedAt,
		entity.UpdatedAt,
	).Scan(&id)
	log.C(ctx).Debugf("created todo with ID: %v", id)
	if err != nil {
		log.C(ctx).Errorf("failed to insert todo: %v", err)
		return "", fmt.Errorf("failed to create todo: %w", err)
	}
	return id, nil
}

func (r *SQLXTodoRepository) Get(ctx context.Context, id string) (models.Todo, error) {
	log.C(ctx).Info("getting todo")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return models.Todo{}, err
	}

	query := `
		SELECT id, title, description, list_id, priority, due_date, start_date, completed, tags, created_at, updated_at, assigned_to
		FROM todos
		WHERE id = $1
`
	var entity Entity
	err = tx.GetContext(ctx, &entity, query, id)
	if err != nil {
		log.C(ctx).Errorf("failed to get todo: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return models.Todo{}, fmt.Errorf("todo not found: %w", err)
		}
		return models.Todo{}, fmt.Errorf("failed to get todo: %w", err)
	}

	log.C(ctx).Debugf("got entity in repo layer: %v", entity)

	return r.converter.ConvertTodoToModel(entity), nil
}

func (r *SQLXTodoRepository) Update(ctx context.Context, todo models.Todo) error {
	log.C(ctx).Info("updating todo")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return err
	}
	entity := r.converter.ConvertTodoToEntity(todo)
	log.C(ctx).Debugf("converted entity in repo todo: %v", entity)

	updateTodoQuery := `
		UPDATE todos
		SET title = $1, description = $2, 
		    priority = $3, due_date = $4, start_date = $5, completed = $6, tags = $7, assigned_to = $8
		WHERE id = $9
	`

	_, err = tx.ExecContext(ctx, updateTodoQuery,
		entity.Title,
		entity.Description,
		entity.Priority,
		entity.DueDate,
		entity.StartDate,
		entity.Completed,
		entity.Tags,
		pkg.NullIfEmpty(*entity.AssignedTo),
		entity.ID,
	)
	if err != nil {
		log.C(ctx).Errorf("failed to update todo: %v", err)
		return fmt.Errorf("failed to update todo: %w", err)
	}
	return nil
}

func (r *SQLXTodoRepository) Delete(ctx context.Context, id string) error {
	log.C(ctx).Info("deleting todo")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return err
	}
	deleteQuery := `DELETE FROM todos WHERE id = $1`
	_, err = tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		log.C(ctx).Errorf("failed to delete todo: %v", err)
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	log.C(ctx).Debugf("deleted todo with ID: %v", id)
	return nil
}

func (r *SQLXTodoRepository) GetAllByListID(ctx context.Context, listID string) ([]models.Todo, error) {
	log.C(ctx).Info("getting all todos")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return []models.Todo{}, err
	}
	query := `
		SELECT id, title, description, list_id, priority, due_date, start_date, completed, tags, created_at, updated_at
		FROM todos
		WHERE list_id = $1
	`

	var todos []Entity

	err = tx.SelectContext(ctx, &todos, query, listID)
	if err != nil {
		log.C(ctx).Errorf("failed to get todos: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no todos found for list with id %s", listID)
		}
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}

	result := make([]models.Todo, 0)
	for _, entity := range todos {
		log.C(ctx).Debugf("got entity in repo layer: %v", entity)
		result = append(result, r.converter.ConvertTodoToModel(entity))
	}

	return result, nil
}

func (r *SQLXTodoRepository) GetAll(ctx context.Context) ([]models.Todo, error) {
	log.C(ctx).Info("getting all todos")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return []models.Todo{}, err
	}
	query := `
		SELECT id, title, description, list_id, priority, due_date, start_date, completed, tags, created_at, updated_at
		FROM todos
	`

	var todos []Entity

	err = tx.SelectContext(ctx, &todos, query)
	if err != nil {
		log.C(ctx).Errorf("failed to get todos: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no todos found")
		}
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}

	result := make([]models.Todo, 0)
	for _, entity := range todos {
		log.C(ctx).Debugf("got entity in repo layer: %v", entity)
		result = append(result, r.converter.ConvertTodoToModel(entity))
	}

	return result, nil
}

func (r *SQLXTodoRepository) CompleteTodo(ctx context.Context, id string) (models.Todo, error) {
	log.C(ctx).Info("completing todo repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return models.Todo{}, err
	}

	updateListQuery := `
		UPDATE todos
		SET completed = true
		WHERE id = $1
	`

	_, err = tx.ExecContext(ctx, updateListQuery, id)
	if err != nil {
		log.C(ctx).Errorf("failed to complete todo: %v", err)
		return models.Todo{}, err
	}

	updatedTodo, err := r.Get(ctx, id)
	if err != nil {
		log.C(ctx).Errorf("error while fetching completed todo: %v", err)
		return models.Todo{}, err
	}

	log.C(ctx).Info("todo completed successfully")
	return updatedTodo, nil
}

func (r *SQLXTodoRepository) UpdateTodoDescription(ctx context.Context, todoID string, description string) (models.Todo, error) {
	log.C(ctx).Info("updating todo description repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return models.Todo{}, err
	}

	updateTodoQuery := `
		UPDATE todos
		SET description = $1
		WHERE id = $2
	`

	_, err = tx.ExecContext(ctx, updateTodoQuery, description, todoID)
	if err != nil {
		log.C(ctx).Errorf("failed to update todo description: %v", err)
		return models.Todo{}, err
	}

	updatedTodo, err := r.Get(ctx, todoID)
	if err != nil {
		log.C(ctx).Errorf("error while fetching updated todo: %v", err)
		return models.Todo{}, err
	}

	log.C(ctx).Info("todo description updated successfully")
	return updatedTodo, nil
}

func (r *SQLXTodoRepository) UpdateTodoTitle(ctx context.Context, todoID string, title string) (models.Todo, error) {
	log.C(ctx).Info("updating todo title repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return models.Todo{}, err
	}

	updateTodoQuery := `
		UPDATE todos
		SET title = $1
		WHERE id = $2
	`

	_, err = tx.ExecContext(ctx, updateTodoQuery, title, todoID)
	if err != nil {
		log.C(ctx).Errorf("failed to update todo title: %v", err)
		return models.Todo{}, err
	}

	updatedTodo, err := r.Get(ctx, todoID)
	if err != nil {
		log.C(ctx).Errorf("error while fetching updated todo: %v", err)
		return models.Todo{}, err
	}

	log.C(ctx).Info("todo title updated successfully")
	return updatedTodo, nil
}

func (r *SQLXTodoRepository) UpdateTodoPriority(ctx context.Context, todoID string, priority constants.PriorityLevel) (models.Todo, error) {
	log.C(ctx).Info("updating todo priority repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return models.Todo{}, err
	}

	updateTodoQuery := `
		UPDATE todos
		SET priority = $1
		WHERE id = $2
	`

	_, err = tx.ExecContext(ctx, updateTodoQuery, priority, todoID)
	if err != nil {
		log.C(ctx).Errorf("failed to update todo priority: %v", err)
		return models.Todo{}, err
	}

	updatedTodo, err := r.Get(ctx, todoID)
	if err != nil {
		log.C(ctx).Errorf("error while fetching updated todo: %v", err)
		return models.Todo{}, err
	}

	log.C(ctx).Info("todo priority updated successfully")
	return updatedTodo, nil
}

func (r *SQLXTodoRepository) UpdateAssignedTo(ctx context.Context, todoID string, userID string) (models.Todo, error) {
	log.C(ctx).Info("updating todo assigned_to repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return models.Todo{}, err
	}

	updateTodoQuery := `
		UPDATE todos
		SET assigned_to = $1
		WHERE id = $2
	`

	_, err = tx.ExecContext(ctx, updateTodoQuery, userID, todoID)
	if err != nil {
		log.C(ctx).Errorf("failed to update todo assigned_to: %v", err)
		return models.Todo{}, err
	}

	updatedTodo, err := r.Get(ctx, todoID)
	if err != nil {
		log.C(ctx).Errorf("error while fetching updated todo: %v", err)
		return models.Todo{}, err
	}

	log.C(ctx).Info("todo assigned_to updated successfully")
	return updatedTodo, nil
}
