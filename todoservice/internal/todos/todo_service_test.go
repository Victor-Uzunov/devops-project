package todos_test

import (
	"context"
	"errors"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/todos"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/todos/automock"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestServiceCreateTodo(t *testing.T) {
	id := "1"
	mockTime := time.Time{}
	err := errors.New("error")
	ctx := context.Background()

	modelInput := models.Todo{
		Title:       "Test Todo",
		Description: "Test description",
		ListID:      "1",
		Priority:    constants.PriorityLow,
	}

	model := models.Todo{
		ID:          id,
		Title:       "Test Todo",
		Description: "Test description",
		ListID:      "1",
		Priority:    constants.PriorityLow,
		CreatedAt:   mockTime,
		UpdatedAt:   mockTime,
	}

	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.TodoRepository
		timeService   func() *automock.TimeService
		input         models.Todo
		expectedError error
	}{
		{
			name: "Create new todo",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				uuidService.EXPECT().Generate().Return(id).Once()
				return uuidService
			},
			repo: func() *automock.TodoRepository {
				repo := &automock.TodoRepository{}
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
			repo: func() *automock.TodoRepository {
				repo := &automock.TodoRepository{}
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
			uuidService := tt.uuidService()
			repo := tt.repo()
			timeService := tt.timeService()
			defer mock.AssertExpectationsForObjects(t, uuidService, repo, timeService)

			svc := todos.NewService(repo, uuidService, timeService)
			_, err := svc.CreateTodo(ctx, tt.input)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceGetTodo(t *testing.T) {
	id := "1"
	err := errors.New("error")
	ctx := context.Background()

	model := models.Todo{
		Title:       "Test Todo",
		Description: "Test description",
		ListID:      "1",
		Priority:    constants.PriorityLow,
	}

	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.TodoRepository
		timeService   func() *automock.TimeService
		expectedError error
	}{
		{
			name: "Get Todo",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.TodoRepository {
				repo := &automock.TodoRepository{}
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
			repo: func() *automock.TodoRepository {
				repo := &automock.TodoRepository{}
				repo.EXPECT().Get(ctx, id).Return(models.Todo{}, err).Once()
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
			defer mock.AssertExpectationsForObjects(t, uuidService, repo, timeService)

			svc := todos.NewService(repo, uuidService, timeService)
			_, err := svc.GetTodo(ctx, id)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceUpdateTodo(t *testing.T) {
	id := "1"
	err := errors.New("error")
	ctx := context.Background()
	mockTime := time.Time{}

	modelInput := models.Todo{
		ID:          id,
		Title:       "Test Todo",
		Description: "Test description",
		ListID:      "1",
		Priority:    constants.PriorityLow,
	}

	model := models.Todo{
		ID:          id,
		Title:       "Test Todo",
		Description: "Test description",
		ListID:      "1",
		Priority:    constants.PriorityLow,
		CreatedAt:   mockTime,
		UpdatedAt:   mockTime,
	}

	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.TodoRepository
		timeService   func() *automock.TimeService
		input         models.Todo
		expectedError error
	}{
		{
			name: "Update Todo",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.TodoRepository {
				repo := &automock.TodoRepository{}
				repo.EXPECT().Update(ctx, model).Return(nil).Once()
				repo.EXPECT().Get(ctx, id).Return(model, nil).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				timeService.EXPECT().Now().Return(mockTime).Once()
				return timeService
			},
			input:         modelInput,
			expectedError: nil,
		},
		{
			name: "Error when repo update fails",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.TodoRepository {
				repo := &automock.TodoRepository{}
				repo.EXPECT().Get(ctx, id).Return(models.Todo{}, err).Once()
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
			uuidService := tt.uuidService()
			repo := tt.repo()
			timeService := tt.timeService()
			defer mock.AssertExpectationsForObjects(t, uuidService, repo, timeService)

			svc := todos.NewService(repo, uuidService, timeService)
			err := svc.UpdateTodo(ctx, modelInput)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceDeleteTodo(t *testing.T) {
	id := "1"
	err := errors.New("error")
	ctx := context.Background()

	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.TodoRepository
		timeService   func() *automock.TimeService
		expectedError error
	}{
		{
			name: "Delete existing todo",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.TodoRepository {
				repo := &automock.TodoRepository{}
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
			repo: func() *automock.TodoRepository {
				repo := &automock.TodoRepository{}
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
			defer mock.AssertExpectationsForObjects(t, uuidService, repo, timeService)

			svc := todos.NewService(repo, uuidService, timeService)
			err := svc.DeleteTodo(ctx, id)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServiceListTodosByListID(t *testing.T) {
	id := "1"
	err := errors.New("error")
	ctx := context.Background()

	model := models.Todo{
		Title:       "Test Todo",
		Description: "Test description",
		ListID:      "1",
		Priority:    constants.PriorityLow,
	}

	tests := []struct {
		name          string
		uuidService   func() *automock.UUIDService
		repo          func() *automock.TodoRepository
		timeService   func() *automock.TimeService
		expectedError error
	}{
		{
			name: "Get todos with listID",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.TodoRepository {
				repo := &automock.TodoRepository{}
				repo.EXPECT().GetAllByListID(ctx, id).Return([]models.Todo{model}, nil).Once()
				return repo
			},
			timeService: func() *automock.TimeService {
				timeService := &automock.TimeService{}
				return timeService
			},
			expectedError: nil,
		},
		{
			name: "Error when repo list fails",
			uuidService: func() *automock.UUIDService {
				uuidService := &automock.UUIDService{}
				return uuidService
			},
			repo: func() *automock.TodoRepository {
				repo := &automock.TodoRepository{}
				repo.EXPECT().GetAllByListID(ctx, id).Return([]models.Todo{}, err).Once()
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
			defer mock.AssertExpectationsForObjects(t, uuidService, repo, timeService)

			svc := todos.NewService(repo, uuidService, timeService)
			_, err := svc.ListTodosByListID(ctx, id)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
