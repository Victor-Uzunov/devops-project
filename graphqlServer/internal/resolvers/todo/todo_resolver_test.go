package todo_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/generated/graphql"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/converters/automock"
	mock2 "github.com/Victor-Uzunov/devops-project/graphqlServer/internal/resolvers/mock"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/resolvers/todo"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestTodos_TodoResolver(t *testing.T) {
	expectedTodo := graphql.Todo{
		ID:    "1",
		Title: "Test Todo",
	}
	inputTodo := models.Todo{
		ID:    "1",
		Title: "Test Todo",
	}

	tests := []struct {
		name          string
		mockMethod    string
		mockURL       string
		mockResp      []byte
		mockErr       error
		expectError   bool
		expectTodos   []*graphql.Todo
		todoConverter func() *automock.TodoConverter
	}{
		{
			name:        "successful todos fetch",
			mockMethod:  "GET",
			mockURL:     "todos/all",
			mockResp:    []byte(`[{"ID": "1", "Title": "Test Todo"}]`),
			mockErr:     nil,
			expectError: false,
			expectTodos: []*graphql.Todo{
				{ID: "1", Title: "Test Todo"},
			},
			todoConverter: func() *automock.TodoConverter {
				todoConverter := &automock.TodoConverter{}
				todoConverter.EXPECT().ConvertTodoToGraphQL(inputTodo).Return(&expectedTodo, nil)
				return todoConverter
			},
		},
		{
			name:        "failed HTTP request",
			mockMethod:  "GET",
			mockURL:     "todos/all",
			mockResp:    nil,
			mockErr:     errors.New("failed to fetch todos"),
			expectError: true,
			expectTodos: nil,
			todoConverter: func() *automock.TodoConverter {
				return &automock.TodoConverter{}
			},
		},
		{
			name:        "failed to unmarshal response",
			mockMethod:  "GET",
			mockURL:     "todos/all",
			mockResp:    []byte(`invalid JSON`),
			mockErr:     nil,
			expectError: true,
			expectTodos: nil,
			todoConverter: func() *automock.TodoConverter {
				return &automock.TodoConverter{}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(mock2.ClientMock)

			mockClient.On("Do", mock.Anything, tt.mockMethod, tt.mockURL, mock.Anything).Return(tt.mockResp, tt.mockErr)

			todoConverter := tt.todoConverter()

			r := todo.NewResolver(mockClient, todoConverter, nil, nil)

			result, err := r.Todos(context.Background())

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectTodos, result)
			}

			mockClient.AssertCalled(t, "Do", mock.Anything, tt.mockMethod, tt.mockURL, mock.Anything)
		})
	}
}

func TestTodo_TodoResolver(t *testing.T) {
	expectedTodo := graphql.Todo{
		ID:    "1",
		Title: "Test Todo",
	}
	inputTodo := models.Todo{
		ID:    "1",
		Title: "Test Todo",
	}

	tests := []struct {
		name          string
		mockMethod    string
		mockURL       string
		mockResp      []byte
		mockErr       error
		expectError   bool
		expectTodo    *graphql.Todo
		todoConverter func() *automock.TodoConverter
	}{
		{
			name:        "successful todo fetch",
			mockMethod:  "GET",
			mockURL:     "todos/1",
			mockResp:    []byte(`{"ID": "1", "Title": "Test Todo"}`),
			mockErr:     nil,
			expectError: false,
			expectTodo:  &expectedTodo,
			todoConverter: func() *automock.TodoConverter {
				todoConverter := &automock.TodoConverter{}
				todoConverter.EXPECT().ConvertTodoToGraphQL(inputTodo).Return(&expectedTodo, nil)
				return todoConverter
			},
		},
		{
			name:        "failed http request",
			mockMethod:  "GET",
			mockURL:     "todos/1",
			mockResp:    nil,
			mockErr:     errors.New("failed to fetch todo"),
			expectError: true,
			expectTodo:  nil,
			todoConverter: func() *automock.TodoConverter {
				return &automock.TodoConverter{}
			},
		},
		{
			name:        "failed to unmarshal response",
			mockMethod:  "GET",
			mockURL:     "todos/1",
			mockResp:    []byte(`invalid JSON`),
			mockErr:     nil,
			expectError: true,
			expectTodo:  nil,
			todoConverter: func() *automock.TodoConverter {
				return &automock.TodoConverter{}
			},
		},
		{
			name:        "failed to convert todo",
			mockMethod:  "GET",
			mockURL:     "todos/1",
			mockResp:    []byte(`{"ID": "1", "Title": "Test Todo"}`),
			mockErr:     nil,
			expectError: true,
			expectTodo:  nil,
			todoConverter: func() *automock.TodoConverter {
				todoConverter := &automock.TodoConverter{}
				todoConverter.EXPECT().ConvertTodoToGraphQL(inputTodo).Return(nil, errors.New("conversion failed"))
				return todoConverter
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(mock2.ClientMock)
			mockClient.On("Do", mock.Anything, tt.mockMethod, tt.mockURL, mock.Anything).Return(tt.mockResp, tt.mockErr)

			todoConverter := tt.todoConverter()
			r := todo.NewResolver(mockClient, todoConverter, nil, nil)

			result, err := r.Todo(context.Background(), "1")

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectTodo, result)
			}

			mockClient.AssertCalled(t, "Do", mock.Anything, tt.mockMethod, tt.mockURL, mock.Anything)
		})
	}
}

func TestCreateTodo_TodoResolver(t *testing.T) {
	expectedTodo := graphql.Todo{
		ID:    "1",
		Title: "New Todo",
	}

	inputTodo := models.Todo{
		Title: "New Todo",
	}

	convertModel := models.Todo{
		ID:    "1",
		Title: "New Todo",
	}

	createTodoInput := graphql.CreateTodoInput{
		Title: "New Todo",
	}

	httpInput := models.Todo{
		Title: "New Todo",
	}

	data, _ := json.Marshal(inputTodo)

	tests := []struct {
		name        string
		input       graphql.CreateTodoInput
		mockConvert func() *automock.TodoConverter
		mockDo      func() *mock2.ClientMock
		mockRespID  []byte
		mockRespGet []byte
		mockErr     error
		expectError bool
		expectTodo  *graphql.Todo
	}{
		{
			name:  "successful todo creation",
			input: createTodoInput,
			mockConvert: func() *automock.TodoConverter {
				todoConverter := &automock.TodoConverter{}
				todoConverter.EXPECT().ConvertCreateTodoInput(createTodoInput).Return(httpInput, nil)
				todoConverter.EXPECT().ConvertTodoToGraphQL(convertModel).Return(&expectedTodo, nil)
				return todoConverter
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "POST", "todos/create", data).Return([]byte(`"1"`), nil)
				mockClient.On("Do", mock.Anything, "GET", "todos/1", mock.Anything).Return([]byte(`{"ID": "1", "Title": "New Todo"}`), nil)
				return mockClient
			},
			mockRespID:  []byte(`"1"`),
			mockRespGet: []byte(`{"ID": "1", "Title": "New Todo"}`),
			mockErr:     nil,
			expectError: false,
			expectTodo:  &expectedTodo,
		},
		{
			name:  "failed to convert create todo input",
			input: createTodoInput,
			mockConvert: func() *automock.TodoConverter {
				todoConverter := &automock.TodoConverter{}
				todoConverter.EXPECT().ConvertCreateTodoInput(createTodoInput).Return(models.Todo{}, errors.New("conversion error"))
				return todoConverter
			},
			mockDo: func() *mock2.ClientMock {
				return new(mock2.ClientMock)
			},
			mockRespID:  nil,
			mockRespGet: nil,
			mockErr:     nil,
			expectError: true,
			expectTodo:  nil,
		},
		{
			name:  "failed to get todo after creation",
			input: createTodoInput,
			mockConvert: func() *automock.TodoConverter {
				todoConverter := &automock.TodoConverter{}
				todoConverter.EXPECT().ConvertCreateTodoInput(createTodoInput).Return(httpInput, nil)
				todoConverter.EXPECT().ConvertTodoToGraphQL(inputTodo).Return(&expectedTodo, nil)
				return todoConverter
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "POST", "todos/create", data).Return([]byte(`"1"`), nil)
				mockClient.On("Do", mock.Anything, "GET", "todos/1", mock.Anything).Return(data, errors.New("failed to fetch todo"))
				return mockClient
			},
			mockRespID:  []byte(`"1"`),
			mockRespGet: nil,
			mockErr:     errors.New("failed to fetch todo"),
			expectError: true,
			expectTodo:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := tt.mockDo()
			todoConverter := tt.mockConvert()

			r := todo.NewResolver(mockClient, todoConverter, nil, nil)

			result, err := r.CreateTodo(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectTodo, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestUpdateTodo_TodoResolver(t *testing.T) {
	expectedTodo := graphql.Todo{
		ID:    "1",
		Title: "Updated Todo",
	}

	inputTodo := models.Todo{
		ID:    "1",
		Title: "Updated Todo",
	}

	title := "Updated Todo"

	updateTodoInput := graphql.UpdateTodoInput{
		Title: &title,
	}

	httpInput := models.Todo{
		ID:    "1",
		Title: "Updated Todo",
	}

	data, _ := json.Marshal(inputTodo)

	tests := []struct {
		name        string
		id          string
		input       graphql.UpdateTodoInput
		mockConvert func() *automock.TodoConverter
		mockDo      func() *mock2.ClientMock
		mockErr     error
		expectError bool
		expectTodo  *graphql.Todo
	}{
		{
			name:  "successful todo update",
			id:    "1",
			input: updateTodoInput,
			mockConvert: func() *automock.TodoConverter {
				todoConverter := &automock.TodoConverter{}
				todoConverter.EXPECT().ConvertUpdateTodoInput(updateTodoInput).Return(httpInput, nil)
				todoConverter.EXPECT().ConvertTodoToGraphQL(inputTodo).Return(&expectedTodo, nil)
				return todoConverter
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "PUT", "todos/1", data).Return([]byte(`"1"`), nil)
				mockClient.On("Do", mock.Anything, "GET", "todos/1", mock.Anything).Return([]byte(`{"ID": "1", "Title": "Updated Todo"}`), nil)
				return mockClient
			},
			mockErr:     nil,
			expectError: false,
			expectTodo:  &expectedTodo,
		},
		{
			name:  "failed to convert update todo input",
			id:    "1",
			input: updateTodoInput,
			mockConvert: func() *automock.TodoConverter {
				todoConverter := &automock.TodoConverter{}
				todoConverter.EXPECT().ConvertUpdateTodoInput(updateTodoInput).Return(models.Todo{}, errors.New("conversion error"))
				return todoConverter
			},
			mockDo: func() *mock2.ClientMock {
				return new(mock2.ClientMock)
			},
			mockErr:     nil,
			expectError: true,
			expectTodo:  nil,
		},
		{
			name:  "failed to update todo",
			id:    "1",
			input: updateTodoInput,
			mockConvert: func() *automock.TodoConverter {
				todoConverter := &automock.TodoConverter{}
				todoConverter.EXPECT().ConvertUpdateTodoInput(updateTodoInput).Return(httpInput, nil)
				return todoConverter
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "PUT", "todos/1", data).Return([]byte(`"1"`), errors.New("failed to update todo"))
				return mockClient
			},
			mockErr:     errors.New("failed to update todo"),
			expectError: true,
			expectTodo:  nil,
		},
		{
			name:  "failed to fetch todo after update",
			id:    "1",
			input: updateTodoInput,
			mockConvert: func() *automock.TodoConverter {
				todoConverter := &automock.TodoConverter{}
				todoConverter.EXPECT().ConvertUpdateTodoInput(updateTodoInput).Return(httpInput, nil)
				todoConverter.EXPECT().ConvertTodoToGraphQL(inputTodo).Return(&expectedTodo, nil)
				return todoConverter
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "PUT", "todos/1", data).Return([]byte(`"1"`), nil)
				mockClient.On("Do", mock.Anything, "GET", "todos/1", mock.Anything).Return([]byte(`"1"`), errors.New("failed to fetch todo"))
				return mockClient
			},
			mockErr:     errors.New("failed to fetch todo"),
			expectError: true,
			expectTodo:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := tt.mockDo()
			todoConverter := tt.mockConvert()

			r := todo.NewResolver(mockClient, todoConverter, nil, nil)

			result, err := r.UpdateTodo(context.Background(), tt.id, tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectTodo, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestDeleteTodo_TodoResolver(t *testing.T) {
	expectedTodo := graphql.Todo{
		ID:    "1",
		Title: "Test Todo",
	}

	existingTodo := models.Todo{
		ID:    "1",
		Title: "Test Todo",
	}

	tests := []struct {
		name        string
		id          string
		mockConvert func() *automock.TodoConverter
		mockDo      func() *mock2.ClientMock
		mockRespGet []byte
		mockErr     error
		expectError bool
		expectTodo  *graphql.Todo
	}{
		{
			name: "successful todo deletion",
			id:   "1",
			mockConvert: func() *automock.TodoConverter {
				todoConverter := &automock.TodoConverter{}
				todoConverter.EXPECT().ConvertTodoToGraphQL(existingTodo).Return(&expectedTodo, nil)
				return todoConverter
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "GET", "todos/1", mock.Anything).Return([]byte(`{"ID": "1", "Title": "Test Todo"}`), nil)
				mockClient.On("Do", mock.Anything, "DELETE", "todos/1", mock.Anything).Return([]byte(`"1"`), nil)
				return mockClient
			},
			mockRespGet: []byte(`{"ID": "1", "Title": "Test Todo"}`),
			mockErr:     nil,
			expectError: false,
			expectTodo:  &expectedTodo,
		},
		{
			name: "failed to fetch todo before deletion",
			id:   "1",
			mockConvert: func() *automock.TodoConverter {
				return &automock.TodoConverter{}
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "GET", "todos/1", mock.Anything).Return([]byte(`"1"`), errors.New("failed to fetch todo"))
				return mockClient
			},
			mockRespGet: nil,
			mockErr:     errors.New("failed to fetch todo"),
			expectError: true,
			expectTodo:  nil,
		},
		{
			name: "failed to delete todo",
			id:   "1",
			mockConvert: func() *automock.TodoConverter {
				todoConverter := &automock.TodoConverter{}
				todoConverter.EXPECT().ConvertTodoToGraphQL(existingTodo).Return(&expectedTodo, nil)
				return todoConverter
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "GET", "todos/1", mock.Anything).Return([]byte(`{"ID": "1", "Title": "Test Todo"}`), nil)
				mockClient.On("Do", mock.Anything, "DELETE", "todos/1", mock.Anything).Return([]byte(`"1"`), errors.New("failed to delete todo"))
				return mockClient
			},
			mockRespGet: []byte(`{"ID": "1", "Title": "Test Todo"}`),
			mockErr:     errors.New("failed to delete todo"),
			expectError: true,
			expectTodo:  nil,
		},
		{
			name: "failed to unmarshal todo",
			id:   "1",
			mockConvert: func() *automock.TodoConverter {
				return &automock.TodoConverter{}
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "GET", "todos/1", mock.Anything).Return([]byte(`invalid_json`), nil)
				return mockClient
			},
			mockRespGet: []byte(`invalid_json`),
			mockErr:     errors.New("failed to unmarshal todo"),
			expectError: true,
			expectTodo:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := tt.mockDo()
			todoConverter := tt.mockConvert()
			r := todo.NewResolver(mockClient, todoConverter, nil, nil)

			result, err := r.DeleteTodo(context.Background(), tt.id)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectTodo, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
