package todo_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/http/todo"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/todos/automock"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTodoHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	err = errors.New("error")
	id := "1"
	modelInput := models.Todo{
		ID:          id,
		Title:       "Test Todo",
		Description: "Test Todo",
		ListID:      "list1",
		Tags:        json.RawMessage{0x6e, 0x75, 0x6c, 0x6c},
	}

	tests := []struct {
		name               string
		mockService        func() *automock.TodoService
		mockDatabase       func()
		expectedStatusCode int
		expectedError      error
	}{
		{
			name: "Create Todo",
			mockService: func() *automock.TodoService {
				mockService := &automock.TodoService{}
				mockService.EXPECT().CreateTodo(mock.Anything, modelInput).Return(id, nil).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectCommit()
			},
			expectedStatusCode: http.StatusCreated,
			expectedError:      nil,
		},
		{
			name: "Error when create todo fails",
			mockService: func() *automock.TodoService {
				mockService := &automock.TodoService{}
				mockService.EXPECT().CreateTodo(mock.Anything, modelInput).Return("", err).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectRollback()
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.mockService()
			handler := todo.NewHandler(mockService, db)

			body, _ := json.Marshal(modelInput)
			req, _ := http.NewRequest(http.MethodPost, "/todos/create", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			w := httptest.NewRecorder()
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.CreateTodo(w, req)
			resp := w.Result()
			defer func() {
				err := resp.Body.Close()
				if err != nil {
					t.Error(err)
				}
			}()

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				var respID string
				_ = json.NewDecoder(resp.Body).Decode(&respID)
				assert.Equal(t, "1", respID)
			}

			if err := mockDatabase.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetTodoHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	err = errors.New("error")
	id := "1"
	modelInput := models.Todo{
		ID:          id,
		Title:       "Test Todo",
		Description: "Test Todo",
		ListID:      "list1",
	}
	model := models.Todo{
		ID:          id,
		Title:       "Test Todo",
		Description: "Test Todo",
		ListID:      "list1",
	}
	tests := []struct {
		name               string
		mockService        func() *automock.TodoService
		mockDatabase       func()
		expectedStatusCode int
		urlVars            map[string]string
		expectedError      error
	}{
		{
			name: "Get List",
			mockService: func() *automock.TodoService {
				mockService := &automock.TodoService{}
				mockService.EXPECT().GetTodo(mock.Anything, id).Return(model, nil).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			urlVars:            map[string]string{"id": id},
			expectedError:      nil,
		},
		{
			name: "Error when get todo fails",
			mockService: func() *automock.TodoService {
				mockService := &automock.TodoService{}
				mockService.EXPECT().GetTodo(mock.Anything, id).Return(models.Todo{}, err).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
			urlVars:            map[string]string{"id": id},
			expectedError:      err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.mockService()
			handler := todo.NewHandler(mockService, db)

			body, _ := json.Marshal(modelInput)
			req, _ := http.NewRequest(http.MethodGet, "/todos/1", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			w := httptest.NewRecorder()
			req = mux.SetURLVars(req, tt.urlVars)
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.GetTodo(w, req)
			resp := w.Result()
			defer func() {
				err := resp.Body.Close()
				if err != nil {
					t.Error(err)
				}
			}()

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				expectedResponse, _ := json.Marshal(model)
				var actualResponse bytes.Buffer
				if _, err := actualResponse.ReadFrom(resp.Body); err != nil {
					t.Error(err)
				}
				assert.JSONEq(t, string(expectedResponse), actualResponse.String())
			}

			if err := mockDatabase.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateTodoHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	err = errors.New("error")
	id := ""
	modelInput := models.Todo{
		ID:          id,
		Title:       "Test Todo",
		Description: "Test Todo",
		ListID:      "list1",
		Tags:        json.RawMessage{0x6e, 0x75, 0x6c, 0x6c},
	}
	tests := []struct {
		name               string
		mockService        func() *automock.TodoService
		mockDatabase       func()
		expectedStatusCode int
		urlVars            map[string]string
		expectedError      error
	}{
		{
			name: "Update Todo",
			mockService: func() *automock.TodoService {
				mockService := &automock.TodoService{}
				mockService.EXPECT().GetTodo(mock.Anything, id).Return(modelInput, nil).Once()
				mockService.EXPECT().UpdateTodo(mock.Anything, modelInput).Return(nil).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			urlVars:            map[string]string{"id": id},
			expectedError:      nil,
		},
		{
			name: "Error when update todo fails",
			mockService: func() *automock.TodoService {
				mockService := &automock.TodoService{}
				mockService.EXPECT().GetTodo(mock.Anything, id).Return(models.Todo{}, err).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
			urlVars:            map[string]string{"id": id},
			expectedError:      err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.mockService()
			handler := todo.NewHandler(mockService, db)

			body, _ := json.Marshal(modelInput)
			req, _ := http.NewRequest(http.MethodPut, "/todos/update/1", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			w := httptest.NewRecorder()
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.UpdateTodo(w, req)
			resp := w.Result()
			defer func() {
				err := resp.Body.Close()
				if err != nil {
					t.Error(err)
				}
			}()

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				var response string
				_ = json.NewDecoder(resp.Body).Decode(&response)
				assert.Equal(t, "", response)
			}

			if err := mockDatabase.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteTodoHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	err = errors.New("error")
	id := "1"

	tests := []struct {
		name               string
		mockService        func() *automock.TodoService
		mockDatabase       func()
		expectedStatusCode int
		urlVars            map[string]string
		expectedError      error
	}{
		{
			name: "Delete Todo",
			mockService: func() *automock.TodoService {
				mockService := &automock.TodoService{}
				mockService.EXPECT().GetTodo(mock.Anything, id).Return(models.Todo{}, nil).Once()
				mockService.EXPECT().DeleteTodo(mock.Anything, id).Return(nil).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectCommit()
			},
			expectedStatusCode: http.StatusNoContent,
			urlVars:            map[string]string{"id": id},
			expectedError:      nil,
		},
		{
			name: "Error when delete todo fails",
			mockService: func() *automock.TodoService {
				mockService := &automock.TodoService{}
				mockService.EXPECT().GetTodo(mock.Anything, id).Return(models.Todo{}, err).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
			urlVars:            map[string]string{"id": id},
			expectedError:      err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.mockService()
			handler := todo.NewHandler(mockService, db)

			req, _ := http.NewRequest(http.MethodDelete, "/todos/delete/1", nil)
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			req = mux.SetURLVars(req, tt.urlVars)
			w := httptest.NewRecorder()
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.DeleteTodo(w, req)
			resp := w.Result()
			defer func() {
				err := resp.Body.Close()
				if err != nil {
					t.Error(err)
				}
			}()

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				var respID string
				_ = json.NewDecoder(resp.Body).Decode(&respID)
				assert.Equal(t, "", respID)
			}

			if err := mockDatabase.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestListTodosByListIDHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	err = errors.New("error")
	id := "1"
	model := models.Todo{
		ID:          id,
		Title:       "Test List",
		Description: "Test List",
		ListID:      "list1",
	}
	tests := []struct {
		name               string
		mockService        func() *automock.TodoService
		mockDatabase       func()
		expectedStatusCode int
		urlVars            map[string]string
		expectedError      error
	}{
		{
			name: "List all by listID",
			mockService: func() *automock.TodoService {
				mockService := &automock.TodoService{}
				mockService.EXPECT().ListTodosByListID(mock.Anything, id).Return([]models.Todo{model}, nil).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			urlVars:            map[string]string{"list_id": id},
			expectedError:      nil,
		},
		{
			name: "Error when list all by listID fails",
			mockService: func() *automock.TodoService {
				mockService := &automock.TodoService{}
				mockService.EXPECT().ListTodosByListID(mock.Anything, id).Return([]models.Todo{}, err).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
			urlVars:            map[string]string{"list_id": id},
			expectedError:      err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.mockService()
			handler := todo.NewHandler(mockService, db)

			req, _ := http.NewRequest(http.MethodGet, "/lists/1/todos", nil)
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			req = mux.SetURLVars(req, tt.urlVars)
			w := httptest.NewRecorder()
			req = mux.SetURLVars(req, tt.urlVars)
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.ListTodosByListID(w, req)
			resp := w.Result()
			defer func() {
				err := resp.Body.Close()
				if err != nil {
					t.Error(err)
				}
			}()

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				expectedResponse, _ := json.Marshal([]models.Todo{model})
				var actualResponse bytes.Buffer
				if _, err := actualResponse.ReadFrom(resp.Body); err != nil {
					t.Error(err)
				}
				assert.JSONEq(t, string(expectedResponse), actualResponse.String())
			}

			if err := mockDatabase.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
