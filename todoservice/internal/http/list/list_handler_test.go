package list_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/http/list"
	"github.com/Victor-Uzunov/devops-project/todoservice/internal/lists/automock"
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

func TestCreateListHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	err = errors.New("error")
	id := "1"
	modelInput := models.List{
		ID:          id,
		Name:        "Test List",
		Description: "Test List",
		OwnerID:     "user1",
		Tags:        json.RawMessage{0x6e, 0x75, 0x6c, 0x6c},
	}

	tests := []struct {
		name               string
		mockService        func() *automock.ListService
		mockDatabase       func()
		expectedStatusCode int
		expectedError      error
	}{
		{
			name: "Create List",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().CreateList(mock.Anything, modelInput).Return(id, nil).Once()
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
			name: "Error when create list fails",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().CreateList(mock.Anything, modelInput).Return("", err).Once()
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
			handler := list.NewHandler(mockService, db)

			body, _ := json.Marshal(modelInput)
			req, _ := http.NewRequest(http.MethodPost, "/lists/create", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			w := httptest.NewRecorder()
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.CreateList(w, req)
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

func TestGetListHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	err = errors.New("error")
	id := "1"
	modelInput := models.List{
		ID:          id,
		Name:        "Test List",
		Description: "Test List",
		OwnerID:     "user1",
	}
	model := models.List{
		ID:          id,
		Name:        "Test List",
		Description: "Test List",
		OwnerID:     "user1",
	}
	tests := []struct {
		name               string
		mockService        func() *automock.ListService
		mockDatabase       func()
		expectedStatusCode int
		urlVars            map[string]string
		expectedError      error
	}{
		{
			name: "Get List",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().GetList(mock.Anything, id).Return(model, nil).Once()
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
			name: "Error when get list fails",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().GetList(mock.Anything, id).Return(models.List{}, err).Once()
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
			handler := list.NewHandler(mockService, db)

			body, _ := json.Marshal(modelInput)
			req, _ := http.NewRequest(http.MethodGet, "/lists/1", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			w := httptest.NewRecorder()
			req = mux.SetURLVars(req, tt.urlVars)
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.GetList(w, req)
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

func TestGetAllListsHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	err = errors.New("error")

	model := []models.List{
		{
			ID:          "1",
			Name:        "Test List",
			Description: "Test List",
			OwnerID:     "user1",
		},
	}
	tests := []struct {
		name               string
		mockService        func() *automock.ListService
		mockDatabase       func()
		expectedStatusCode int
		expectedError      error
	}{
		{
			name: "Get All List",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().GetAllLists(mock.Anything).Return(model, nil).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			expectedError:      nil,
		},
		{
			name: "Error when get all lists fails",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().GetAllLists(mock.Anything).Return([]models.List{}, err).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
			expectedError:      err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.mockService()
			handler := list.NewHandler(mockService, db)

			req, _ := http.NewRequest(http.MethodGet, "/lists/all", nil)
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			w := httptest.NewRecorder()
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.GetAllLists(w, req)
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

func TestGetListOwnerIDHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	err = errors.New("error")
	id := "1"
	modelInput := models.List{
		ID:          id,
		Name:        "Test List",
		Description: "Test List",
		OwnerID:     "user1",
	}
	ownerID := "user1"
	tests := []struct {
		name               string
		mockService        func() *automock.ListService
		mockDatabase       func()
		expectedStatusCode int
		urlVars            map[string]string
		expectedError      error
	}{
		{
			name: "Get list ownerID",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().GetListOwnerID(mock.Anything, id).Return(ownerID, nil).Once()
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
			name: "Error when get list ownerID fails",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().GetListOwnerID(mock.Anything, id).Return("", err).Once()
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
			handler := list.NewHandler(mockService, db)

			body, _ := json.Marshal(modelInput)
			req, _ := http.NewRequest(http.MethodGet, "/lists/1/owner", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			w := httptest.NewRecorder()
			req = mux.SetURLVars(req, tt.urlVars)
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.GetListOwnerID(w, req)
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
				expectedResponse, _ := json.Marshal(ownerID)
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

func TestGetUsersByListIDHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	err = errors.New("error")
	id := "1"
	modelInput := models.List{
		ID:          id,
		Name:        "Test List",
		Description: "Test List",
		OwnerID:     "user1",
		SharedWith:  []string{"user1", "user2"},
	}
	model := []models.Access{
		{
			ListID: id,
			UserID: "user1",
			Role:   "reader",
		},
		{
			ListID: id,
			UserID: "user2",
			Role:   "reader",
		},
	}
	tests := []struct {
		name               string
		mockService        func() *automock.ListService
		mockDatabase       func()
		expectedStatusCode int
		urlVars            map[string]string
		expectedError      error
	}{
		{
			name: "Get users by listID",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().GetUsersByListID(mock.Anything, id).Return(model, nil).Once()
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
			name: "Error when get users by listID fails",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().GetUsersByListID(mock.Anything, id).Return([]models.Access{}, err).Once()
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
			handler := list.NewHandler(mockService, db)

			body, _ := json.Marshal(modelInput)
			req, _ := http.NewRequest(http.MethodGet, "/lists/1", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			w := httptest.NewRecorder()
			req = mux.SetURLVars(req, tt.urlVars)
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.GetUsersByListID(w, req)
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

func TestUpdateListHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	err = errors.New("error")
	id := ""
	modelInput := models.List{
		ID:          id,
		Name:        "Test List",
		Description: "Test List",
		OwnerID:     "user1",
		Tags:        json.RawMessage(`"null"`),
	}
	tests := []struct {
		name               string
		mockService        func() *automock.ListService
		mockDatabase       func()
		expectedStatusCode int
		urlVars            map[string]string
		expectedError      error
	}{
		{
			name: "Update List",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().GetList(mock.Anything, id).Return(modelInput, nil).Once()
				mockService.EXPECT().UpdateList(mock.Anything, mock.MatchedBy(func(input models.List) bool {
					return input.Name == modelInput.Name &&
						input.Description == modelInput.Description
				})).Return(nil).Once()
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
			name: "Error when update list fails",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().GetList(mock.Anything, id).Return(models.List{}, err).Once()
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
			handler := list.NewHandler(mockService, db)

			body, _ := json.Marshal(modelInput)
			req, _ := http.NewRequest(http.MethodPut, "/lists/update/1", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			w := httptest.NewRecorder()
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.UpdateList(w, req)
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

func TestDeleteListHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	err = errors.New("error")
	id := "1"

	tests := []struct {
		name               string
		mockService        func() *automock.ListService
		mockDatabase       func()
		expectedStatusCode int
		urlVars            map[string]string
		expectedError      error
	}{
		{
			name: "Delete List",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().GetList(mock.Anything, id).Return(models.List{}, nil).Once()
				mockService.EXPECT().DeleteList(mock.Anything, id).Return(nil).Once()
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
			name: "Error when delete list fails",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().GetList(mock.Anything, id).Return(models.List{}, err).Once()
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
			handler := list.NewHandler(mockService, db)

			req, _ := http.NewRequest(http.MethodDelete, "/lists/delete/1", nil)
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			req = mux.SetURLVars(req, tt.urlVars)
			w := httptest.NewRecorder()
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.DeleteList(w, req)
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

func TestListAllByUserIDHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	err = errors.New("error")
	id := "1"
	accessModel := models.Access{
		ListID: id,
		UserID: "user1",
	}
	tests := []struct {
		name               string
		mockService        func() *automock.ListService
		mockDatabase       func()
		expectedStatusCode int
		urlVars            map[string]string
		expectedError      error
	}{
		{
			name: "List all by userID",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().ListAllByUserID(mock.Anything, id).Return([]models.Access{accessModel}, nil).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			urlVars:            map[string]string{"user_id": id},
			expectedError:      nil,
		},
		{
			name: "Error when list all by userID fails",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().ListAllByUserID(mock.Anything, id).Return([]models.Access{}, err).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
			urlVars:            map[string]string{"user_id": id},
			expectedError:      err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.mockService()
			handler := list.NewHandler(mockService, db)

			req, _ := http.NewRequest(http.MethodGet, "/users/1/lists", nil)
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			req = mux.SetURLVars(req, tt.urlVars)
			w := httptest.NewRecorder()
			req = mux.SetURLVars(req, tt.urlVars)
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.ListAllByUser(w, req)
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
				expectedResponse, _ := json.Marshal([]models.Access{accessModel})
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

func TestCreateAccessHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	err = errors.New("error")
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
		name               string
		mockService        func() *automock.ListService
		mockDatabase       func()
		expectedStatusCode int
		expectedError      error
	}{
		{
			name: "Create Access",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().CreateAccess(mock.Anything, modelInput).Return(model, nil).Once()
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
			name: "Error when create access fails",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().CreateAccess(mock.Anything, modelInput).Return(models.Access{}, err).Once()
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
			handler := list.NewHandler(mockService, db)

			body, _ := json.Marshal(modelInput)
			req, _ := http.NewRequest(http.MethodPost, "/lists_access/create", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			w := httptest.NewRecorder()
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.CreateAccess(w, req)
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
				var respID models.Access
				_ = json.NewDecoder(resp.Body).Decode(&respID)
				assert.Equal(t, model, respID)
			}
			if err := mockDatabase.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetAccessHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	err = errors.New("error")
	listID := "listID"
	userID := "userID"
	model := models.Access{
		ListID: listID,
		UserID: userID,
		Role:   "reader",
	}
	tests := []struct {
		name               string
		mockService        func() *automock.ListService
		mockDatabase       func()
		expectedStatusCode int
		urlVars            map[string]string
		expectedError      error
	}{
		{
			name: "Get Access",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().GetAccess(mock.Anything, listID, userID).Return(model, nil).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectCommit()
			},
			expectedStatusCode: http.StatusOK,
			urlVars:            map[string]string{"list_id": listID, "user_id": userID},
			expectedError:      nil,
		},
		{
			name: "Error when get access fails",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().GetAccess(mock.Anything, listID, userID).Return(models.Access{}, err).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectRollback()
			},
			expectedStatusCode: http.StatusNotFound,
			urlVars:            map[string]string{"list_id": listID, "user_id": userID},
			expectedError:      err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.mockService()
			handler := list.NewHandler(mockService, db)

			req, _ := http.NewRequest(http.MethodGet, "/lists/listID/userID", nil)
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			w := httptest.NewRecorder()
			req = mux.SetURLVars(req, tt.urlVars)
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.GetAccess(w, req)
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

func TestDeleteAccessHandler(t *testing.T) {
	db, mockDatabase, err := sqlxmock.Newx()
	require.NoError(t, err)
	err = errors.New("error")
	listID := "listID"
	userID := "userID"

	tests := []struct {
		name               string
		mockService        func() *automock.ListService
		mockDatabase       func()
		expectedStatusCode int
		urlVars            map[string]string
		expectedError      error
	}{
		{
			name: "Delete Access",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().DeleteAccess(mock.Anything, listID, userID).Return(nil).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectCommit()
			},
			expectedStatusCode: http.StatusNoContent,
			urlVars:            map[string]string{"list_id": listID, "user_id": userID},
			expectedError:      nil,
		},
		{
			name: "Error when delete list fails",
			mockService: func() *automock.ListService {
				mockService := &automock.ListService{}
				mockService.EXPECT().DeleteAccess(mock.Anything, listID, userID).Return(err).Once()
				return mockService
			},
			mockDatabase: func() {
				mockDatabase.ExpectBegin()
				mockDatabase.ExpectRollback()
			},
			expectedStatusCode: http.StatusBadRequest,
			urlVars:            map[string]string{"list_id": listID, "user_id": userID},
			expectedError:      err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.mockService()
			handler := list.NewHandler(mockService, db)

			req, _ := http.NewRequest(http.MethodDelete, "/lists_access/listID/userID", nil)
			req.Header.Set("Content-Type", constants.ContentTypeJSON)
			req = mux.SetURLVars(req, tt.urlVars)
			w := httptest.NewRecorder()
			defer mock.AssertExpectationsForObjects(t, mockService)

			tt.mockDatabase()

			handler.DeleteAccess(w, req)
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
