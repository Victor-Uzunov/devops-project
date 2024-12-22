package list_test

import (
	"context"
	"errors"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/generated/graphql"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/converters/automock"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/resolvers/list"
	mock2 "github.com/Victor-Uzunov/devops-project/graphqlServer/internal/resolvers/mock"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestLists_ListResolver(t *testing.T) {
	expectedList := graphql.List{
		ID:   "1",
		Name: "Test List",
	}
	inputList := models.List{
		ID: "1",
	}

	tests := []struct {
		name          string
		mockMethod    string
		mockURL       string
		mockResp      []byte
		mockErr       error
		expectError   bool
		expectLists   []*graphql.List
		listConverter func() *automock.ListConverter
	}{
		{
			name:        "successful lists fetch",
			mockMethod:  "GET",
			mockURL:     "lists/all",
			mockResp:    []byte(`[{"ID": "1", "Title": "Test List"}]`),
			mockErr:     nil,
			expectError: false,
			expectLists: []*graphql.List{
				{ID: "1", Name: "Test List"},
			},
			listConverter: func() *automock.ListConverter {
				listConverter := &automock.ListConverter{}
				listConverter.EXPECT().ConvertListToGraphQL(inputList).Return(&expectedList, nil)
				return listConverter
			},
		},
		{
			name:        "failed HTTP request",
			mockMethod:  "GET",
			mockURL:     "lists/all",
			mockResp:    nil,
			mockErr:     errors.New("failed to fetch lists"),
			expectError: true,
			expectLists: nil,
			listConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
		},
		{
			name:        "failed to unmarshal response",
			mockMethod:  "GET",
			mockURL:     "lists/all",
			mockResp:    []byte(`invalid JSON`),
			mockErr:     nil,
			expectError: true,
			expectLists: nil,
			listConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(mock2.ClientMock)

			mockClient.On("Do", mock.Anything, tt.mockMethod, tt.mockURL, mock.Anything).Return(tt.mockResp, tt.mockErr)

			listConverter := tt.listConverter()

			r := list.NewResolver(mockClient, listConverter, nil)

			result, err := r.Lists(context.Background())

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectLists, result)
			}

			mockClient.AssertCalled(t, "Do", mock.Anything, tt.mockMethod, tt.mockURL, mock.Anything)
		})
	}
}

func TestList_ListResolver(t *testing.T) {
	expectedList := graphql.List{
		ID:   "1",
		Name: "Test List",
	}
	inputList := models.List{
		ID: "1",
	}

	tests := []struct {
		name          string
		id            string
		mockMethod    string
		mockURL       string
		mockResp      []byte
		mockErr       error
		expectError   bool
		expectList    *graphql.List
		listConverter func() *automock.ListConverter
	}{
		{
			name:        "successful list fetch by id",
			id:          "1",
			mockMethod:  "GET",
			mockURL:     "lists/1",
			mockResp:    []byte(`{"ID": "1", "Title": "Test List"}`),
			mockErr:     nil,
			expectError: false,
			expectList:  &expectedList,
			listConverter: func() *automock.ListConverter {
				listConverter := &automock.ListConverter{}
				listConverter.EXPECT().ConvertListToGraphQL(inputList).Return(&expectedList, nil)
				return listConverter
			},
		},
		{
			name:        "failed HTTP request",
			id:          "1",
			mockMethod:  "GET",
			mockURL:     "lists/1",
			mockResp:    nil,
			mockErr:     errors.New("failed to fetch list by id"),
			expectError: true,
			expectList:  nil,
			listConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
		},
		{
			name:        "failed to unmarshal response",
			id:          "1",
			mockMethod:  "GET",
			mockURL:     "lists/1",
			mockResp:    []byte(`invalid JSON`),
			mockErr:     nil,
			expectError: true,
			expectList:  nil,
			listConverter: func() *automock.ListConverter {
				return &automock.ListConverter{}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(mock2.ClientMock)

			mockClient.On("Do", mock.Anything, tt.mockMethod, tt.mockURL, mock.Anything).Return(tt.mockResp, tt.mockErr)

			listConverter := tt.listConverter()

			r := list.NewResolver(mockClient, listConverter, nil)

			result, err := r.List(context.Background(), tt.id)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectList, result)
			}

			mockClient.AssertCalled(t, "Do", mock.Anything, tt.mockMethod, tt.mockURL, mock.Anything)
		})
	}
}

func TestCreateList_ListResolver(t *testing.T) {
	expectedList := graphql.List{
		ID:   "1",
		Name: "Test List",
	}

	input := graphql.CreateListInput{
		Name: "Test List",
	}

	httpInput := models.List{
		Name: "Test List",
	}

	tests := []struct {
		name          string
		input         graphql.CreateListInput
		mockMethod    string
		mockURLPost   string
		mockURLGet    string
		mockRespPost  []byte
		mockRespGet   []byte
		mockErrPost   error
		mockErrGet    error
		expectError   bool
		expectList    *graphql.List
		listConverter func() *automock.ListConverter
	}{
		{
			name:         "successful list creation",
			input:        input,
			mockMethod:   "POST",
			mockURLPost:  "lists/create",
			mockURLGet:   "lists/1",
			mockRespPost: []byte(`"1"`),
			mockRespGet:  []byte(`{"ID": "1", "Title": "Test List"}`),
			mockErrPost:  nil,
			mockErrGet:   nil,
			expectError:  false,
			expectList:   &expectedList,
			listConverter: func() *automock.ListConverter {
				listConverter := &automock.ListConverter{}
				listConverter.EXPECT().ConvertCreateListInput(input).Return(httpInput, nil)
				listConverter.EXPECT().ConvertListToGraphQL(models.List{ID: "1"}).Return(&expectedList, nil)
				return listConverter
			},
		},
		{
			name:         "failed POST request",
			input:        input,
			mockMethod:   "POST",
			mockURLPost:  "lists/create",
			mockRespPost: nil,
			mockRespGet:  nil,
			mockErrPost:  errors.New("failed POST request"),
			expectError:  true,
			expectList:   nil,
			listConverter: func() *automock.ListConverter {
				listConverter := &automock.ListConverter{}
				listConverter.EXPECT().ConvertCreateListInput(input).Return(httpInput, nil)
				return listConverter
			},
		},
		{
			name:         "failed GET request after list creation",
			input:        input,
			mockMethod:   "POST",
			mockURLPost:  "lists/create",
			mockURLGet:   "lists/1",
			mockRespPost: []byte(`"1"`),
			mockRespGet:  nil,
			mockErrPost:  nil,
			mockErrGet:   errors.New("failed GET request"),
			expectError:  true,
			expectList:   nil,
			listConverter: func() *automock.ListConverter {
				listConverter := &automock.ListConverter{}
				listConverter.EXPECT().ConvertCreateListInput(input).Return(httpInput, nil)
				return listConverter
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(mock2.ClientMock)

			mockClient.On("Do", mock.Anything, "POST", tt.mockURLPost, mock.Anything).Return(tt.mockRespPost, tt.mockErrPost)

			if tt.mockURLGet != "" {
				mockClient.On("Do", mock.Anything, "GET", tt.mockURLGet, mock.Anything).Return(tt.mockRespGet, tt.mockErrGet)
			}

			listConverter := tt.listConverter()

			r := list.NewResolver(mockClient, listConverter, nil)

			result, err := r.CreateList(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectList, result)
			}

			mockClient.AssertCalled(t, "Do", mock.Anything, "POST", tt.mockURLPost, mock.Anything)
			if tt.mockURLGet != "" {
				mockClient.AssertCalled(t, "Do", mock.Anything, "GET", tt.mockURLGet, mock.Anything)
			}
		})
	}
}

func TestUpdateList_ListResolver(t *testing.T) {
	expectedList := graphql.List{
		ID:   "1",
		Name: "Updated Test List",
	}
	name := "Updated Test List"
	input := graphql.UpdateListInput{
		Name: &name,
	}

	httpInput := models.List{
		Name: "Updated Test List",
	}

	tests := []struct {
		name          string
		id            string
		input         graphql.UpdateListInput
		mockMethodPut string
		mockMethodGet string
		mockURLPut    string
		mockURLGet    string
		mockRespPut   []byte
		mockRespGet   []byte
		mockErrPut    error
		mockErrGet    error
		expectError   bool
		expectList    *graphql.List
		listConverter func() *automock.ListConverter
	}{
		{
			name:          "successful list update",
			id:            "1",
			input:         input,
			mockMethodPut: "PUT",
			mockMethodGet: "GET",
			mockURLPut:    "lists/1",
			mockURLGet:    "lists/1",
			mockRespPut:   nil,
			mockRespGet:   []byte(`{"ID": "1", "Title": "Updated Test List"}`),
			mockErrPut:    nil,
			mockErrGet:    nil,
			expectError:   false,
			expectList:    &expectedList,
			listConverter: func() *automock.ListConverter {
				listConverter := &automock.ListConverter{}
				listConverter.EXPECT().ConvertUpdateListInput(input).Return(httpInput, nil)
				listConverter.EXPECT().ConvertListToGraphQL(models.List{ID: "1"}).Return(&expectedList, nil)
				return listConverter
			},
		},
		{
			name:          "failed GET request after update",
			id:            "1",
			input:         input,
			mockMethodPut: "PUT",
			mockMethodGet: "GET",
			mockURLPut:    "lists/1",
			mockURLGet:    "lists/1",
			mockRespPut:   nil,
			mockRespGet:   nil,
			mockErrPut:    nil,
			mockErrGet:    errors.New("failed GET request"),
			expectError:   true,
			expectList:    nil,
			listConverter: func() *automock.ListConverter {
				listConverter := &automock.ListConverter{}
				listConverter.EXPECT().ConvertUpdateListInput(input).Return(httpInput, nil)
				return listConverter
			},
		},
		{
			name:          "failed to unmarshal GET response",
			id:            "1",
			input:         input,
			mockMethodPut: "PUT",
			mockMethodGet: "GET",
			mockURLPut:    "lists/1",
			mockURLGet:    "lists/1",
			mockRespPut:   nil,
			mockRespGet:   []byte(`invalid JSON`),
			mockErrPut:    nil,
			expectError:   true,
			expectList:    nil,
			listConverter: func() *automock.ListConverter {
				listConverter := &automock.ListConverter{}
				listConverter.EXPECT().ConvertUpdateListInput(input).Return(httpInput, nil)
				return listConverter
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(mock2.ClientMock)

			mockClient.On("Do", mock.Anything, "PUT", tt.mockURLPut, mock.Anything).Return(tt.mockRespPut, tt.mockErrPut)

			mockClient.On("Do", mock.Anything, "GET", tt.mockURLGet, mock.Anything).Return(tt.mockRespGet, tt.mockErrGet)

			listConverter := tt.listConverter()

			r := list.NewResolver(mockClient, listConverter, nil)

			result, err := r.UpdateList(context.Background(), tt.id, tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectList, result)
			}

			mockClient.AssertCalled(t, "Do", mock.Anything, "PUT", tt.mockURLPut, mock.Anything)
			mockClient.AssertCalled(t, "Do", mock.Anything, "GET", tt.mockURLGet, mock.Anything)
		})
	}
}

func TestDeleteList_ListResolver(t *testing.T) {
	expectedList := graphql.List{
		ID:   "1",
		Name: "Test List",
	}

	tests := []struct {
		name          string
		id            string
		mockMethodGet string
		mockMethodDel string
		mockURLGet    string
		mockURLDel    string
		mockRespGet   []byte
		mockRespDel   []byte
		mockErrGet    error
		mockErrDel    error
		expectError   bool
		expectList    *graphql.List
		listConverter func() *automock.ListConverter
	}{
		{
			name:          "successful list deletion",
			id:            "1",
			mockMethodGet: "GET",
			mockMethodDel: "DELETE",
			mockURLGet:    "lists/1",
			mockURLDel:    "lists/1",
			mockRespGet:   []byte(`{"ID": "1", "Title": "Test List"}`),
			mockRespDel:   nil,
			mockErrGet:    nil,
			mockErrDel:    nil,
			expectError:   false,
			expectList:    &expectedList,
			listConverter: func() *automock.ListConverter {
				listConverter := &automock.ListConverter{}
				listConverter.EXPECT().ConvertListToGraphQL(models.List{ID: "1"}).Return(&expectedList, nil)
				return listConverter
			},
		},
		{
			name:          "failed to delete list",
			id:            "1",
			mockMethodGet: "GET",
			mockMethodDel: "DELETE",
			mockURLGet:    "lists/1",
			mockURLDel:    "lists/1",
			mockRespGet:   []byte(`{"ID": "1", "Title": "Test List"}`),
			mockRespDel:   nil,
			mockErrGet:    nil,
			mockErrDel:    errors.New("failed to delete list"),
			expectError:   true,
			expectList:    nil,
			listConverter: func() *automock.ListConverter {
				listConverter := &automock.ListConverter{}
				listConverter.EXPECT().ConvertListToGraphQL(models.List{ID: "1", Name: "Test List"}).Return(&expectedList, nil)
				return listConverter
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(mock2.ClientMock)

			mockClient.On("Do", mock.Anything, "GET", tt.mockURLGet, mock.Anything).Return(tt.mockRespGet, tt.mockErrGet)

			mockClient.On("Do", mock.Anything, "DELETE", tt.mockURLDel, mock.Anything).Return(tt.mockRespDel, tt.mockErrDel)

			listConverter := tt.listConverter()

			r := list.NewResolver(mockClient, listConverter, nil)

			result, err := r.DeleteList(context.Background(), tt.id)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectList, result)
			}

			mockClient.AssertCalled(t, "Do", mock.Anything, "GET", tt.mockURLGet, mock.Anything)
			mockClient.AssertCalled(t, "Do", mock.Anything, "DELETE", tt.mockURLDel, mock.Anything)
		})
	}
}
