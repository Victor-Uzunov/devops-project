package http_test

//
//import (
//	"bytes"
//	"encoding/json"
//	"errors"
//	http2 "github.com/Victor-Uzunov/devops-project/todoservice/internal/http"
//	middle "github.com/Victor-Uzunov/devops-project/todoservice/internal/http/automock"
//	"github.com/Victor-Uzunov/devops-project/todoservice/internal/lists/automock"
//	automock2 "github.com/Victor-Uzunov/devops-project/todoservice/internal/todos/automock"
//	automock3 "github.com/Victor-Uzunov/devops-project/todoservice/internal/users/automock"
//	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
//	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
//	"github.com/gorilla/mux"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"github.com/stretchr/testify/require"
//	sqlxmock "github.com/zhashkevych/go-sqlxmock"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestCreateRoutes(t *testing.T) {
//	db, mockDatabase, err := sqlxmock.Newx()
//	require.NoError(t, err)
//	tests := []struct {
//		name                string
//		method              string
//		url                 string
//		input               interface{}
//		expectedStatus      int
//		mockServiceFunction func(listService *automock.ListService, todoService *automock2.TodoService, middleware *middle.Middlewares)
//		mockDatabase        func()
//	}{
//		{
//			name:   "Create List",
//			method: http.MethodPost,
//			url:    "/lists/create",
//			input: &models.List{
//				ID:   "1",
//				Tags: json.RawMessage{0x6e, 0x75, 0x6c, 0x6c},
//			},
//			expectedStatus: http.StatusCreated,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService, middleware *middle.Middlewares) {
//				listService.EXPECT().CreateList(mock.Anything, models.List{ID: "1", Tags: json.RawMessage{0x6e, 0x75, 0x6c, 0x6c}}).Return("TestID", nil).Once()
//				middleware.EXPECT().
//					JWTMiddleware(mock.AnythingOfType("http.HandlerFunc")).
//					Return(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//						next := r.Context().Value("nextHandler").(http.HandlerFunc)
//						next.ServeHTTP(w, r)
//					}))
//
//				middleware.EXPECT().
//					Protected(mock.Anything, mock.Anything, mock.Anything).
//					Return(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//						next := r.Context().Value("nextHandler").(http.HandlerFunc)
//						next.ServeHTTP(w, r)
//					}))
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectCommit()
//			},
//		},
//		{
//			name:   "Error when create list fails",
//			method: http.MethodPost,
//			url:    "/lists/create",
//			input: &models.List{
//				ID:   "1",
//				Tags: json.RawMessage{0x6e, 0x75, 0x6c, 0x6c},
//			},
//			expectedStatus: http.StatusBadRequest,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService, middleware *middle.Middlewares) {
//				listService.EXPECT().CreateList(mock.Anything, models.List{ID: "1", Tags: json.RawMessage{0x6e, 0x75, 0x6c, 0x6c}}).Return("", errors.New("error")).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectRollback()
//			},
//		},
//		{
//			name:   "Create Todo",
//			method: http.MethodPost,
//			url:    "/todos/create",
//			input: &models.List{
//				ID:   "1",
//				Tags: json.RawMessage{0x6e, 0x75, 0x6c, 0x6c},
//			},
//			expectedStatus: http.StatusCreated,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService, middleware *middle.Middlewares) {
//				todoService.EXPECT().CreateTodo(mock.Anything, models.Todo{ID: "1", Tags: json.RawMessage{0x6e, 0x75, 0x6c, 0x6c}}).Return("TestID", nil).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectCommit()
//			},
//		},
//		{
//			name:   "Error when create todo fails",
//			method: http.MethodPost,
//			url:    "/todos/create",
//			input: &models.List{
//				ID:   "1",
//				Tags: json.RawMessage{0x6e, 0x75, 0x6c, 0x6c},
//			},
//			expectedStatus: http.StatusBadRequest,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService, middleware *middle.Middlewares) {
//				todoService.EXPECT().CreateTodo(mock.Anything, models.Todo{ID: "1", Tags: json.RawMessage{0x6e, 0x75, 0x6c, 0x6c}}).Return("", errors.New("error")).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectRollback()
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			listService := new(automock.ListService)
//			todoService := new(automock2.TodoService)
//			userService := new(automock3.UserService)
//			middleware := new(middle.Middlewares)
//			test.mockServiceFunction(listService, todoService, middleware)
//
//			test.mockDatabase()
//
//			req, _ := http.NewRequest(test.method, test.url, func() *bytes.Buffer {
//				if test.input != nil {
//					body, _ := json.Marshal(test.input)
//					return bytes.NewBuffer(body)
//				}
//				return nil
//			}())
//			mockToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4ifQ.UHnffynBjuE3dcwEUyqldVbN-5QzMT-oiyXqkRbWJOI"
//			req.Header.Set("Content-Type", constants.ContentTypeJSON)
//			req.Header.Set("Authorization", "Bearer "+mockToken)
//
//			w := httptest.NewRecorder()
//
//			server := http2.NewServerWithServices(db, listService, todoService, userService, middleware)
//			router := mux.NewRouter()
//			server.RegisterRoutes(router)
//			router.ServeHTTP(w, req)
//
//			assert.Equal(t, test.expectedStatus, w.Code)
//			listService.AssertExpectations(t)
//			todoService.AssertExpectations(t)
//
//			if err := mockDatabase.ExpectationsWereMet(); err != nil {
//				t.Errorf("there were unfulfilled expectations: %s", err)
//			}
//		})
//	}
//}
//
//func TestGetRoutes(t *testing.T) {
//	db, mockDatabase, err := sqlxmock.Newx()
//	require.NoError(t, err)
//	tests := []struct {
//		name                string
//		method              string
//		url                 string
//		expectedStatus      int
//		mockServiceFunction func(listService *automock.ListService, todoService *automock2.TodoService)
//		mockDatabase        func()
//	}{
//		{
//			name:           "Get List",
//			method:         http.MethodGet,
//			url:            "/lists/1",
//			expectedStatus: http.StatusOK,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService) {
//				listService.EXPECT().GetList(mock.Anything, "1").Return(models.List{}, nil).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectCommit()
//			},
//		},
//		{
//			name:           "Error when get list fails",
//			method:         http.MethodGet,
//			url:            "/lists/1",
//			expectedStatus: http.StatusNotFound,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService) {
//				listService.EXPECT().GetList(mock.Anything, "1").Return(models.List{}, errors.New("error")).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectRollback()
//			},
//		},
//		{
//			name:           "Get Todo",
//			method:         http.MethodGet,
//			url:            "/todos/1",
//			expectedStatus: http.StatusOK,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService) {
//				todoService.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{}, nil).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectCommit()
//			},
//		},
//		{
//			name:           "Error when get todo fails",
//			method:         http.MethodGet,
//			url:            "/todos/1",
//			expectedStatus: http.StatusNotFound,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService) {
//				todoService.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{}, errors.New("error")).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectRollback()
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			listService := new(automock.ListService)
//			todoService := new(automock2.TodoService)
//			userService := new(automock3.UserService)
//			middleware := new(middle.Middlewares)
//			test.mockServiceFunction(listService, todoService)
//
//			test.mockDatabase()
//
//			req, _ := http.NewRequest(test.method, test.url, nil)
//			req.Header.Set("Content-Type", constants.ContentTypeJSON)
//
//			w := httptest.NewRecorder()
//
//			server := http2.NewServerWithServices(db, listService, todoService, userService, middleware)
//			router := mux.NewRouter()
//			server.RegisterRoutes(router)
//			router.ServeHTTP(w, req)
//
//			assert.Equal(t, test.expectedStatus, w.Code)
//			listService.AssertExpectations(t)
//			todoService.AssertExpectations(t)
//
//			if err := mockDatabase.ExpectationsWereMet(); err != nil {
//				t.Errorf("there were unfulfilled expectations: %s", err)
//			}
//		})
//	}
//}
//
//func TestUpdateRoutes(t *testing.T) {
//	db, mockDatabase, err := sqlxmock.Newx()
//	require.NoError(t, err)
//	modelTodo := models.Todo{
//		ID:          "1",
//		Title:       "Test",
//		Description: "Test",
//		ListID:      "user1",
//		Tags:        json.RawMessage{0x6e, 0x75, 0x6c, 0x6c},
//	}
//	tests := []struct {
//		name                string
//		method              string
//		url                 string
//		input               interface{}
//		expectedStatus      int
//		mockServiceFunction func(listService *automock.ListService, todoService *automock2.TodoService)
//		mockDatabase        func()
//	}{
//		{
//			name:   "Update List",
//			method: http.MethodPut,
//			url:    "/lists/1",
//			input: &models.List{
//				ID:          "1",
//				Name:        "Test",
//				Description: "Test",
//				OwnerID:     "user1",
//				Tags:        json.RawMessage{0x6e, 0x75, 0x6c, 0x6c},
//			},
//			expectedStatus: http.StatusOK,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService) {
//				listService.EXPECT().GetList(mock.Anything, "1").Return(models.List{}, nil).Once()
//				listService.EXPECT().UpdateList(mock.Anything, models.List{}).Return(nil).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectCommit()
//			},
//		},
//		{
//			name:   "Error when update list fails",
//			method: http.MethodPut,
//			url:    "/lists/1",
//			input: &models.List{
//				ID:          "1",
//				Name:        "Test",
//				Description: "Test",
//				OwnerID:     "user1",
//				Tags:        json.RawMessage{0x6e, 0x75, 0x6c, 0x6c},
//			},
//			expectedStatus: http.StatusNotFound,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService) {
//				listService.EXPECT().GetList(mock.Anything, "1").Return(models.List{}, errors.New("error")).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectRollback()
//			},
//		},
//		{
//			name:           "Update Todo",
//			method:         http.MethodPut,
//			url:            "/todos/1",
//			input:          modelTodo,
//			expectedStatus: http.StatusOK,
//			mockServiceFunction: func(lisService *automock.ListService, todoService *automock2.TodoService) {
//				todoService.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{}, nil).Once()
//				todoService.EXPECT().UpdateTodo(mock.Anything, modelTodo).Return(nil).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectCommit()
//			},
//		},
//		{
//			name:           "Error when update todo fails",
//			method:         http.MethodPut,
//			url:            "/todos/1",
//			input:          modelTodo,
//			expectedStatus: http.StatusNotFound,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService) {
//				todoService.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{}, errors.New("error")).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectRollback()
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			listService := new(automock.ListService)
//			todoService := new(automock2.TodoService)
//			userService := new(automock3.UserService)
//			middleware := new(middle.Middlewares)
//			test.mockServiceFunction(listService, todoService)
//
//			test.mockDatabase()
//
//			req, _ := http.NewRequest(test.method, test.url, func() *bytes.Buffer {
//				if test.input != nil {
//					body, _ := json.Marshal(test.input)
//					return bytes.NewBuffer(body)
//				}
//				return nil
//			}())
//			req.Header.Set("Content-Type", constants.ContentTypeJSON)
//
//			w := httptest.NewRecorder()
//
//			server := http2.NewServerWithServices(db, listService, todoService, userService, middleware)
//			router := mux.NewRouter()
//			server.RegisterRoutes(router)
//			router.ServeHTTP(w, req)
//
//			assert.Equal(t, test.expectedStatus, w.Code)
//			listService.AssertExpectations(t)
//			todoService.AssertExpectations(t)
//
//			if err := mockDatabase.ExpectationsWereMet(); err != nil {
//				t.Errorf("there were unfulfilled expectations: %s", err)
//			}
//		})
//	}
//}
//
//func TestDeleteRoutes(t *testing.T) {
//	db, mockDatabase, err := sqlxmock.Newx()
//	require.NoError(t, err)
//	tests := []struct {
//		name                string
//		method              string
//		url                 string
//		expectedStatus      int
//		mockServiceFunction func(listService *automock.ListService, todoService *automock2.TodoService)
//		mockDatabase        func()
//	}{
//		{
//			name:           "Delete List",
//			method:         http.MethodDelete,
//			url:            "/lists/1",
//			expectedStatus: http.StatusNoContent,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService) {
//				listService.EXPECT().GetList(mock.Anything, "1").Return(models.List{}, nil).Once()
//				listService.EXPECT().DeleteList(mock.Anything, "1").Return(nil).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectCommit()
//			},
//		},
//		{
//			name:           "Error when delete list fails",
//			method:         http.MethodDelete,
//			url:            "/lists/1",
//			expectedStatus: http.StatusNotFound,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService) {
//				listService.EXPECT().GetList(mock.Anything, "1").Return(models.List{}, errors.New("error")).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectRollback()
//			},
//		},
//		{
//			name:           "Delete Todo",
//			method:         http.MethodDelete,
//			url:            "/todos/1",
//			expectedStatus: http.StatusNoContent,
//			mockServiceFunction: func(lisService *automock.ListService, todoService *automock2.TodoService) {
//				todoService.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{}, nil).Once()
//				todoService.EXPECT().DeleteTodo(mock.Anything, "1").Return(nil).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectCommit()
//			},
//		},
//		{
//			name:           "Error when delete todo fails",
//			method:         http.MethodDelete,
//			url:            "/todos/1",
//			expectedStatus: http.StatusNotFound,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService) {
//				todoService.EXPECT().GetTodo(mock.Anything, "1").Return(models.Todo{}, errors.New("error")).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectRollback()
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			listService := new(automock.ListService)
//			todoService := new(automock2.TodoService)
//			userService := new(automock3.UserService)
//			middleware := new(middle.Middlewares)
//			test.mockServiceFunction(listService, todoService)
//
//			test.mockDatabase()
//
//			req, _ := http.NewRequest(test.method, test.url, nil)
//			req.Header.Set("Content-Type", constants.ContentTypeJSON)
//
//			w := httptest.NewRecorder()
//
//			server := http2.NewServerWithServices(db, listService, todoService, userService, middleware)
//			router := mux.NewRouter()
//			server.RegisterRoutes(router)
//			router.ServeHTTP(w, req)
//
//			assert.Equal(t, test.expectedStatus, w.Code)
//			listService.AssertExpectations(t)
//			todoService.AssertExpectations(t)
//
//			if err := mockDatabase.ExpectationsWereMet(); err != nil {
//				t.Errorf("there were unfulfilled expectations: %s", err)
//			}
//		})
//	}
//}
//
//func TestListAllRoutes(t *testing.T) {
//	db, mockDatabase, err := sqlxmock.Newx()
//	require.NoError(t, err)
//	tests := []struct {
//		name                string
//		method              string
//		url                 string
//		expectedStatus      int
//		mockServiceFunction func(listService *automock.ListService, todoService *automock2.TodoService)
//		mockDatabase        func()
//	}{
//		{
//			name:           "List all lists by userID",
//			method:         http.MethodGet,
//			url:            "/users/1/lists",
//			expectedStatus: http.StatusOK,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService) {
//				listService.EXPECT().ListAllByUserID(mock.Anything, "1").Return([]models.Access{}, nil).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectCommit()
//			},
//		},
//		{
//			name:           "Error when list all lists by user fails",
//			method:         http.MethodGet,
//			url:            "/users/1/lists",
//			expectedStatus: http.StatusNotFound,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService) {
//				listService.EXPECT().ListAllByUserID(mock.Anything, "1").Return(nil, errors.New("error")).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectRollback()
//			},
//		},
//		{
//			name:           "List todos by listID",
//			method:         http.MethodGet,
//			url:            "/lists/1/todos",
//			expectedStatus: http.StatusOK,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService) {
//				todoService.EXPECT().ListTodosByListID(mock.Anything, "1").Return([]models.Todo{}, nil).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectCommit()
//			},
//		},
//		{
//			name:           "Error when list todos by list fails",
//			method:         http.MethodGet,
//			url:            "/lists/1/todos",
//			expectedStatus: http.StatusNotFound,
//			mockServiceFunction: func(listService *automock.ListService, todoService *automock2.TodoService) {
//				todoService.EXPECT().ListTodosByListID(mock.Anything, "1").Return(nil, errors.New("error")).Once()
//			},
//			mockDatabase: func() {
//				mockDatabase.ExpectBegin()
//				mockDatabase.ExpectRollback()
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			listService := new(automock.ListService)
//			todoService := new(automock2.TodoService)
//			userService := new(automock3.UserService)
//			middleware := new(middle.Middlewares)
//			test.mockServiceFunction(listService, todoService)
//
//			test.mockDatabase()
//
//			req, _ := http.NewRequest(test.method, test.url, nil)
//			req.Header.Set("Content-Type", constants.ContentTypeJSON)
//
//			w := httptest.NewRecorder()
//
//			server := http2.NewServerWithServices(db, listService, todoService, userService, middleware)
//			router := mux.NewRouter()
//			server.RegisterRoutes(router)
//			router.ServeHTTP(w, req)
//
//			assert.Equal(t, test.expectedStatus, w.Code)
//			listService.AssertExpectations(t)
//			todoService.AssertExpectations(t)
//
//			if err := mockDatabase.ExpectationsWereMet(); err != nil {
//				t.Errorf("there were unfulfilled expectations: %s", err)
//			}
//		})
//	}
//}
