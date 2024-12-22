package lists

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/db"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/todos"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/jmoiron/sqlx"
)

//go:generate mockery --name=ListRepository --output=automock --with-expecter=true --outpkg=automock --case=underscore --disable-version-string
type ListRepository interface {
	Update(ctx context.Context, list models.List) error
	Get(ctx context.Context, id string) (models.List, error)
	GetAll(ctx context.Context) ([]models.List, error)
	GetAccess(ctx context.Context, listID string, userID string) (models.Access, error)
	GetUsersByListID(ctx context.Context, listID string) ([]models.Access, error)
	GetListOwnerID(ctx context.Context, listID string) (string, error)
	ListAllByUserID(ctx context.Context, userID string) ([]models.Access, error)
	Delete(ctx context.Context, id string) error
	DeleteAccess(ctx context.Context, listID string, userID string) error
	Create(ctx context.Context, list models.List) (string, error)
	CreateAccess(ctx context.Context, access models.Access) (models.Access, error)
	UpdateListDescription(ctx context.Context, listID string, description string) (models.List, error)
	UpdateListName(ctx context.Context, listID string, name string) (models.List, error)
	GetAllTodosForList(ctx context.Context, listID string) ([]models.Todo, error)
	GetPendingLists(ctx context.Context, userID string) ([]models.Access, error)
	AcceptList(ctx context.Context, listID string, userID string) error
	GetAccessesByListID(ctx context.Context, listID string) ([]models.Access, error)
	GetAcceptedLists(ctx context.Context, listID string) ([]models.Access, error)
}

type SQLXListRepository struct {
	converter *Converter
}

var _ ListRepository = &SQLXListRepository{}

func NewSQLXListRepository() ListRepository {
	return &SQLXListRepository{converter: NewConverter()}
}

func (r *SQLXListRepository) Create(ctx context.Context, list models.List) (string, error) {
	log.C(ctx).Info("creating list repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return "", err
	}

	entity := r.converter.ConvertListToEntity(list)
	log.C(ctx).Debugf("successfull converted list: %v", entity)

	insertListQuery := `
		INSERT INTO lists (id, name, description, owner_id, visibility, tags, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var id string
	err = tx.QueryRowContext(ctx, insertListQuery,
		entity.ID,
		entity.Name,
		entity.Description,
		entity.OwnerID,
		entity.Visibility,
		entity.Tags,
		entity.CreatedAt,
		entity.UpdatedAt,
	).Scan(&id)
	log.C(ctx).Debugf("list id: %v", id)
	if err != nil {
		log.C(ctx).Errorf("failed to create list: %v", err)
		return "", fmt.Errorf("failed to create list: %w", err)
	}

	var role string
	status := constants.StatusOwner

	insertAccessQuery := `
			INSERT INTO list_access (user_id, list_id, access_level, status)
			VALUES ($1, $2, $3, $4)
		`
	if _, err = tx.ExecContext(ctx, insertAccessQuery, entity.OwnerID, entity.ID, "admin", status); err != nil {
		log.C(ctx).Errorf("failed to create list_access: %v", err)
		return "", fmt.Errorf("failed to create list_access: %w", err)
	}

	status = constants.StatusPending
	for _, userID := range entity.SharedWith {
		getRoleQuery := `
			SELECT role 
			FROM users 
			WHERE id = $1
		`
		err = tx.GetContext(ctx, &role, getRoleQuery, userID)
		if err != nil {
			log.C(ctx).Errorf("failed to create list_access: %v", err)
			return "", fmt.Errorf("failed to get user's role: %w", err)
		}
		if _, err = tx.ExecContext(ctx, insertAccessQuery, userID, entity.ID, role, status); err != nil {
			log.C(ctx).Errorf("failed to create list_access: %v", err)
			return "", fmt.Errorf("failed to insert access: %w", err)
		}
	}

	return id, nil
}

func (r *SQLXListRepository) Get(ctx context.Context, id string) (models.List, error) {
	log.C(ctx).Info("getting list repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return models.List{}, err
	}

	query := `
		SELECT id, name, description, owner_id, visibility, tags, created_at, updated_at
		FROM lists
		WHERE id = $1
`
	var entity Entity
	err = tx.GetContext(ctx, &entity, query, id)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch list: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return models.List{}, fmt.Errorf("list not found: %w", err)
		}
		return models.List{}, fmt.Errorf("failed to get list: %w", err)
	}

	log.C(ctx).Debugf("list: %v", entity)

	sharedWithQuery := `
		SELECT user_id
		FROM list_access
		WHERE list_id = $1
	`
	var sharedWith []string
	err = tx.SelectContext(ctx, &sharedWith, sharedWithQuery, id)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch list_access: %v", err)
		return models.List{}, fmt.Errorf("failed to fetch list_access: %w", err)
	}

	entity.SharedWith = sharedWith
	return r.converter.ConvertListToModel(entity), nil
}

func (r *SQLXListRepository) Update(ctx context.Context, list models.List) error {
	log.C(ctx).Info("updating list repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return err
	}
	entity := r.converter.ConvertListToEntity(list)

	updateListQuery := `
		UPDATE lists
		SET name = $1, description = $2, visibility = $3, tags = $4
		WHERE id = $5
	`

	_, err = tx.ExecContext(ctx, updateListQuery, entity.Name, entity.Description,
		entity.Visibility, entity.Tags, entity.ID)
	if err != nil {
		log.C(ctx).Errorf("failed to update list: %v", err)
		return fmt.Errorf("failed to update list: %w", err)
	}

	if entity.Visibility == constants.VisibilityPrivate {
		err = r.privateList(ctx, tx, list.ID)
		if err != nil {
			log.C(ctx).Errorf("failed to remove access for a private list: %v", err)
			return err
		}
		log.C(ctx).Debugf("removed access for a private list: %s", list.ID)
	}

	log.C(ctx).Debugf("successfully updated list: %v", entity)
	return nil
}

func (r *SQLXListRepository) privateList(ctx context.Context, tx *sqlx.Tx, listID string) error {
	log.C(ctx).Info("remove accesses for updated private list repository")
	deleteQuery := `DELETE FROM list_access WHERE list_id = $1 AND status IN ($2, $3)`
	_, err := tx.ExecContext(ctx, deleteQuery, listID, constants.StatusAccepted, constants.StatusPending)
	if err != nil {
		log.C(ctx).Errorf("failed to delete list_access for private list: %v", err)
		return fmt.Errorf("failed to delete list for private list: %w", err)
	}
	return nil
}

func (r *SQLXListRepository) Delete(ctx context.Context, id string) error {
	log.C(ctx).Info("deleting list repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return err
	}
	deleteQuery := `DELETE FROM lists WHERE id = $1`
	_, err = tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		log.C(ctx).Errorf("failed to delete lists: %v", err)
		return fmt.Errorf("failed to delete list: %w", err)
	}
	log.C(ctx).Debugf("deleted list: %v", id)
	return nil
}

func (r *SQLXListRepository) ListAllByUserID(ctx context.Context, userID string) ([]models.Access, error) {
	log.C(ctx).Info("listing all lists by user repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return []models.Access{}, err
	}
	query := `
		SELECT list_id, user_id, access_level, status
		FROM list_access
		WHERE user_id = $1 AND status = 'owner'
	`

	var lists []AccessEntity

	err = tx.SelectContext(ctx, &lists, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.C(ctx).Errorf("failed to list all lists by user repository: %v", err)
			return nil, fmt.Errorf("no lists found for user with id %s", userID)
		}
		return nil, fmt.Errorf("failed to get lists: %w", err)
	}

	result := make([]models.Access, 0)
	for _, entity := range lists {
		result = append(result, r.converter.ConvertAccessToModel(entity))
	}
	log.C(ctx).Debugf("list all by user ID repository: %v", result)
	return result, nil
}

func (r *SQLXListRepository) GetAcceptedLists(ctx context.Context, userID string) ([]models.Access, error) {
	log.C(ctx).Info("get accepted by user repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return []models.Access{}, err
	}
	query := `
		SELECT list_id, user_id, access_level, status
		FROM list_access
		WHERE user_id = $1 AND status = 'accepted'
	`

	var lists []AccessEntity

	err = tx.SelectContext(ctx, &lists, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.C(ctx).Errorf("failed to list all lists by user repository: %v", err)
			return nil, fmt.Errorf("no lists found for user with id %s", userID)
		}
		return nil, fmt.Errorf("failed to get lists: %w", err)
	}

	result := make([]models.Access, 0)
	for _, entity := range lists {
		result = append(result, r.converter.ConvertAccessToModel(entity))
	}
	log.C(ctx).Debugf("list all by user ID repository: %v", result)
	return result, nil
}

func (r *SQLXListRepository) GetPendingLists(ctx context.Context, userID string) ([]models.Access, error) {
	log.C(ctx).Info("listing all pending lists by user repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return []models.Access{}, err
	}
	query := `
		SELECT list_id, user_id, access_level, status
		FROM list_access
		WHERE user_id = $1 AND status = 'pending'
	`

	var lists []AccessEntity

	err = tx.SelectContext(ctx, &lists, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.C(ctx).Errorf("failed to list all pending lists by user repository: %v", err)
			return nil, fmt.Errorf("no lists found for user with id %s", userID)
		}
		return nil, fmt.Errorf("failed to get pending lists: %w", err)
	}

	result := make([]models.Access, 0)
	for _, entity := range lists {
		result = append(result, r.converter.ConvertAccessToModel(entity))
	}
	log.C(ctx).Debugf("list all by user ID repository: %v", result)
	return result, nil
}

func (r *SQLXListRepository) GetAll(ctx context.Context) ([]models.List, error) {
	log.C(ctx).Info("listing all lists repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return []models.List{}, err
	}
	query := `
		SELECT id, name, description, owner_id, visibility, tags, created_at, updated_at
		FROM lists
	`

	var lists []Entity

	err = tx.SelectContext(ctx, &lists, query)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch all lists: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no lists found")
		}
		return nil, fmt.Errorf("failed to get lists: %w", err)
	}

	log.C(ctx).Debugf("list all: %v", lists)

	sharedWithQuery := `
		SELECT user_id
		FROM list_access
		WHERE list_id = $1
	`

	result := make([]models.List, 0)
	for _, entity := range lists {
		var sharedWith []string
		err = tx.SelectContext(ctx, &sharedWith, sharedWithQuery, entity.ID)
		if err != nil {
			log.C(ctx).Errorf("failed to fetch lists: %v", err)
			return []models.List{}, fmt.Errorf("failed to fetch list_access: %w", err)
		}
		entity.SharedWith = sharedWith
		result = append(result, r.converter.ConvertListToModel(entity))
	}

	log.C(ctx).Debugf("list all: %v", result)

	return result, nil
}

func (r *SQLXListRepository) GetUsersByListID(ctx context.Context, listID string) ([]models.Access, error) {
	log.C(ctx).Info("listing all lists repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return []models.Access{}, err
	}
	query := `
		SELECT list_id, user_id, access_level, status
		FROM list_access
		WHERE list_id = $1
	`

	var access []AccessEntity

	err = tx.SelectContext(ctx, &access, query, listID)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch list_access: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no access found")
		}
		return nil, fmt.Errorf("failed to get access: %w", err)
	}

	log.C(ctx).Debugf("list access: %v", access)

	result := make([]models.Access, 0)
	for _, entity := range access {
		result = append(result, r.converter.ConvertAccessToModel(entity))
	}

	log.C(ctx).Debugf("list access: %v", result)

	return result, nil
}

func (r *SQLXListRepository) GetListOwnerID(ctx context.Context, listID string) (string, error) {
	log.C(ctx).Info("listing all lists repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return "", err
	}
	query := `
		SELECT owner_id
		FROM lists
		WHERE id = $1
	`

	var ownerID string

	err = tx.GetContext(ctx, &ownerID, query, listID)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch list_owner_id: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("no access found")
		}
		return "", fmt.Errorf("failed to get access: %w", err)
	}
	log.C(ctx).Debugf("list owner: %v", ownerID)
	return ownerID, nil
}

func (r *SQLXListRepository) CreateAccess(ctx context.Context, access models.Access) (models.Access, error) {
	log.C(ctx).Info("creating list repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return models.Access{}, err
	}

	entity := r.converter.ConvertAccessToEntity(access)
	status := constants.StatusPending

	insertListQuery := `
		INSERT INTO list_access (list_id, user_id, access_level, status)
		VALUES ($1, $2, $3, $4)
	`

	_, err = tx.ExecContext(ctx, insertListQuery,
		entity.ListID,
		entity.UserID,
		entity.Role,
		status,
	)
	if err != nil {
		log.C(ctx).Errorf("failed to create list_access: %v", err)
		return models.Access{}, fmt.Errorf("failed to create list access: %w", err)
	}

	return access, nil
}

func (r *SQLXListRepository) AcceptList(ctx context.Context, listID string, userID string) error {
	log.C(ctx).Info("accept list repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return err
	}

	status := constants.StatusAccepted

	updateListAccessQuery := `
		UPDATE list_access
		SET status = $1
		WHERE user_id = $2 AND list_id = $3
	`

	_, err = tx.ExecContext(ctx, updateListAccessQuery,
		status,
		userID,
		listID,
	)
	if err != nil {
		log.C(ctx).Errorf("failed to accept list_access: %v", err)
		return fmt.Errorf("failed to create list access: %w", err)
	}

	return nil
}

func (r *SQLXListRepository) DeleteAccess(ctx context.Context, listID string, userID string) error {
	log.C(ctx).Info("deleting list repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return err
	}
	deleteQuery := `DELETE FROM list_access WHERE list_id = $1 AND user_id = $2`
	_, err = tx.ExecContext(ctx, deleteQuery, listID, userID)
	if err != nil {
		log.C(ctx).Errorf("failed to delete list_access: %v", err)
		return fmt.Errorf("failed to delete list: %w", err)
	}
	return nil
}

func (r *SQLXListRepository) GetAccess(ctx context.Context, listID string, userID string) (models.Access, error) {
	log.C(ctx).Info("getting list repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return models.Access{}, err
	}

	query := `
		SELECT list_id, user_id, access_level, status
		FROM list_access
		WHERE list_id = $1 AND user_id = $2
`
	var entity AccessEntity
	err = tx.GetContext(ctx, &entity, query, listID, userID)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch list_access: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return models.Access{}, fmt.Errorf("list not found: %w", err)
		}
		return models.Access{}, fmt.Errorf("failed to get list: %w", err)
	}

	log.C(ctx).Debugf("list access: %v", entity)

	return r.converter.ConvertAccessToModel(entity), nil
}

func (r *SQLXListRepository) UpdateListDescription(ctx context.Context, listID string, description string) (models.List, error) {
	log.C(ctx).Info("updating list description repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return models.List{}, err
	}

	updateListQuery := `
		UPDATE lists
		SET description = $1
		WHERE id = $2
	`

	_, err = tx.ExecContext(ctx, updateListQuery, description, listID)
	if err != nil {
		log.C(ctx).Errorf("failed to update list description: %v", err)
		return models.List{}, err
	}

	updatedList, err := r.Get(ctx, listID)
	if err != nil {
		log.C(ctx).Errorf("error while fetching updated list: %v", err)
		return models.List{}, err
	}

	log.C(ctx).Info("list description updated successfully")
	return updatedList, nil
}

func (r *SQLXListRepository) UpdateListName(ctx context.Context, listID string, name string) (models.List, error) {
	log.C(ctx).Info("updating list name repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return models.List{}, err
	}

	updateListQuery := `
		UPDATE lists
		SET name = $1
		WHERE id = $2
	`

	_, err = tx.ExecContext(ctx, updateListQuery, name, listID)
	if err != nil {
		log.C(ctx).Errorf("failed to update list name: %v", err)
		return models.List{}, err
	}

	updatedList, err := r.Get(ctx, listID)
	if err != nil {
		log.C(ctx).Errorf("error while fetching updated list: %v", err)
		return models.List{}, err
	}

	log.C(ctx).Info("list name updated successfully")
	return updatedList, nil
}

func (r *SQLXListRepository) GetAllTodosForList(ctx context.Context, listID string) ([]models.Todo, error) {
	log.C(ctx).Info("getting all todos for a list")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return []models.Todo{}, err
	}
	query := `
		SELECT id, title, description, list_id, priority, start_date, due_date, completed, tags, created_at, updated_at
		FROM todos
		WHERE list_id = $1
	`

	var allTodos []todos.Entity

	err = tx.SelectContext(ctx, &allTodos, query, listID)
	if err != nil {
		log.C(ctx).Errorf("failed to get todos: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no todos found for list with id %s", listID)
		}
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}

	result := make([]models.Todo, 0)
	converter := todos.Converter{}
	for _, entity := range allTodos {
		log.C(ctx).Debugf("got entity in repo layer: %v", entity)
		result = append(result, converter.ConvertTodoToModel(entity))
	}

	return result, nil
}

func (r *SQLXListRepository) GetAccessesByListID(ctx context.Context, listID string) ([]models.Access, error) {
	log.C(ctx).Info("getting list accesses by list id repository")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %v", err)
		return []models.Access{}, err
	}

	query := `
		SELECT list_id, user_id, access_level, status
		FROM list_access
		WHERE list_id = $1
`
	var entities []AccessEntity
	err = tx.SelectContext(ctx, &entities, query, listID)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch list_access: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Access{}, fmt.Errorf("list not found: %w", err)
		}
		return []models.Access{}, fmt.Errorf("failed to get list: %w", err)
	}

	log.C(ctx).Debugf("list access: %v", entities)

	var accesses []models.Access
	for _, entity := range entities {
		accesses = append(accesses, r.converter.ConvertAccessToModel(entity))
	}

	return accesses, nil
}
