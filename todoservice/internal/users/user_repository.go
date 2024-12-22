package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/db"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"time"
)

//go:generate mockery --name=UserRepository --output=automock --with-expecter=true --outpkg=automock --case=underscore --disable-version-string
type UserRepository interface {
	Update(ctx context.Context, user models.User) error
	Get(ctx context.Context, id string) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, user models.User) (string, error)
	SaveRefreshToken(ctx context.Context, email string, refreshToken string, expirationTime time.Time) error
	FindByRefreshToken(ctx context.Context, refreshToken string) (models.User, error)
	Logout(ctx context.Context, email string) error
}

type SQLXUserRepository struct {
	converter *Converter
}

var _ UserRepository = &SQLXUserRepository{}

func NewSQLXUserRepository() UserRepository {
	return &SQLXUserRepository{converter: NewConverter()}
}

func (r *SQLXUserRepository) Create(ctx context.Context, user models.User) (string, error) {
	log.C(ctx).Infof("creating user in repo user: %+v", user)
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %+v", err)
		return "", err
	}

	entity := r.converter.ConvertUserToEntity(user)
	log.C(ctx).Debugf("converting user to entity in repo user: %+v", entity)

	insertUserQuery := `
		INSERT INTO users (id, email, github_id, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id string
	err = tx.QueryRowxContext(ctx, insertUserQuery,
		entity.ID,
		entity.Email,
		entity.GithubID,
		entity.Role,
		entity.CreatedAt,
		entity.UpdatedAt,
	).Scan(&id)
	if err != nil {
		log.C(ctx).Errorf("error in creating user in repo user: %+v", err)
		return "", fmt.Errorf("failed to create user: %w", err)
	}
	log.C(ctx).Debugf("created user in repo: %+v", id)
	return id, nil
}

func (r *SQLXUserRepository) Get(ctx context.Context, id string) (models.User, error) {
	log.C(ctx).Infof("getting user from repo user: %+v", id)
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %+v", err)
		return models.User{}, err
	}

	query := `
		SELECT id, email, github_id, role, created_at, updated_at
		FROM users
		WHERE id = $1
`
	var entity Entity
	err = tx.GetContext(ctx, &entity, query, id)
	if err != nil {
		log.C(ctx).Errorf("error in getting user from repo: %+v and id: %s", err, id)
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("user not found: %w", err)
		}
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	log.C(ctx).Debugf("got user from repo with id: %+v", entity.ID)

	return r.converter.ConvertUserToModel(entity), nil
}

func (r *SQLXUserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	log.C(ctx).Infof("getting user from repo user with email: %+v", email)
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %+v", err)
		return models.User{}, err
	}

	query := `
		SELECT id, email, github_id, role, created_at, updated_at
		FROM users
		WHERE email = $1
`
	var entity Entity
	err = tx.GetContext(ctx, &entity, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("user not found: %w", err)
		}
		log.C(ctx).Errorf("error in getting user from repo: %+v", err)
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	log.C(ctx).Debugf("got user from repo with id: %+v", entity.ID)

	return r.converter.ConvertUserToModel(entity), nil
}

func (r *SQLXUserRepository) Update(ctx context.Context, user models.User) error {
	log.C(ctx).Infof("updating user in repo user: %+v", user)
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %+v", err)
		return err
	}
	entity := r.converter.ConvertUserToEntity(user)

	updateUserQuery := `
		UPDATE users
		SET email = :email, github_id = :github_id, role = :role
		WHERE id = :id
	`

	_, err = tx.NamedExecContext(ctx, updateUserQuery, entity)
	if err != nil {
		log.C(ctx).Errorf("error in updating user in repo user: %+v", err)
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *SQLXUserRepository) Delete(ctx context.Context, id string) error {
	log.C(ctx).Infof("deleting user from repo user: %+v", id)
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %+v", err)
		return err
	}
	deleteQuery := `DELETE FROM users WHERE id = $1`
	_, err = tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		log.C(ctx).Errorf("error in deleting user in repo user: %+v", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *SQLXUserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	log.C(ctx).Infof("getting all users in repo users")
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %+v", err)
		return []models.User{}, err
	}
	query := `
		SELECT id, id, email, github_id, role, created_at, updated_at
		FROM users
	`

	var users []Entity

	err = tx.SelectContext(ctx, &users, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no todos found")
		}
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}

	result := make([]models.User, 0)
	for _, entity := range users {
		log.C(ctx).Debugf("got user: %+v", entity)
		result = append(result, r.converter.ConvertUserToModel(entity))
	}

	return result, nil
}

func (r *SQLXUserRepository) SaveRefreshToken(ctx context.Context, email string, refreshToken string, expirationTime time.Time) error {
	log.C(ctx).Infof("saving refresh token for user: %+v", email)
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %+v", err)
		return err
	}

	query := `
		UPDATE users 
		SET refresh_token = $1, refresh_token_expiration = $2 
		WHERE email = $3
	`

	_, err = tx.ExecContext(ctx, query, refreshToken, expirationTime, email)
	if err != nil {
		log.C(ctx).Errorf("failed to save refresh token: %+v", err)
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	log.C(ctx).Info("refresh token saved successfully")
	return nil
}

func (r *SQLXUserRepository) FindByRefreshToken(ctx context.Context, refreshToken string) (models.User, error) {
	log.C(ctx).Infof("finding user by refresh token: %+v", refreshToken)
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %+v", err)
		return models.User{}, err
	}

	query := `
        SELECT id, email, github_id, role, created_at, updated_at
        FROM users
        WHERE refresh_token = $1
    `

	var entity Entity
	err = tx.GetContext(ctx, &entity, query, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("no user found for refresh token")
		}
		return models.User{}, fmt.Errorf("failed to find user by refresh token: %w", err)
	}

	return r.converter.ConvertUserToModel(entity), nil
}

func (r *SQLXUserRepository) Logout(ctx context.Context, email string) error {
	log.C(ctx).Infof("logouting user in repo user: %+v", email)
	tx, err := db.FromContext(ctx)
	if err != nil {
		log.C(ctx).Errorf("error while getting transaction from context: %+v", err)
		return err
	}
	updateQuery := `
		UPDATE users 
		SET refresh_token = NULL, refresh_token_expiration = NULL 
		WHERE email = $1`

	_, err = tx.ExecContext(ctx, updateQuery, email)
	if err != nil {
		log.C(ctx).Errorf("error in updating user refresh token in repo user: %+v", err)
		return fmt.Errorf("failed to clear refresh token: %w", err)
	}
	return nil
}
