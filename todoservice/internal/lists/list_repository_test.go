package lists_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/db"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/lists"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"time"
)

func TestSQLXListRepositoryCreate(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := lists.NewSQLXListRepository()

	testCases := []struct {
		name          string
		input         models.List
		setupMocks    func()
		expectedID    string
		expectedError error
	}{
		{
			name: "Successful creation",
			input: models.List{
				ID:          "1",
				Name:        "Test List",
				Description: "Test Description",
				OwnerID:     "owner-id",
				Visibility:  constants.VisibilityShared,
				Tags:        pkg.JSONRawMessageFromNullableString(pkg.NewValidNullableString("tag1, tag2")),
				SharedWith:  []string{"user1", "user2"},
			},
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(`^INSERT INTO lists`).WithArgs(
					"1", "Test List", "Test Description", "owner-id", constants.VisibilityShared, "tag1, tag2", sqlxmock.AnyArg(), sqlxmock.AnyArg(),
				).WillReturnRows(sqlxmock.NewRows([]string{"id"}).AddRow("1"))

				mockDB.ExpectExec(`^INSERT INTO list_access`).WithArgs("owner-id", "1", "admin", "owner").WillReturnResult(sqlxmock.NewResult(1, 1))

				for _, userID := range []string{"user1", "user2"} {
					mockDB.ExpectQuery(`^SELECT role FROM users`).WithArgs(userID).WillReturnRows(sqlxmock.NewRows([]string{"role"}).AddRow("role"))
					mockDB.ExpectExec(`^INSERT INTO list_access`).WithArgs(userID, "1", "role", "pending").WillReturnResult(sqlxmock.NewResult(1, 1))
				}
				mockDB.ExpectCommit()
			},
			expectedID:    "1",
			expectedError: nil,
		},
		{
			name: "Failed creation due to database error",
			input: models.List{
				ID:          "1",
				Name:        "Test List",
				Description: "Test Description",
				OwnerID:     "owner-id",
				Visibility:  constants.VisibilityShared,
				Tags:        pkg.JSONRawMessageFromNullableString(pkg.NewValidNullableString("tag1, tag2")),
				SharedWith:  []string{"user1", "user2"},
			},
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(`^INSERT INTO lists`).WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			expectedID:    "",
			expectedError: fmt.Errorf("failed to create list: %w", errors.New("db error")),
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

func TestSQLXListRepositoryGet(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := lists.NewSQLXListRepository()

	testCases := []struct {
		name          string
		id            string
		setupMocks    func()
		expectedList  models.List
		expectedError error
	}{
		{
			name: "Successful get of a list",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT id, name, description, owner_id, visibility, tags, created_at, updated_at FROM lists").WithArgs(
					"1").WillReturnRows(sqlxmock.NewRows([]string{"id", "name", "description", "owner_id", "visibility", "tags", "created_at", "updated_at"}).
					AddRow("1", "Test List", "Test Description", "owner-id", constants.VisibilityShared, "tag1, tag2", time.Time{}, time.Time{}))

				mockDB.ExpectQuery(`^SELECT user_id FROM list_access`).
					WithArgs("1").
					WillReturnRows(sqlxmock.NewRows([]string{"user_id"}).AddRow("user1").AddRow("user2"))
				mockDB.ExpectCommit()
			},
			expectedList: models.List{
				ID:          "1",
				Name:        "Test List",
				Description: "Test Description",
				OwnerID:     "owner-id",
				Visibility:  constants.VisibilityShared,
				Tags:        pkg.JSONRawMessageFromNullableString(pkg.NewValidNullableString("tag1, tag2")),
				SharedWith:  []string{"user1", "user2"},
			},
			expectedError: nil,
		},
		{
			name: "Failed get list due to database error",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT id, name, description, owner_id, visibility, tags, created_at, updated_at FROM lists").WithArgs("1").WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			expectedList:  models.List{},
			expectedError: fmt.Errorf("failed to get list: %w", errors.New("db error")),
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

func TestSQLXListRepositoryUpdate(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := lists.NewSQLXListRepository()

	testCases := []struct {
		name          string
		id            string
		setupMocks    func()
		expectedError error
		input         models.List
	}{
		{
			name: "Successful get of a list",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectExec("UPDATE lists").
					WithArgs("Test List", "Test Description", constants.VisibilityShared, "tag1, tag2", "1").
					WillReturnResult(sqlxmock.NewResult(1, 1))
				mockDB.ExpectCommit()
			},
			input: models.List{
				ID:          "1",
				Name:        "Test List",
				Description: "Test Description",
				OwnerID:     "owner-id",
				Visibility:  constants.VisibilityShared,
				Tags:        pkg.JSONRawMessageFromNullableString(pkg.NewValidNullableString("tag1, tag2")),
				SharedWith:  []string{"user1", "user2"},
			},
			expectedError: nil,
		},
		{
			name: "Failed to update list due to database error",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectExec("UPDATE lists").WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			input: models.List{
				ID:          "1",
				Name:        "Test List",
				Description: "Test Description",
				OwnerID:     "owner-id",
				Visibility:  constants.VisibilityShared,
				Tags:        pkg.JSONRawMessageFromNullableString(pkg.NewValidNullableString("tag1, tag2")),
				SharedWith:  []string{"user1", "user2"},
			},
			expectedError: fmt.Errorf("failed to update list: %w", errors.New("db error")),
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

func TestSQLXListRepositoryDelete(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := lists.NewSQLXListRepository()

	testCases := []struct {
		name          string
		id            string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful delete of a list",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectExec(`^DELETE FROM lists`).
					WithArgs("1").
					WillReturnResult(sqlxmock.NewResult(1, 1))
				mockDB.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Failed delete list due to database error",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectExec(`^DELETE FROM lists`).WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			expectedError: fmt.Errorf("failed to delete list: %w", errors.New("db error")),
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

func TestSQLXListRepositoryListAllByUserID(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := lists.NewSQLXListRepository()

	testCases := []struct {
		name          string
		userID        string
		setupMocks    func()
		expectedLists []models.Access
		expectedError error
	}{
		{
			name:   "Successful fetch of lists",
			userID: "user_id",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT list_id, user_id, access_level, status FROM list_access").
					WithArgs("user_id").
					WillReturnRows(sqlxmock.NewRows([]string{"list_id", "user_id", "access_level", "status"}).
						AddRow("1", "user_id", constants.Reader, "owner").
						AddRow("2", "user_id", constants.Reader, "owner"))
				mockDB.ExpectCommit()
			},
			expectedLists: []models.Access{
				{
					ListID: "1",
					UserID: "user_id",
					Role:   constants.Reader,
					Status: "owner",
				},
				{
					ListID: "2",
					UserID: "user_id",
					Role:   constants.Reader,
					Status: "owner",
				},
			},
			expectedError: nil,
		},
		{
			name:   "No lists found for user",
			userID: "owner-id",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT list_id, user_id, access_level, status FROM list_access").
					WithArgs("owner-id").
					WillReturnError(sql.ErrNoRows)
				mockDB.ExpectRollback()
			},
			expectedLists: nil,
			expectedError: fmt.Errorf("no lists found for user with id %s", "owner-id"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			ctx := context.Background()
			tx, err := database.BeginTxx(ctx, nil)
			require.NoError(t, err)

			ctx = db.SaveToContext(ctx, tx)

			listsResult, err := repo.ListAllByUserID(ctx, tc.userID)

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

func TestSQLXListRepositoryListAll(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := lists.NewSQLXListRepository()

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedList  []models.List
		expectedError error
	}{
		{
			name: "Successful get of all lists",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT id, name, description, owner_id, visibility, tags, created_at, updated_at FROM lists").
					WillReturnRows(sqlxmock.NewRows([]string{"id", "name", "description", "owner_id", "visibility", "tags", "created_at", "updated_at"}).
						AddRow("1", "Test List", "Test Description", "owner-id", constants.VisibilityShared, "tag1, tag2", time.Time{}, time.Time{}).
						AddRow("1", "Test List", "Test Description", "owner-id", constants.VisibilityShared, "tag1, tag2", time.Time{}, time.Time{}))

				mockDB.ExpectQuery(`^SELECT user_id FROM list_access`).
					WithArgs("1").
					WillReturnRows(sqlxmock.NewRows([]string{"user_id"}).AddRow("user1").AddRow("user2"))
				mockDB.ExpectQuery(`^SELECT user_id FROM list_access`).
					WithArgs("1").
					WillReturnRows(sqlxmock.NewRows([]string{"user_id"}).AddRow("user1").AddRow("user2"))
				mockDB.ExpectCommit()
			},
			expectedList: []models.List{
				{
					ID:          "1",
					Name:        "Test List",
					Description: "Test Description",
					OwnerID:     "owner-id",
					Visibility:  constants.VisibilityShared,
					Tags:        pkg.JSONRawMessageFromNullableString(pkg.NewValidNullableString("tag1, tag2")),
					SharedWith:  []string{"user1", "user2"},
				},
				{
					ID:          "1",
					Name:        "Test List",
					Description: "Test Description",
					OwnerID:     "owner-id",
					Visibility:  constants.VisibilityShared,
					Tags:        pkg.JSONRawMessageFromNullableString(pkg.NewValidNullableString("tag1, tag2")),
					SharedWith:  []string{"user1", "user2"},
				},
			},
			expectedError: nil,
		},
		{
			name: "Failed get all lists due to database error",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT id, name, description, owner_id, visibility, tags, created_at, updated_at FROM lists").WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			expectedList:  []models.List{},
			expectedError: fmt.Errorf("failed to get lists: %w", errors.New("db error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			ctx := context.Background()
			tx, err := database.BeginTxx(ctx, nil)
			require.NoError(t, err)

			ctx = db.SaveToContext(ctx, tx)

			allLists, err := repo.GetAll(ctx)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				err = tx.Rollback()
				require.NoError(t, err)

			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedList, allLists)
				err = tx.Commit()
				require.NoError(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestSQLXListRepositoryGetUsersByListID(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := lists.NewSQLXListRepository()

	testCases := []struct {
		name           string
		id             string
		setupMocks     func()
		expectedAccess []models.Access
		expectedError  error
	}{
		{
			name: "Successful get of users for a list",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT list_id, user_id, access_level, status FROM list_access").WithArgs(
					"1").WillReturnRows(sqlxmock.NewRows([]string{"list_id", "user_id", "access_level"}).
					AddRow("1", "user_id", constants.Reader))

				mockDB.ExpectCommit()
			},
			expectedAccess: []models.Access{
				{
					ListID: "1",
					UserID: "user_id",
					Role:   constants.Reader,
				},
			},
			expectedError: nil,
		},
		{
			name: "Failed get users for a list due to database error",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT list_id, user_id, access_level, status FROM list_access").WithArgs("1").WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			expectedAccess: []models.Access{},
			expectedError:  fmt.Errorf("failed to get access: %w", errors.New("db error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			ctx := context.Background()
			tx, err := database.BeginTxx(ctx, nil)
			require.NoError(t, err)

			ctx = db.SaveToContext(ctx, tx)

			users, err := repo.GetUsersByListID(ctx, tc.id)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				err = tx.Rollback()
				require.NoError(t, err)

			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedAccess, users)
				err = tx.Commit()
				require.NoError(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestSQLXListRepositoryGetListOwnerID(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := lists.NewSQLXListRepository()

	testCases := []struct {
		name          string
		id            string
		setupMocks    func()
		expectedID    string
		expectedError error
	}{
		{
			name: "Successful get of the ownerID for a list",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT owner_id FROM lists").WithArgs(
					"1").WillReturnRows(sqlxmock.NewRows([]string{"owner_id"}).
					AddRow("user1"))

				mockDB.ExpectCommit()
			},
			expectedID:    "user1",
			expectedError: nil,
		},
		{
			name: "Failed get ownerID for a list due to database error",
			id:   "1",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT owner_id FROM lists").WithArgs("1").WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			expectedID:    "",
			expectedError: fmt.Errorf("failed to get access: %w", errors.New("db error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			ctx := context.Background()
			tx, err := database.BeginTxx(ctx, nil)
			require.NoError(t, err)

			ctx = db.SaveToContext(ctx, tx)

			ownerID, err := repo.GetListOwnerID(ctx, tc.id)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				err = tx.Rollback()
				require.NoError(t, err)

			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedID, ownerID)
				err = tx.Commit()
				require.NoError(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestSQLXListRepositoryCreateAccess(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := lists.NewSQLXListRepository()

	testCases := []struct {
		name           string
		input          models.Access
		setupMocks     func()
		expectedAccess models.Access
		expectedError  error
	}{
		{
			name: "Successful creation of a list access",
			input: models.Access{
				ListID: "list_id",
				UserID: "user_id",
				Role:   constants.Reader,
			},
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectExec(`^INSERT INTO list_access`).WithArgs(
					"list_id", "user_id", constants.Reader, "pending",
				).WillReturnResult(sqlxmock.NewResult(1, 1))

				mockDB.ExpectCommit()
			},
			expectedAccess: models.Access{
				ListID: "list_id",
				UserID: "user_id",
				Role:   constants.Reader,
			},
			expectedError: nil,
		},
		{
			name: "Failed creation of a list access due to database error",
			input: models.Access{
				ListID: "list_id",
				UserID: "user_id",
				Role:   constants.Reader,
			},
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectExec(`^INSERT INTO list_access`).WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			expectedAccess: models.Access{},
			expectedError:  fmt.Errorf("failed to create list access: %w", errors.New("db error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			ctx := context.Background()
			tx, err := database.BeginTxx(ctx, nil)
			require.NoError(t, err)

			ctx = db.SaveToContext(ctx, tx)

			access, err := repo.CreateAccess(ctx, tc.input)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				err = tx.Rollback()
				require.NoError(t, err)

			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedAccess, access)
				err = tx.Commit()
				require.NoError(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestSQLXListRepositoryGetAccess(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := lists.NewSQLXListRepository()
	listID := "list1"
	userID := "user1"

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedList  models.Access
		expectedError error
	}{
		{
			name: "Successful get of a list access",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT list_id, user_id, access_level, status FROM list_access").WithArgs(
					"list1", "user1").WillReturnRows(sqlxmock.NewRows([]string{"list_id", "user_id", "access_level"}).
					AddRow("list1", "user1", constants.Reader))

				mockDB.ExpectCommit()
			},
			expectedList: models.Access{
				ListID: "list1",
				UserID: "user1",
				Role:   constants.Reader,
			},
			expectedError: nil,
		},
		{
			name: "Failed get list access due to database error",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("SELECT list_id, user_id, access_level, status FROM list_access").
					WithArgs("list1", "user1").WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			expectedList:  models.Access{},
			expectedError: fmt.Errorf("failed to get list: %w", errors.New("db error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			ctx := context.Background()
			tx, err := database.BeginTxx(ctx, nil)
			require.NoError(t, err)

			ctx = db.SaveToContext(ctx, tx)

			access, err := repo.GetAccess(ctx, listID, userID)

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				err = tx.Rollback()
				require.NoError(t, err)

			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedList, access)
				err = tx.Commit()
				require.NoError(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestSQLXListRepositoryDeleteAccess(t *testing.T) {
	database, mockDB, err := sqlxmock.Newx()
	require.NoError(t, err)
	repo := lists.NewSQLXListRepository()
	listID := "list1"
	userID := "user1"

	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful delete of a list access",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectExec(`^DELETE FROM list_access`).
					WithArgs("list1", "user1").
					WillReturnResult(sqlxmock.NewResult(1, 1))
				mockDB.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Failed delete list access due to database error",
			setupMocks: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectExec(`^DELETE FROM list_access`).WillReturnError(errors.New("db error"))
				mockDB.ExpectRollback()
			},
			expectedError: fmt.Errorf("failed to delete list: %w", errors.New("db error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			ctx := context.Background()
			tx, err := database.BeginTxx(ctx, nil)
			require.NoError(t, err)

			ctx = db.SaveToContext(ctx, tx)

			err = repo.DeleteAccess(ctx, listID, userID)

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
