package users_test

import (
	"context"
	"errors"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/users"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/users/automock"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestServiceCreateUser(t *testing.T) {
	id := "1"
	mockTime := time.Time{}
	ctx := context.Background()
	err := errors.New("error")

	modelInput := models.User{
		Email:    "test",
		GithubID: "github1",
		Role:     "user",
	}

	model := models.User{
		ID:        id,
		Email:     "test",
		GithubID:  "github1",
		Role:      "user",
		CreatedAt: mockTime,
		UpdatedAt: mockTime,
	}

	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.UserRepository
		timeService   func() *automock.TimeService
		input         models.User
		expectedError error
	}{
		{
			name: "Create new user",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				uuidService.EXPECT().Generate().Return(id).Once()
				return uuidService
			},
			repo: func() *automock.UserRepository {
				repo := &automock.UserRepository{}
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
			repo: func() *automock.UserRepository {
				repo := &automock.UserRepository{}
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

			svc := users.NewService(repo, uuidService, timeService)
			_, err := svc.CreateUser(ctx, tt.input)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceGetUser(t *testing.T) {
	id := "1"
	mockTime := time.Time{}
	ctx := context.Background()
	err := errors.New("error")

	model := models.User{
		ID:        id,
		Email:     "test",
		GithubID:  "github1",
		Role:      "user",
		CreatedAt: mockTime,
		UpdatedAt: mockTime,
	}
	tests := []struct {
		name          string
		repo          func() *automock.UserRepository
		expectedError error
	}{
		{
			name: "Get user",
			repo: func() *automock.UserRepository {
				repo := &automock.UserRepository{}
				repo.EXPECT().Get(ctx, id).Return(model, nil).Once()
				return repo
			},
			expectedError: nil,
		},
		{
			name: "Error when repo get fails",
			repo: func() *automock.UserRepository {
				repo := &automock.UserRepository{}
				repo.EXPECT().Get(ctx, id).Return(models.User{}, err).Once()
				return repo
			},
			expectedError: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.repo()
			defer mock.AssertExpectationsForObjects(t, repo)

			svc := users.NewService(repo, nil, nil)
			_, err := svc.GetUser(ctx, id)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceUpdateUser(t *testing.T) {
	id := "1"
	mockTime := time.Time{}
	ctx := context.Background()
	err := errors.New("error")

	modelInput := models.User{
		ID:       id,
		Email:    "test",
		GithubID: "github1",
		Role:     "user",
	}

	model := models.User{
		ID:        id,
		Email:     "test",
		GithubID:  "github1",
		Role:      "user",
		CreatedAt: mockTime,
		UpdatedAt: mockTime,
	}

	tests := []struct {
		name          string
		repo          func() *automock.UserRepository
		input         models.User
		expectedError error
	}{
		{
			name:  "Update existing user",
			input: modelInput,
			repo: func() *automock.UserRepository {
				repo := &automock.UserRepository{}
				repo.EXPECT().Get(ctx, id).Return(model, nil).Once()
				repo.EXPECT().Update(ctx, model).Return(nil).Once()
				return repo
			},
			expectedError: nil,
		},
		{
			name: "Error when repo get fails",
			repo: func() *automock.UserRepository {
				repo := &automock.UserRepository{}
				repo.EXPECT().Get(ctx, id).Return(models.User{}, err).Once()
				return repo
			},
			input:         modelInput,
			expectedError: err,
		},
		{
			name: "Error when repo update fails",
			repo: func() *automock.UserRepository {
				repo := &automock.UserRepository{}
				repo.EXPECT().Get(ctx, id).Return(model, nil).Once()
				repo.EXPECT().Update(ctx, model).Return(err).Once()
				return repo
			},
			input:         modelInput,
			expectedError: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.repo()
			defer mock.AssertExpectationsForObjects(t, repo)

			svc := users.NewService(repo, nil, nil)
			err := svc.UpdateUser(ctx, tt.input)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceDeleteUser(t *testing.T) {
	id := "1"
	err := errors.New("error")
	ctx := context.Background()

	tests := []struct {
		name          string
		repo          func() *automock.UserRepository
		expectedError error
	}{
		{
			name: "Delete existing user",
			repo: func() *automock.UserRepository {
				repo := &automock.UserRepository{}
				repo.EXPECT().Delete(ctx, id).Return(nil).Once()
				return repo
			},
			expectedError: nil,
		},
		{
			name: "Error when repo delete fails",
			repo: func() *automock.UserRepository {
				repo := &automock.UserRepository{}
				repo.EXPECT().Delete(ctx, id).Return(err).Once()
				return repo
			},
			expectedError: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.repo()
			defer mock.AssertExpectationsForObjects(t, repo)

			svc := users.NewService(repo, nil, nil)
			err := svc.DeleteUser(ctx, id)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
