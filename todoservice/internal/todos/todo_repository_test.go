package todos_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/db"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/todos"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"time"
)

func TestSQLXTodoRepositoryCreate(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := todos.NewSQLXTodoRepository()
	date := time.Time{}
	assignedTo := "someone"

	testCases := []struct {
		name          string
		input         models.Todo
		setupMocks    func()
		expectedID    string
		expectedError error
	}{
		{
			name: "Successful creation",
			input: models.Todo{
				ID:          "4cfd7e64-7431-4690-a2a0-1268917cedf3",
				Title:       "Test Todo",
				Description: "Test Description",
				ListID:      "4cfd7e64-7431-4690-a2a0-1268917cedf3",
				Priority:    constants.PriorityLow,
				Tags:        pkg.JSONRawMessageFromNullableString(pkg.NewValidNullableString("tag1, tag2")),
				DueDate:     &date,
				StartDate:   &date,
				AssignedTo:  &assignedTo,
			},
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(`^INSERT INTO todos`).WithArgs(
					"4cfd7e64-7431-4690-a2a0-1268917cedf3", "Test Todo", "Test Description", "4cfd7e64-7431-4690-a2a0-1268917cedf3", false, "tag1, tag2", constants.PriorityLow, nil, nil, "someone", sqlxmock.AnyArg(), sqlxmock.AnyArg(),
				).WillReturnRows(sqlxmock.NewRows([]string{"id"}).AddRow("1"))

				mockDB.ExpectCommit()
			},
			expectedID:    "1",
			expectedError: nil,
		},
		{
			name: "Failed creation due to database error",
			input: models.Todo{
				ID:          "4cfd7e64-7431-4690-a2a0-1268917cedf3",
				Title:       "Test Todo",
				Description: "Test Description",
				ListID:      "4cfd7e64-7431-4690-a2a0-1268917cedf3",
				Priority:    constants.PriorityLow,
				Tags:        pkg.JSONRawMessageFromNullableString(pkg.NewValidNullableString("tag1, tag2")),
				DueDate:     &date,
				StartDate:   &date,
				AssignedTo:  &assignedTo,
			},
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(`^INSERT INTO todos`).WithArgs(
					"4cfd7e64-7431-4690-a2a0-1268917cedf3", "Test Todo", "Test Description", "4cfd7e64-7431-4690-a2a0-1268917cedf3", false, "tag1, tag2", constants.PriorityLow, nil, nil, "someone", sqlxmock.AnyArg(), sqlxmock.AnyArg(),
				).WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			expectedID:    "",
			expectedError: fmt.Errorf("failed to create todo: %w", errors.New("db error")),
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

func TestSQLXTodoRepositoryGet(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := todos.NewSQLXTodoRepository()
	assignedTo := ""

	testCases := []struct {
		name          string
		id            string
		setupMocks    func()
		expectedList  models.Todo
		expectedError error
	}{
		{
			name: "Successful get of a todo",
			id:   "4cfd7e64-7431-4690-a2a0-1268917cedf3",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT id, title, description, list_id, priority, due_date, start_date, completed, tags, created_at, updated_at, assigned_to FROM todos").WithArgs(
					"4cfd7e64-7431-4690-a2a0-1268917cedf3").WillReturnRows(sqlxmock.NewRows([]string{"id", "title", "description", "list_id", "priority", "due_date", "start_date", "completed", "tags", "assigned_to", "created_at", "updated_at"}).
					AddRow("4cfd7e64-7431-4690-a2a0-1268917cedf3", "Test Todo", "Test Description", "4cfd7e64-7431-4690-a2a0-1268917cedf3", constants.PriorityLow, time.Time{}, time.Time{}, false, "tag1, tag2", "", time.Time{}, time.Time{}))

				mockDB.ExpectCommit()
			},
			expectedList: models.Todo{
				ID:          "4cfd7e64-7431-4690-a2a0-1268917cedf3",
				Title:       "Test Todo",
				Description: "Test Description",
				ListID:      "4cfd7e64-7431-4690-a2a0-1268917cedf3",
				Priority:    constants.PriorityLow,
				Tags:        pkg.JSONRawMessageFromNullableString(pkg.NewValidNullableString("tag1, tag2")),
				DueDate:     &time.Time{},
				StartDate:   &time.Time{},
				AssignedTo:  &assignedTo,
			},
			expectedError: nil,
		},
		{
			name: "Failed get todo due to database error",
			id:   "4cfd7e64-7431-4690-a2a0-1268917cedf3",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT id, title, description, list_id, priority, due_date, start_date, completed, tags, created_at, updated_at, assigned_to FROM todos").WithArgs(
					"4cfd7e64-7431-4690-a2a0-1268917cedf3").WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			expectedList:  models.Todo{},
			expectedError: fmt.Errorf("failed to get todo: %w", errors.New("db error")),
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

func TestSQLXTodoRepositoryUpdate(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := todos.NewSQLXTodoRepository()
	date := time.Time{}
	assignedTo := "someone"

	testCases := []struct {
		name          string
		id            string
		setupMocks    func()
		expectedError error
		input         models.Todo
	}{
		{
			name: "Successful get of a todo",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectExec("UPDATE todos").
					WithArgs("Test Todo", "Test Description", constants.PriorityLow, nil, nil, false, "tag1, tag2", "someone", "4cfd7e64-7431-4690-a2a0-1268917cedf3").
					WillReturnResult(sqlxmock.NewResult(1, 1))

				mockDB.ExpectCommit()
			},
			input: models.Todo{
				ID:          "4cfd7e64-7431-4690-a2a0-1268917cedf3",
				Title:       "Test Todo",
				Description: "Test Description",
				ListID:      "4cfd7e64-7431-4690-a2a0-1268917cedf3",
				Priority:    constants.PriorityLow,
				Tags:        pkg.JSONRawMessageFromNullableString(pkg.NewValidNullableString("tag1, tag2")),
				DueDate:     &date,
				StartDate:   &date,
				AssignedTo:  &assignedTo,
			},
			expectedError: nil,
		},
		{
			name: "Failed to update list due to database error",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectExec("UPDATE todos").
					WithArgs("Test Todo", "Test Description", constants.PriorityLow, nil, nil, false, "tag1, tag2", "someone", "4cfd7e64-7431-4690-a2a0-1268917cedf3").
					WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			input: models.Todo{
				ID:          "4cfd7e64-7431-4690-a2a0-1268917cedf3",
				Title:       "Test Todo",
				Description: "Test Description",
				ListID:      "4cfd7e64-7431-4690-a2a0-1268917cedf3",
				Priority:    constants.PriorityLow,
				Tags:        pkg.JSONRawMessageFromNullableString(pkg.NewValidNullableString("tag1, tag2")),
				DueDate:     &date,
				StartDate:   &date,
				AssignedTo:  &assignedTo,
			},
			expectedError: fmt.Errorf("failed to update todo: %w", errors.New("db error")),
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

func TestSQLXTodoRepositoryDelete(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := todos.NewSQLXTodoRepository()

	testCases := []struct {
		name          string
		id            string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful delete of a todo",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectExec(`^DELETE FROM todos`).
					WithArgs("1").
					WillReturnResult(sqlxmock.NewResult(1, 1))
				mockDB.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Failed delete todo due to database error",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectExec(`^DELETE FROM todos`).WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			expectedError: fmt.Errorf("failed to delete todo: %w", errors.New("db error")),
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

func TestSQLXTodoRepositoryGetAllByListID(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := todos.NewSQLXTodoRepository()

	testCases := []struct {
		name          string
		userID        string
		setupMocks    func()
		expectedLists []models.Todo
		expectedError error
	}{
		{
			name:   "Successful fetch of todos",
			userID: "owner_id",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT id, title, description, list_id, priority, due_date, start_date, completed, tags, created_at, updated_at FROM todos").
					WithArgs("owner_id").
					WillReturnRows(sqlxmock.NewRows([]string{"id", "title", "description", "list_id", "priority", "due_date", "start_date", "completed", "tags", "created_at", "updated_at"}).
						AddRow("1", "Todo 1", "Desc 1", "owner_id", constants.PriorityLow, nil, nil, false, "tag1, tag2", time.Time{}, time.Time{}).
						AddRow("2", "Todo 2", "Desc 2", "owner_id", constants.PriorityLow, nil, nil, false, "tag1, tag2", time.Time{}, time.Time{}))
				mockDB.ExpectCommit()
			},
			expectedLists: []models.Todo{
				{
					ID:          "1",
					Title:       "Todo 1",
					Description: "Desc 1",
					ListID:      "owner_id",
					Priority:    constants.PriorityLow,
					Tags:        pkg.JSONRawMessageFromNullableString(pkg.NewValidNullableString("tag1, tag2")),
				},
				{
					ID:          "2",
					Title:       "Todo 2",
					Description: "Desc 2",
					ListID:      "owner_id",
					Priority:    constants.PriorityLow,
					Tags:        pkg.JSONRawMessageFromNullableString(pkg.NewValidNullableString("tag1, tag2")),
				},
			},
			expectedError: nil,
		},
		{
			name:   "No todos found for user",
			userID: "owner-id",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT id, title, description, list_id, priority, due_date, start_date, completed, tags, created_at, updated_at FROM todos").
					WithArgs("owner-id").
					WillReturnError(sql.ErrNoRows)
				mockDB.ExpectRollback()
			},
			expectedLists: nil,
			expectedError: fmt.Errorf("no todos found for list with id %s", "owner-id"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			ctx := context.Background()
			tx, err := database.BeginTxx(ctx, nil)
			require.NoError(t, err)

			ctx = db.SaveToContext(ctx, tx)

			listsResult, err := repo.GetAllByListID(ctx, tc.userID)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				err = tx.Rollback()
				require.NoError(t, err)

			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedLists, listsResult)
				err = tx.Commit()
				require.NoError(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}
