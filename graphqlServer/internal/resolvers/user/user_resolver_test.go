package user_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/generated/graphql"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/converters/automock"
	mock2 "github.com/Victor-Uzunov/devops-project/graphqlServer/internal/resolvers/mock"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/resolvers/user"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsers_UserResolver(t *testing.T) {
	expectedUser := graphql.User{
		ID:    "1",
		Email: "victor",
	}
	inputUser := models.User{
		ID:    "1",
		Email: "victor",
	}
	tests := []struct {
		name          string
		mockMethod    string
		mockURL       string
		mockBody      []byte
		mockResp      []byte
		mockErr       error
		expectError   bool
		expectUsers   []*graphql.User
		userConverter func() *automock.UserConverter
		listConverter func() *automock.ListConverter
	}{
		{
			name:        "successful users fetch",
			mockMethod:  "GET",
			mockURL:     "users/all",
			mockResp:    []byte(`[{"ID": "1", "Email": "victor"}]`),
			mockErr:     nil,
			expectError: false,
			expectUsers: []*graphql.User{
				{ID: "1", Email: "victor"},
			},
			userConverter: func() *automock.UserConverter {
				userConverter := &automock.UserConverter{}
				userConverter.EXPECT().ConvertUserToGraphQL(inputUser).Return(&expectedUser, nil)
				return userConverter
			},
			listConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
		},
		{
			name:        "failed http request",
			mockMethod:  "GET",
			mockURL:     "users/all",
			mockResp:    nil,
			mockErr:     errors.New("failed to fetch users"),
			expectError: true,
			expectUsers: nil,
			userConverter: func() *automock.UserConverter {
				return &automock.UserConverter{}
			},
			listConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
		},
		{
			name:        "failed to unmarshal response",
			mockMethod:  "GET",
			mockURL:     "users/all",
			mockResp:    []byte(`invalid JSON`),
			mockErr:     nil,
			expectError: true,
			expectUsers: nil,
			userConverter: func() *automock.UserConverter {
				return &automock.UserConverter{}
			},
			listConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(mock2.ClientMock)

			mockClient.On("Do", mock.Anything, tt.mockMethod, tt.mockURL, mock.Anything).Return(tt.mockResp, tt.mockErr)
			userConverter := tt.userConverter()
			listConverter := tt.listConverter()
			r := user.NewResolver(mockClient, userConverter, listConverter)

			result, err := r.Users(context.Background())

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectUsers, result)
			}

			mockClient.AssertCalled(t, "Do", mock.Anything, tt.mockMethod, tt.mockURL, mock.Anything)
		})
	}
}

func TestUser_UserResolver(t *testing.T) {
	expectedUser := graphql.User{
		ID:    "1",
		Email: "victor",
	}

	inputUser := models.User{
		ID:    "1",
		Email: "victor",
	}

	tests := []struct {
		name          string
		userID        string
		mockMethod    string
		mockURL       string
		mockResp      []byte
		mockErr       error
		expectError   bool
		expectUser    *graphql.User
		userConverter func() *automock.UserConverter
		listConverter func() *automock.ListConverter
	}{
		{
			name:        "successful user fetch",
			userID:      "1",
			mockMethod:  "GET",
			mockURL:     "users/1",
			mockResp:    []byte(`{"ID": "1", "Email": "victor"}`),
			mockErr:     nil,
			expectError: false,
			expectUser:  &expectedUser,
			userConverter: func() *automock.UserConverter {
				userConverter := &automock.UserConverter{}
				userConverter.EXPECT().ConvertUserToGraphQL(inputUser).Return(&expectedUser, nil)
				return userConverter
			},
			listConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
		},
		{
			name:        "failed http request",
			userID:      "1",
			mockMethod:  "GET",
			mockURL:     "users/1",
			mockResp:    nil,
			mockErr:     errors.New("failed to fetch user"),
			expectError: true,
			expectUser:  nil,
			userConverter: func() *automock.UserConverter {
				return &automock.UserConverter{}
			},
			listConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
		},
		{
			name:        "failed to unmarshal response",
			userID:      "1",
			mockMethod:  "GET",
			mockURL:     "users/1",
			mockResp:    []byte(`invalid JSON`),
			mockErr:     nil,
			expectError: true,
			expectUser:  nil,
			userConverter: func() *automock.UserConverter {
				return &automock.UserConverter{}
			},
			listConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
		},
		{
			name:        "failed to convert user",
			userID:      "1",
			mockMethod:  "GET",
			mockURL:     "users/1",
			mockResp:    []byte(`{"ID": "1", "Email": "victor"}`),
			mockErr:     nil,
			expectError: true,
			expectUser:  nil,
			userConverter: func() *automock.UserConverter {
				userConverter := &automock.UserConverter{}
				userConverter.EXPECT().ConvertUserToGraphQL(inputUser).Return(nil, errors.New("conversion error"))
				return userConverter
			},
			listConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(mock2.ClientMock)

			mockClient.On("Do", mock.Anything, tt.mockMethod, tt.mockURL, mock.Anything).Return(tt.mockResp, tt.mockErr)

			userConverter := tt.userConverter()
			listConverter := tt.listConverter()

			r := user.NewResolver(mockClient, userConverter, listConverter)

			result, err := r.User(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectUser, result)
			}

			mockClient.AssertCalled(t, "Do", mock.Anything, tt.mockMethod, tt.mockURL, mock.Anything)
		})
	}
}

func TestCreateUser_UserResolver(t *testing.T) {
	expectedUser := graphql.User{
		ID:    "1",
		Email: "victor",
	}

	inputUser := models.User{
		Email: "victor",
	}

	convertModel := models.User{
		ID:    "1",
		Email: "victor",
	}

	createUserInput := graphql.CreateUserInput{
		Email: "victor",
	}

	httpInput := models.User{
		Email: "victor",
	}

	data, _ := json.Marshal(inputUser)

	tests := []struct {
		name              string
		input             graphql.CreateUserInput
		mockConvert       func() *automock.UserConverter
		mockListConverter func() *automock.ListConverter
		mockDo            func() *mock2.ClientMock
		mockRespID        []byte
		mockRespGet       []byte
		mockErr           error
		expectError       bool
		expectUser        *graphql.User
	}{
		{
			name:  "successful user creation",
			input: createUserInput,
			mockConvert: func() *automock.UserConverter {
				userConverter := &automock.UserConverter{}
				userConverter.EXPECT().ConvertCreateUserInput(createUserInput).Return(httpInput, nil)
				userConverter.EXPECT().ConvertUserToGraphQL(convertModel).Return(&expectedUser, nil)
				return userConverter
			},
			mockListConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "POST", "users/create", data).Return([]byte(`"1"`), nil)
				mockClient.On("Do", mock.Anything, "GET", "users/1", mock.Anything).Return([]byte(`{"ID": "1", "Email": "victor"}`), nil)
				return mockClient
			},
			mockRespID:  []byte(`"1"`),
			mockRespGet: []byte(`"ID": "1", Email": "victor"}`),
			mockErr:     nil,
			expectError: false,
			expectUser:  &expectedUser,
		},
		{
			name:  "failed to convert create user input",
			input: createUserInput,
			mockConvert: func() *automock.UserConverter {
				userConverter := &automock.UserConverter{}
				userConverter.EXPECT().ConvertCreateUserInput(createUserInput).Return(models.User{}, errors.New("conversion error"))
				return userConverter
			},
			mockListConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
			mockDo: func() *mock2.ClientMock {
				return new(mock2.ClientMock)
			},
			mockRespID:  nil,
			mockRespGet: nil,
			mockErr:     nil,
			expectError: true,
			expectUser:  nil,
		},
		{
			name:  "failed to get user after creation",
			input: createUserInput,
			mockConvert: func() *automock.UserConverter {
				userConverter := &automock.UserConverter{}
				userConverter.EXPECT().ConvertCreateUserInput(createUserInput).Return(httpInput, nil)
				userConverter.EXPECT().ConvertUserToGraphQL(inputUser).Return(&expectedUser, nil)
				return userConverter
			},
			mockListConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "POST", "users/create", data).Return([]byte(`"1"`), nil)
				mockClient.On("Do", mock.Anything, "GET", "users/1", mock.Anything).Return(data, errors.New("failed to fetch user"))
				return mockClient
			},
			mockRespID:  []byte(`"1"`),
			mockRespGet: nil,
			mockErr:     errors.New("failed to fetch user"),
			expectError: true,
			expectUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := tt.mockDo()
			userConverter := tt.mockConvert()
			listConverter := tt.mockListConverter()
			r := user.NewResolver(mockClient, userConverter, listConverter)

			result, err := r.CreateUser(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectUser, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestUpdateUser_UserResolver(t *testing.T) {
	expectedUser := graphql.User{
		ID:    "1",
		Email: "victor",
	}

	updatedUser := models.User{
		ID:    "1",
		Email: "updated_victor",
	}
	email := "updated_victor"

	updateUserInput := graphql.UpdateUserInput{
		Email: &email,
	}

	httpInput := models.User{
		Email: "updated_victor",
	}

	data, _ := json.Marshal(httpInput)

	tests := []struct {
		name              string
		id                string
		input             graphql.UpdateUserInput
		mockConvert       func() *automock.UserConverter
		mockListConverter func() *automock.ListConverter
		mockDo            func() *mock2.ClientMock
		mockRespGet       []byte
		mockErr           error
		expectError       bool
		expectUser        *graphql.User
	}{
		{
			name:  "successful user update",
			id:    "1",
			input: updateUserInput,
			mockConvert: func() *automock.UserConverter {
				userConverter := &automock.UserConverter{}
				userConverter.EXPECT().ConvertUpdateUserInput(updateUserInput).Return(httpInput, nil)
				userConverter.EXPECT().ConvertUserToGraphQL(updatedUser).Return(&expectedUser, nil)
				return userConverter
			},
			mockListConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "PUT", "users/1", data).Return([]byte(`"1"`), nil)
				mockClient.On("Do", mock.Anything, "GET", "users/1", mock.Anything).Return([]byte(`{"ID": "1", "Email": "updated_victor"}`), nil)
				return mockClient
			},
			mockRespGet: []byte(`{"ID": "1", "Email": "updated_victor"}`),
			mockErr:     nil,
			expectError: false,
			expectUser:  &expectedUser,
		},
		{
			name:  "failed to convert update user input",
			id:    "1",
			input: updateUserInput,
			mockConvert: func() *automock.UserConverter {
				userConverter := &automock.UserConverter{}
				userConverter.EXPECT().ConvertUpdateUserInput(updateUserInput).Return(models.User{}, errors.New("conversion error"))
				return userConverter
			},
			mockListConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
			mockDo: func() *mock2.ClientMock {
				return new(mock2.ClientMock)
			},
			mockRespGet: nil,
			mockErr:     nil,
			expectError: true,
			expectUser:  nil,
		},
		{
			name:  "failed to update user",
			id:    "1",
			input: updateUserInput,
			mockConvert: func() *automock.UserConverter {
				userConverter := &automock.UserConverter{}
				userConverter.EXPECT().ConvertUpdateUserInput(updateUserInput).Return(httpInput, nil)
				return userConverter
			},
			mockListConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "PUT", "users/1", data).Return([]byte(`"1"`), errors.New("failed to update user"))
				return mockClient
			},
			mockRespGet: nil,
			mockErr:     errors.New("failed to update user"),
			expectError: true,
			expectUser:  nil,
		},
		{
			name:  "failed to fetch updated user after update",
			id:    "1",
			input: updateUserInput,
			mockConvert: func() *automock.UserConverter {
				userConverter := &automock.UserConverter{}
				userConverter.EXPECT().ConvertUpdateUserInput(updateUserInput).Return(httpInput, nil)
				return userConverter
			},
			mockListConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "PUT", "users/1", data).Return([]byte(`"1"`), nil)
				mockClient.On("Do", mock.Anything, "GET", "users/1", mock.Anything).Return([]byte(`"1"`), errors.New("failed to fetch user"))
				return mockClient
			},
			mockRespGet: nil,
			mockErr:     errors.New("failed to fetch user"),
			expectError: true,
			expectUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := tt.mockDo()
			userConverter := tt.mockConvert()
			listConverter := tt.mockListConverter()
			r := user.NewResolver(mockClient, userConverter, listConverter)

			result, err := r.UpdateUser(context.Background(), tt.id, tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectUser, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestDeleteUser_UserResolver(t *testing.T) {
	expectedUser := graphql.User{
		ID:    "1",
		Email: "victor",
	}

	existingUser := models.User{
		ID:    "1",
		Email: "victor",
	}

	tests := []struct {
		name        string
		id          string
		mockConvert func() *automock.UserConverter
		mockDo      func() *mock2.ClientMock
		mockRespGet []byte
		mockErr     error
		expectError bool
		expectUser  *graphql.User
	}{
		{
			name: "successful user deletion",
			id:   "1",
			mockConvert: func() *automock.UserConverter {
				userConverter := &automock.UserConverter{}
				userConverter.EXPECT().ConvertUserToGraphQL(existingUser).Return(&expectedUser, nil)
				return userConverter
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "GET", "users/1", mock.Anything).Return([]byte(`{"ID": "1", "Email": "victor"}`), nil)
				mockClient.On("Do", mock.Anything, "DELETE", "users/1", mock.Anything).Return([]byte(`"1"`), nil)
				return mockClient
			},
			mockRespGet: []byte(`{"ID": "1", "Email": "victor"}`),
			mockErr:     nil,
			expectError: false,
			expectUser:  &expectedUser,
		},
		{
			name: "failed to fetch user before deletion",
			id:   "1",
			mockConvert: func() *automock.UserConverter {
				return &automock.UserConverter{}
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "GET", "users/1", mock.Anything).Return([]byte(`"1"`), errors.New("failed to fetch user"))
				return mockClient
			},
			mockRespGet: nil,
			mockErr:     errors.New("failed to fetch user"),
			expectError: true,
			expectUser:  nil,
		},
		{
			name: "failed to delete user",
			id:   "1",
			mockConvert: func() *automock.UserConverter {
				userConverter := &automock.UserConverter{}
				userConverter.EXPECT().ConvertUserToGraphQL(existingUser).Return(&expectedUser, nil)
				return userConverter
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "GET", "users/1", mock.Anything).Return([]byte(`{"ID": "1", "Email": "victor"}`), nil)
				mockClient.On("Do", mock.Anything, "DELETE", "users/1", mock.Anything).Return([]byte(`"1"`), errors.New("failed to delete user"))
				return mockClient
			},
			mockRespGet: []byte(`{"ID": "1", "Email": "victor"}`),
			mockErr:     errors.New("failed to delete user"),
			expectError: true,
			expectUser:  nil,
		},
		{
			name: "failed to unmarshal user",
			id:   "1",
			mockConvert: func() *automock.UserConverter {
				return &automock.UserConverter{}
			},
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "GET", "users/1", mock.Anything).Return([]byte(`invalid_json`), nil)
				return mockClient
			},
			mockRespGet: []byte(`invalid_json`),
			mockErr:     errors.New("failed to unmarshal user"),
			expectError: true,
			expectUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := tt.mockDo()
			userConverter := tt.mockConvert()
			r := user.NewResolver(mockClient, userConverter, nil)

			result, err := r.DeleteUser(context.Background(), tt.id)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectUser, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestGetList_UserResolver(t *testing.T) {
	listID := "101"
	expectedList := &graphql.List{
		ID:   listID,
		Name: "Test List",
	}

	modelList := models.List{
		ID:   listID,
		Name: "Test List",
	}

	tests := []struct {
		name        string
		listID      string
		mockDo      func() *mock2.ClientMock
		mockConvert func() *automock.ListConverter
		mockResp    []byte
		mockErr     error
		expectError bool
		expectList  *graphql.List
	}{
		{
			name:   "successful list retrieval",
			listID: listID,
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "GET", fmt.Sprintf("lists/%s", listID), mock.Anything).Return([]byte(`{"ID": "101", "Name": "Test List"}`), nil)
				return mockClient
			},
			mockConvert: func() *automock.ListConverter {
				listConverter := &automock.ListConverter{}
				listConverter.EXPECT().ConvertListToGraphQL(modelList).Return(expectedList, nil)
				return listConverter
			},
			mockResp:    []byte(`{"ID": "101", "Name": "Test List"}`),
			mockErr:     nil,
			expectError: false,
			expectList:  expectedList,
		},
		{
			name:   "failed to fetch list",
			listID: listID,
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "GET", fmt.Sprintf("lists/%s", listID), mock.Anything).Return([]byte(`"1"`), errors.New("failed to fetch list"))
				return mockClient
			},
			mockConvert: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
			mockResp:    nil,
			mockErr:     errors.New("failed to fetch list"),
			expectError: true,
			expectList:  nil,
		},
		{
			name:   "failed to unmarshal response",
			listID: listID,
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "GET", fmt.Sprintf("lists/%s", listID), mock.Anything).Return([]byte(`invalid_json`), nil)
				return mockClient
			},
			mockConvert: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
			mockResp:    []byte(`invalid_json`),
			mockErr:     errors.New("failed to unmarshal list"),
			expectError: true,
			expectList:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := tt.mockDo()
			mockListConverter := tt.mockConvert()

			r := user.NewResolver(mockClient, nil, mockListConverter)

			result, err := r.GetList(context.Background(), tt.listID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectList, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestGetUser_UserResolver(t *testing.T) {
	userID := "1"
	expectedUser := &graphql.User{
		ID:    userID,
		Email: "test@example.com",
	}

	modelUser := models.User{
		ID:    userID,
		Email: "test@example.com",
	}

	tests := []struct {
		name        string
		userID      string
		mockDo      func() *mock2.ClientMock
		mockConvert func() *automock.UserConverter
		mockResp    []byte
		mockErr     error
		expectError bool
		expectUser  *graphql.User
	}{
		{
			name:   "successful user retrieval",
			userID: userID,
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "GET", fmt.Sprintf("users/%s", userID), mock.Anything).Return([]byte(`{"ID": "1", "Email": "test@example.com"}`), nil)
				return mockClient
			},
			mockConvert: func() *automock.UserConverter {
				userConverter := &automock.UserConverter{}
				userConverter.EXPECT().ConvertUserToGraphQL(modelUser).Return(expectedUser, nil)
				return userConverter
			},
			mockResp:    []byte(`{"ID": "1", "Email": "test@example.com"}`),
			mockErr:     nil,
			expectError: false,
			expectUser:  expectedUser,
		},
		{
			name:   "failed to fetch user",
			userID: userID,
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "GET", fmt.Sprintf("users/%s", userID), mock.Anything).Return([]byte(`"1"`), errors.New("failed to fetch user"))
				return mockClient
			},
			mockConvert: func() *automock.UserConverter {
				return &automock.UserConverter{}
			},
			mockResp:    nil,
			mockErr:     errors.New("failed to fetch user"),
			expectError: true,
			expectUser:  nil,
		},
		{
			name:   "failed to unmarshal response",
			userID: userID,
			mockDo: func() *mock2.ClientMock {
				mockClient := new(mock2.ClientMock)
				mockClient.On("Do", mock.Anything, "GET", fmt.Sprintf("users/%s", userID), mock.Anything).Return([]byte(`invalid_json`), nil)
				return mockClient
			},
			mockConvert: func() *automock.UserConverter {
				return &automock.UserConverter{}
			},
			mockResp:    []byte(`invalid_json`),
			mockErr:     errors.New("failed to unmarshal user"),
			expectError: true,
			expectUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := tt.mockDo()
			mockUserConverter := tt.mockConvert()

			r := user.NewResolver(mockClient, mockUserConverter, nil)

			result, err := r.GetUser(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectUser, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
