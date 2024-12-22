package user_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/http/user"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/users/automock"
	constants2 "github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
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

func TestCreateUserHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	err = errors.New("error")
	id := "1"
	modelInput := models.User{
		ID:       id,
		Email:    "test@example.com",
		Role:     constants2.Reader,
		GithubID: "github-id",
	}

	tests := []struct {
		name               string
		mockService        func() *automock.UserService
		mockDatabase       func()
		expectedStatusCode int
		expectedError      error
	}{
		{
			name: "Create User",
			mockService: func() *automock.UserService {
				mockService := &automock.UserService{}
				mockService.EXPECT().CreateUser(mock.Anything, modelInput).Return(id, nil).Once()
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
			name: "Error when create user fails",
			mockService: func() *automock.UserService {
				mockService := &automock.UserService{}
				mockService.EXPECT().CreateUser(mock.Anything, modelInput).Return("", err).Once()
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
			handler := user.NewHandler(mockService, db)

			body, _ := json.Marshal(modelInput)
			req, _ := http.NewRequest(http.MethodPost, "/users/create", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants2.ContentTypeJSON)
			w := httptest.NewRecorder()
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.CreateUser(w, req)
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
				assert.Equal(t, id, respID)
			}
			if err := mockDatabase.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetUserHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	err = errors.New("error")
	id := "1"
	model := models.User{
		ID:       id,
		Email:    "test",
		GithubID: "github-id",
		Role:     constants2.Reader,
	}

	tests := []struct {
		name               string
		mockService        func() *automock.UserService
		mockDatabase       func()
		expectedStatusCode int
		urlVars            map[string]string
		expectedError      error
	}{
		{
			name: "Get User",
			mockService: func() *automock.UserService {
				mockService := &automock.UserService{}
				mockService.EXPECT().GetUser(mock.Anything, id).Return(model, nil).Once()
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
			name: "Error when get user fails",
			mockService: func() *automock.UserService {
				mockService := &automock.UserService{}
				mockService.EXPECT().GetUser(mock.Anything, id).Return(models.User{}, err).Once()
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
			handler := user.NewHandler(mockService, db)

			req, _ := http.NewRequest(http.MethodGet, "/users/"+id, nil)
			req.Header.Set("Content-Type", constants2.ContentTypeJSON)
			w := httptest.NewRecorder()
			req = mux.SetURLVars(req, tt.urlVars)
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.GetUser(w, req)
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

func TestUpdateUserHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	err = errors.New("error")
	id := "1"
	modelInput := models.User{
		ID:       id,
		Email:    "test",
		GithubID: "github-id",
		Role:     constants2.Reader,
	}

	tests := []struct {
		name               string
		mockService        func() *automock.UserService
		mockDatabase       func()
		expectedStatusCode int
		urlVars            map[string]string
		expectedError      error
	}{
		{
			name: "Update User",
			mockService: func() *automock.UserService {
				mockService := &automock.UserService{}
				mockService.EXPECT().GetUser(mock.Anything, id).Return(modelInput, nil).Once()
				mockService.EXPECT().UpdateUser(mock.Anything, modelInput).Return(nil).Once()
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
			name: "Error when update user fails",
			mockService: func() *automock.UserService {
				mockService := &automock.UserService{}
				mockService.EXPECT().GetUser(mock.Anything, id).Return(models.User{}, err).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
			},
			expectedStatusCode: http.StatusNotFound,
			urlVars:            map[string]string{"id": id},
			expectedError:      err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.mockService()
			handler := user.NewHandler(mockService, db)

			body, _ := json.Marshal(modelInput)
			req, _ := http.NewRequest(http.MethodPut, "/users/update/"+id, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants2.ContentTypeJSON)
			w := httptest.NewRecorder()
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.UpdateUser(w, req)
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

func TestDeleteUserHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	err = errors.New("error")
	id := "1"

	tests := []struct {
		name               string
		mockService        func() *automock.UserService
		mockDatabase       func()
		expectedStatusCode int
		urlVars            map[string]string
		expectedError      error
	}{
		{
			name: "Delete User",
			mockService: func() *automock.UserService {
				mockService := &automock.UserService{}
				mockService.EXPECT().GetUser(mock.Anything, id).Return(models.User{}, nil).Once()
				mockService.EXPECT().DeleteUser(mock.Anything, id).Return(nil).Once()
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
			name: "Error when delete user fails",
			mockService: func() *automock.UserService {
				mockService := &automock.UserService{}
				mockService.EXPECT().GetUser(mock.Anything, id).Return(models.User{}, err).Once()
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
			handler := user.NewHandler(mockService, db)

			req, _ := http.NewRequest(http.MethodDelete, "/users/delete/"+id, nil)
			req.Header.Set("Content-Type", constants2.ContentTypeJSON)
			req = mux.SetURLVars(req, tt.urlVars)
			w := httptest.NewRecorder()
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.DeleteUser(w, req)
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
