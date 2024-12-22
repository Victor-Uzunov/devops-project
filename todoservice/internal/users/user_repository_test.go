package users_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/db"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/users"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"time"
)

func TestSQLXTUserRepositoryCreate(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := users.NewSQLXUserRepository()

	testCases := []struct {
		name          string
		input         models.User
		setupMocks    func()
		expectedID    string
		expectedError error
	}{
		{
			name: "Successful creation",
			input: models.User{
				ID:       "4cfd7e64-7431-4690-a2a0-1268917cedf3",
				Email:    "Test Todo",
				GithubID: "Test Description",
				Role:     "4cfd7e64-7431-4690-a2a0-1268917cedf3",
			},
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(`^INSERT INTO users`).WithArgs(
					"4cfd7e64-7431-4690-a2a0-1268917cedf3", "Test Todo", "Test Description", "4cfd7e64-7431-4690-a2a0-1268917cedf3", sqlxmock.AnyArg(), sqlxmock.AnyArg(),
				).WillReturnRows(sqlxmock.NewRows([]string{"id"}).AddRow("4cfd7e64-7431-4690-a2a0-1268917cedf3"))

				mockDB.ExpectCommit()
			},
			expectedID:    "4cfd7e64-7431-4690-a2a0-1268917cedf3",
			expectedError: nil,
		},
		{
			name: "Failed creation due to database error",
			input: models.User{
				ID:       "4cfd7e64-7431-4690-a2a0-1268917cedf3",
				Email:    "Test Todo",
				GithubID: "Test Description",
				Role:     "4cfd7e64-7431-4690-a2a0-1268917cedf3",
			},
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(`^INSERT INTO users`).WithArgs(
					"4cfd7e64-7431-4690-a2a0-1268917cedf3", "Test Todo", "Test Description", "4cfd7e64-7431-4690-a2a0-1268917cedf3", sqlxmock.AnyArg(), sqlxmock.AnyArg(),
				).WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			expectedID:    "",
			expectedError: fmt.Errorf("failed to create user: %w", errors.New("db error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			ctx := context.Background()
			tx, err := database.BeginTxx(ctx, nil)
			require.NoError(t, err)

			ctx = db.SaveToContext(ctx, tx)

			id, err := repo.Create(ctx, tc.input)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				err = tx.Rollback()
				require.NoError(t, err)

			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedID, id)
				err = tx.Commit()
				require.NoError(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestSQLXUserRepositoryGet(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := users.NewSQLXUserRepository()

	testCases := []struct {
		name          string
		id            string
		setupMocks    func()
		expectedList  models.User
		expectedError error
	}{
		{
			name: "Successful get of a user",
			id:   "4cfd7e64-7431-4690-a2a0-1268917cedf3",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT id, email, github_id, role, created_at, updated_at FROM users").WithArgs(
					"4cfd7e64-7431-4690-a2a0-1268917cedf3").WillReturnRows(sqlxmock.NewRows([]string{"id", "email", "github_id", "role", "created_at", "updated_at"}).
					AddRow("4cfd7e64-7431-4690-a2a0-1268917cedf3", "Test Todo", "Test Description", "4cfd7e64-7431-4690-a2a0-1268917cedf3", time.Time{}, time.Time{}))

				mockDB.ExpectCommit()
			},
			expectedList: models.User{
				ID:       "4cfd7e64-7431-4690-a2a0-1268917cedf3",
				Email:    "Test Todo",
				GithubID: "Test Description",
				Role:     "4cfd7e64-7431-4690-a2a0-1268917cedf3",
			},
			expectedError: nil,
		},
		{
			name: "Failed get todo due to database error",
			id:   "4cfd7e64-7431-4690-a2a0-1268917cedf3",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT id, email, github_id, role, created_at, updated_at FROM users").WithArgs(
					"4cfd7e64-7431-4690-a2a0-1268917cedf3").WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			expectedList:  models.User{},
			expectedError: fmt.Errorf("failed to get user: %w", errors.New("db error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			ctx := context.Background()
			tx, err := database.BeginTxx(ctx, nil)
			require.NoError(t, err)

			ctx = db.SaveToContext(ctx, tx)

			list, err := repo.Get(ctx, tc.id)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				err = tx.Rollback()
				require.NoError(t, err)

			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedList, list)
				err = tx.Commit()
				require.NoError(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestSQLXUserRepositoryUpdate(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := users.NewSQLXUserRepository()

	testCases := []struct {
		name          string
		id            string
		setupMocks    func()
		expectedError error
		input         models.User
	}{
		{
			name: "Successful get of a user",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectExec("UPDATE users").
					WithArgs("Test Todo", "Test Description", "4cfd7e64-7431-4690-a2a0-1268917cedf3", "4cfd7e64-7431-4690-a2a0-1268917cedf3").
					WillReturnResult(sqlxmock.NewResult(1, 1))

				mockDB.ExpectCommit()
			},
			input: models.User{
				ID:       "4cfd7e64-7431-4690-a2a0-1268917cedf3",
				Email:    "Test Todo",
				GithubID: "Test Description",
				Role:     "4cfd7e64-7431-4690-a2a0-1268917cedf3",
			},
			expectedError: nil,
		},
		{
			name: "Failed to update user due to database error",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectExec("UPDATE users").
					WithArgs("Test Todo", "Test Description", "4cfd7e64-7431-4690-a2a0-1268917cedf3", "4cfd7e64-7431-4690-a2a0-1268917cedf3").
					WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			input: models.User{
				ID:       "4cfd7e64-7431-4690-a2a0-1268917cedf3",
				Email:    "Test Todo",
				GithubID: "Test Description",
				Role:     "4cfd7e64-7431-4690-a2a0-1268917cedf3",
			},
			expectedError: fmt.Errorf("failed to update user: %w", errors.New("db error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			ctx := context.Background()
			tx, err := database.BeginTxx(ctx, nil)
			require.NoError(t, err)

			ctx = db.SaveToContext(ctx, tx)

			err = repo.Update(ctx, tc.input)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				err = tx.Rollback()
				require.NoError(t, err)

			} else {
				require.NoError(t, err)
				err = tx.Commit()
				require.NoError(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestSQLXUserRepositoryDelete(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := users.NewSQLXUserRepository()

	testCases := []struct {
		name          string
		id            string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful delete of a user",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectExec(`^DELETE FROM user`).
					WithArgs("1").
					WillReturnResult(sqlxmock.NewResult(1, 1))
				mockDB.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Failed delete user due to database error",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectExec(`^DELETE FROM user`).WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			expectedError: fmt.Errorf("failed to delete user: %w", errors.New("db error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			ctx := context.Background()
			tx, err := database.BeginTxx(ctx, nil)
			require.NoError(t, err)

			ctx = db.SaveToContext(ctx, tx)

			err = repo.Delete(ctx, tc.id)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				err = tx.Rollback()
				require.NoError(t, err)

			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedError, err)
				err = tx.Commit()
				require.NoError(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}
