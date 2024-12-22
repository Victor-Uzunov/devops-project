package lists_test

import (
	"context"
	"errors"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/lists"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/lists/automock"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestServiceCreateList(t *testing.T) {
	id := "1"
	mockTime := time.Time{}
	ctx := context.Background()
	err := errors.New("error")

	modelInput := models.List{
		Name:        "Test List",
		Description: "Test description",
		OwnerID:     "1",
		SharedWith:  nil,
	}

	model := models.List{
		ID:          id,
		Name:        "Test List",
		Description: "Test description",
		OwnerID:     "1",
		SharedWith:  nil,
		CreatedAt:   mockTime,
		UpdatedAt:   mockTime,
	}

	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.ListRepository
		timeService   func() *automock.TimeService
		input         models.List
		expectedError error
	}{
		{
			name: "Create new list",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				uuidService.EXPECT().Generate().Return(id).Once()
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}

				repo.EXPECT().Create(ctx, model).Return(id, nil).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				timeService.EXPECT().Now().Return(mockTime).Twice()
				return timeService
			},
			input:         modelInput,
			expectedError: nil,
		},
		{
			name: "Error when repo create fails",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				uuidService.EXPECT().Generate().Return(id).Once()
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().Create(ctx, model).Return("", err).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				timeService.EXPECT().Now().Return(mockTime).Twice()
				return timeService
			},
			input:         modelInput,
			expectedError: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeService := tt.timeService()
			uuidService := tt.uuidService()
			repo := tt.repo()
			defer mock.AssertExpectationsForObjects(t, timeService, repo, uuidService)

			svc := lists.NewService(repo, uuidService, timeService)
			_, err := svc.CreateList(ctx, tt.input)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceGetList(t *testing.T) {
	id := "1"
	mockTime := time.Time{}
	ctx := context.Background()
	err := errors.New("error")

	model := models.List{
		ID:          id,
		Name:        "Test List",
		Description: "Test description",
		OwnerID:     "1",
		SharedWith:  nil,
		CreatedAt:   mockTime,
		UpdatedAt:   mockTime,
	}
	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.ListRepository
		timeService   func() *automock.TimeService
		expectedError error
	}{
		{
			name: "Get list",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().Get(ctx, id).Return(model, nil).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: nil,
		},
		{
			name: "Error when repo get fails",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().Get(ctx, id).Return(models.List{}, err).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeService := tt.timeService()
			uuidService := tt.uuidService()
			repo := tt.repo()
			defer mock.AssertExpectationsForObjects(t, repo)

			svc := lists.NewService(repo, uuidService, timeService)
			_, err := svc.GetList(ctx, id)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceUpdateList(t *testing.T) {
	id := "1"
	mockTime := time.Time{}
	ctx := context.Background()
	err := errors.New("error")

	modelInput := models.List{
		ID:          id,
		Name:        "Test List",
		Description: "Test description",
		OwnerID:     "1",
		SharedWith:  nil,
	}

	model := models.List{
		ID:          id,
		Name:        "Test List",
		Description: "Test description",
		OwnerID:     "1",
		SharedWith:  nil,
		CreatedAt:   mockTime,
		UpdatedAt:   mockTime,
	}

	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.ListRepository
		timeService   func() *automock.TimeService
		input         models.List
		expectedError error
	}{
		{
			name:  "Update existing list",
			input: modelInput,
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().Get(ctx, id).Return(model, nil).Once()
				repo.EXPECT().Update(ctx, model).Return(nil).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				timeService.EXPECT().Now().Return(mockTime).Once()
				return timeService
			},
			expectedError: nil,
		},
		{
			name: "Error when repo update fails",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().Get(ctx, id).Return(models.List{}, err).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				timeService.EXPECT().Now().Return(mockTime).Once()
				return timeService
			},
			input:         modelInput,
			expectedError: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeService := tt.timeService()
			uuidService := tt.uuidService()
			repo := tt.repo()
			defer mock.AssertExpectationsForObjects(t, repo, uuidService)

			svc := lists.NewService(repo, uuidService, timeService)
			err := svc.UpdateList(ctx, tt.input)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceDeleteList(t *testing.T) {
	id := "1"
	err := errors.New("error")
	ctx := context.Background()

	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.ListRepository
		timeService   func() *automock.TimeService
		expectedError error
	}{
		{
			name: "Delete existing list",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().Delete(ctx, id).Return(nil).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: nil,
		},
		{
			name: "Error when repo delete fails",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().Delete(ctx, id).Return(err).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuidService := tt.uuidService()
			repo := tt.repo()
			timeService := tt.timeService()
			defer mock.AssertExpectationsForObjects(t, repo)

			svc := lists.NewService(repo, uuidService, timeService)
			err := svc.DeleteList(ctx, id)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceListAllByUserID(t *testing.T) {
	id := "1"
	err := errors.New("error")
	ctx := context.Background()

	accessModel := models.Access{
		ListID: "Test",
		UserID: "1",
	}

	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.ListRepository
		timeService   func() *automock.TimeService
		expectedError error
	}{
		{
			name: "Get all lists by user id",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().ListAllByUserID(ctx, id).Return([]models.Access{accessModel}, nil).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: nil,
		},
		{
			name: "Error when repo list all by userID fails",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().ListAllByUserID(ctx, id).Return([]models.Access{}, err).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuidService := tt.uuidService()
			repo := tt.repo()
			timeService := tt.timeService()
			defer mock.AssertExpectationsForObjects(t, repo, uuidService, timeService)

			svc := lists.NewService(repo, uuidService, timeService)
			_, err := svc.ListAllByUserID(ctx, id)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceGetAllLists(t *testing.T) {
	id := "1"
	mockTime := time.Time{}
	ctx := context.Background()
	err := errors.New("error")

	model := []models.List{
		{
			ID:          id,
			Name:        "Test List",
			Description: "Test description",
			OwnerID:     "1",
			SharedWith:  nil,
			CreatedAt:   mockTime,
			UpdatedAt:   mockTime,
		},
	}
	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.ListRepository
		timeService   func() *automock.TimeService
		expectedError error
	}{
		{
			name: "Get all lists",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().GetAll(ctx).Return(model, nil).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: nil,
		},
		{
			name: "Error when service get all lists fails",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().GetAll(ctx).Return([]models.List{}, err).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeService := tt.timeService()
			uuidService := tt.uuidService()
			repo := tt.repo()
			defer mock.AssertExpectationsForObjects(t, repo)

			svc := lists.NewService(repo, uuidService, timeService)
			_, err := svc.GetAllLists(ctx)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceGetUsersByListID(t *testing.T) {
	id := "1"
	ctx := context.Background()
	err := errors.New("error")
	modelAccess := []models.Access{
		{
			ListID: "1",
			UserID: "user1",
			Role:   "admin",
		},
		{
			ListID: "1",
			UserID: "user2",
			Role:   "writer",
		},
		{
			ListID: "1",
			UserID: "user3",
			Role:   "reader",
		},
	}
	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.ListRepository
		timeService   func() *automock.TimeService
		expectedError error
	}{
		{
			name: "Get users by listID",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().GetUsersByListID(ctx, id).Return(modelAccess, nil).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: nil,
		},
		{
			name: "Error when service get users by listID fails",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().GetUsersByListID(ctx, id).Return([]models.Access{}, err).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeService := tt.timeService()
			uuidService := tt.uuidService()
			repo := tt.repo()
			defer mock.AssertExpectationsForObjects(t, repo)

			svc := lists.NewService(repo, uuidService, timeService)
			_, err := svc.GetUsersByListID(ctx, id)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceListOwnerIDs(t *testing.T) {
	id := "1"
	ctx := context.Background()
	err := errors.New("error")
	ownerID := "1"
	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.ListRepository
		timeService   func() *automock.TimeService
		expectedError error
	}{
		{
			name: "Get list ownerID",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().GetListOwnerID(ctx, id).Return(ownerID, nil).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: nil,
		},
		{
			name: "Error when service get list ownerID fails",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().GetListOwnerID(ctx, id).Return("", err).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeService := tt.timeService()
			uuidService := tt.uuidService()
			repo := tt.repo()
			defer mock.AssertExpectationsForObjects(t, repo)

			svc := lists.NewService(repo, uuidService, timeService)
			_, err := svc.GetListOwnerID(ctx, id)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceCreateAccess(t *testing.T) {
	ctx := context.Background()
	err := errors.New("error")

	modelInput := models.Access{
		ListID: "1",
		UserID: "user1",
		Role:   "reader",
	}

	model := models.Access{
		ListID: "1",
		UserID: "user1",
		Role:   "reader",
	}

	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.ListRepository
		timeService   func() *automock.TimeService
		input         models.Access
		expectedError error
	}{
		{
			name: "Create new list_access",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}

				repo.EXPECT().CreateAccess(ctx, model).Return(model, nil).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			input:         modelInput,
			expectedError: nil,
		},
		{
			name: "Error when service create new access fails",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().CreateAccess(ctx, model).Return(models.Access{}, err).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			input:         modelInput,
			expectedError: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeService := tt.timeService()
			uuidService := tt.uuidService()
			repo := tt.repo()
			defer mock.AssertExpectationsForObjects(t, timeService, repo, uuidService)

			svc := lists.NewService(repo, uuidService, timeService)
			_, err := svc.CreateAccess(ctx, tt.input)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceGetAccess(t *testing.T) {
	ctx := context.Background()
	err := errors.New("error")
	listID := "listID"
	userID := "userID"

	model := models.Access{
		ListID: "1",
		UserID: "user1",
		Role:   "reader",
	}
	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.ListRepository
		timeService   func() *automock.TimeService
		expectedError error
	}{
		{
			name: "Get list access",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().GetAccess(ctx, listID, userID).Return(model, nil).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: nil,
		},
		{
			name: "Error when service get list access fails",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().GetAccess(ctx, listID, userID).Return(models.Access{}, err).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeService := tt.timeService()
			uuidService := tt.uuidService()
			repo := tt.repo()
			defer mock.AssertExpectationsForObjects(t, repo)

			svc := lists.NewService(repo, uuidService, timeService)
			_, err := svc.GetAccess(ctx, listID, userID)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceDeleteAccess(t *testing.T) {
	err := errors.New("error")
	ctx := context.Background()
	listID := "listID"
	userID := "userID"

	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.ListRepository
		timeService   func() *automock.TimeService
		expectedError error
	}{
		{
			name: "Delete existing list access",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().DeleteAccess(ctx, listID, userID).Return(nil).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: nil,
		},
		{
			name: "Error when service delete access fails",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.ListRepository {
				repo := &automock.ListRepository{}
				repo.EXPECT().DeleteAccess(ctx, listID, userID).Return(err).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuidService := tt.uuidService()
			repo := tt.repo()
			timeService := tt.timeService()
			defer mock.AssertExpectationsForObjects(t, repo)

			svc := lists.NewService(repo, uuidService, timeService)
			err := svc.DeleteAccess(ctx, listID, userID)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
