// Code generated by mockery. DO NOT EDIT.

package automock

import (
	constants "github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"

	graphql "github.com/Victor-Uzunov/devops-project/graphqlServer/generated/graphql"

	mock "github.com/stretchr/testify/mock"

	models "github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
)

// ListConverter is an autogenerated mock type for the ListConverter type
type ListConverter struct {
	mock.Mock
}

type ListConverter_Expecter struct {
	mock *mock.Mock
}

func (_m *ListConverter) EXPECT() *ListConverter_Expecter {
	return &ListConverter_Expecter{mock: &_m.Mock}
}

// ConvertAccessLevelFromGraphQL provides a mock function with given fields: role
func (_m *ListConverter) ConvertAccessLevelFromGraphQL(role graphql.AccessLevel) (constants.Role, error) {
	ret := _m.Called(role)

	if len(ret) == 0 {
		panic("no return value specified for ConvertAccessLevelFromGraphQL")
	}

	var r0 constants.Role
	var r1 error
	if rf, ok := ret.Get(0).(func(graphql.AccessLevel) (constants.Role, error)); ok {
		return rf(role)
	}
	if rf, ok := ret.Get(0).(func(graphql.AccessLevel) constants.Role); ok {
		r0 = rf(role)
	} else {
		r0 = ret.Get(0).(constants.Role)
	}

	if rf, ok := ret.Get(1).(func(graphql.AccessLevel) error); ok {
		r1 = rf(role)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListConverter_ConvertAccessLevelFromGraphQL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConvertAccessLevelFromGraphQL'
type ListConverter_ConvertAccessLevelFromGraphQL_Call struct {
	*mock.Call
}

// ConvertAccessLevelFromGraphQL is a helper method to define mock.On call
//   - role graphql.AccessLevel
func (_e *ListConverter_Expecter) ConvertAccessLevelFromGraphQL(role interface{}) *ListConverter_ConvertAccessLevelFromGraphQL_Call {
	return &ListConverter_ConvertAccessLevelFromGraphQL_Call{Call: _e.mock.On("ConvertAccessLevelFromGraphQL", role)}
}

func (_c *ListConverter_ConvertAccessLevelFromGraphQL_Call) Run(run func(role graphql.AccessLevel)) *ListConverter_ConvertAccessLevelFromGraphQL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(graphql.AccessLevel))
	})
	return _c
}

func (_c *ListConverter_ConvertAccessLevelFromGraphQL_Call) Return(_a0 constants.Role, _a1 error) *ListConverter_ConvertAccessLevelFromGraphQL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ListConverter_ConvertAccessLevelFromGraphQL_Call) RunAndReturn(run func(graphql.AccessLevel) (constants.Role, error)) *ListConverter_ConvertAccessLevelFromGraphQL_Call {
	_c.Call.Return(run)
	return _c
}

// ConvertAccessLevelToGraphQL provides a mock function with given fields: role
func (_m *ListConverter) ConvertAccessLevelToGraphQL(role constants.Role) (graphql.AccessLevel, error) {
	ret := _m.Called(role)

	if len(ret) == 0 {
		panic("no return value specified for ConvertAccessLevelToGraphQL")
	}

	var r0 graphql.AccessLevel
	var r1 error
	if rf, ok := ret.Get(0).(func(constants.Role) (graphql.AccessLevel, error)); ok {
		return rf(role)
	}
	if rf, ok := ret.Get(0).(func(constants.Role) graphql.AccessLevel); ok {
		r0 = rf(role)
	} else {
		r0 = ret.Get(0).(graphql.AccessLevel)
	}

	if rf, ok := ret.Get(1).(func(constants.Role) error); ok {
		r1 = rf(role)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListConverter_ConvertAccessLevelToGraphQL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConvertAccessLevelToGraphQL'
type ListConverter_ConvertAccessLevelToGraphQL_Call struct {
	*mock.Call
}

// ConvertAccessLevelToGraphQL is a helper method to define mock.On call
//   - role constants.Role
func (_e *ListConverter_Expecter) ConvertAccessLevelToGraphQL(role interface{}) *ListConverter_ConvertAccessLevelToGraphQL_Call {
	return &ListConverter_ConvertAccessLevelToGraphQL_Call{Call: _e.mock.On("ConvertAccessLevelToGraphQL", role)}
}

func (_c *ListConverter_ConvertAccessLevelToGraphQL_Call) Run(run func(role constants.Role)) *ListConverter_ConvertAccessLevelToGraphQL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(constants.Role))
	})
	return _c
}

func (_c *ListConverter_ConvertAccessLevelToGraphQL_Call) Return(_a0 graphql.AccessLevel, _a1 error) *ListConverter_ConvertAccessLevelToGraphQL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ListConverter_ConvertAccessLevelToGraphQL_Call) RunAndReturn(run func(constants.Role) (graphql.AccessLevel, error)) *ListConverter_ConvertAccessLevelToGraphQL_Call {
	_c.Call.Return(run)
	return _c
}

// ConvertCreateListInput provides a mock function with given fields: input, userID
func (_m *ListConverter) ConvertCreateListInput(input graphql.CreateListInput, userID string) (models.List, error) {
	ret := _m.Called(input, userID)

	if len(ret) == 0 {
		panic("no return value specified for ConvertCreateListInput")
	}

	var r0 models.List
	var r1 error
	if rf, ok := ret.Get(0).(func(graphql.CreateListInput, string) (models.List, error)); ok {
		return rf(input, userID)
	}
	if rf, ok := ret.Get(0).(func(graphql.CreateListInput, string) models.List); ok {
		r0 = rf(input, userID)
	} else {
		r0 = ret.Get(0).(models.List)
	}

	if rf, ok := ret.Get(1).(func(graphql.CreateListInput, string) error); ok {
		r1 = rf(input, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListConverter_ConvertCreateListInput_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConvertCreateListInput'
type ListConverter_ConvertCreateListInput_Call struct {
	*mock.Call
}

// ConvertCreateListInput is a helper method to define mock.On call
//   - input graphql.CreateListInput
//   - userID string
func (_e *ListConverter_Expecter) ConvertCreateListInput(input interface{}, userID interface{}) *ListConverter_ConvertCreateListInput_Call {
	return &ListConverter_ConvertCreateListInput_Call{Call: _e.mock.On("ConvertCreateListInput", input, userID)}
}

func (_c *ListConverter_ConvertCreateListInput_Call) Run(run func(input graphql.CreateListInput, userID string)) *ListConverter_ConvertCreateListInput_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(graphql.CreateListInput), args[1].(string))
	})
	return _c
}

func (_c *ListConverter_ConvertCreateListInput_Call) Return(_a0 models.List, _a1 error) *ListConverter_ConvertCreateListInput_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ListConverter_ConvertCreateListInput_Call) RunAndReturn(run func(graphql.CreateListInput, string) (models.List, error)) *ListConverter_ConvertCreateListInput_Call {
	_c.Call.Return(run)
	return _c
}

// ConvertGrantListAccessInputToModel provides a mock function with given fields: input
func (_m *ListConverter) ConvertGrantListAccessInputToModel(input graphql.GrantListAccessInput) (models.Access, error) {
	ret := _m.Called(input)

	if len(ret) == 0 {
		panic("no return value specified for ConvertGrantListAccessInputToModel")
	}

	var r0 models.Access
	var r1 error
	if rf, ok := ret.Get(0).(func(graphql.GrantListAccessInput) (models.Access, error)); ok {
		return rf(input)
	}
	if rf, ok := ret.Get(0).(func(graphql.GrantListAccessInput) models.Access); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Get(0).(models.Access)
	}

	if rf, ok := ret.Get(1).(func(graphql.GrantListAccessInput) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListConverter_ConvertGrantListAccessInputToModel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConvertGrantListAccessInputToModel'
type ListConverter_ConvertGrantListAccessInputToModel_Call struct {
	*mock.Call
}

// ConvertGrantListAccessInputToModel is a helper method to define mock.On call
//   - input graphql.GrantListAccessInput
func (_e *ListConverter_Expecter) ConvertGrantListAccessInputToModel(input interface{}) *ListConverter_ConvertGrantListAccessInputToModel_Call {
	return &ListConverter_ConvertGrantListAccessInputToModel_Call{Call: _e.mock.On("ConvertGrantListAccessInputToModel", input)}
}

func (_c *ListConverter_ConvertGrantListAccessInputToModel_Call) Run(run func(input graphql.GrantListAccessInput)) *ListConverter_ConvertGrantListAccessInputToModel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(graphql.GrantListAccessInput))
	})
	return _c
}

func (_c *ListConverter_ConvertGrantListAccessInputToModel_Call) Return(_a0 models.Access, _a1 error) *ListConverter_ConvertGrantListAccessInputToModel_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ListConverter_ConvertGrantListAccessInputToModel_Call) RunAndReturn(run func(graphql.GrantListAccessInput) (models.Access, error)) *ListConverter_ConvertGrantListAccessInputToModel_Call {
	_c.Call.Return(run)
	return _c
}

// ConvertListToGraphQL provides a mock function with given fields: list
func (_m *ListConverter) ConvertListToGraphQL(list models.List) (*graphql.List, error) {
	ret := _m.Called(list)

	if len(ret) == 0 {
		panic("no return value specified for ConvertListToGraphQL")
	}

	var r0 *graphql.List
	var r1 error
	if rf, ok := ret.Get(0).(func(models.List) (*graphql.List, error)); ok {
		return rf(list)
	}
	if rf, ok := ret.Get(0).(func(models.List) *graphql.List); ok {
		r0 = rf(list)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graphql.List)
		}
	}

	if rf, ok := ret.Get(1).(func(models.List) error); ok {
		r1 = rf(list)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListConverter_ConvertListToGraphQL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConvertListToGraphQL'
type ListConverter_ConvertListToGraphQL_Call struct {
	*mock.Call
}

// ConvertListToGraphQL is a helper method to define mock.On call
//   - list models.List
func (_e *ListConverter_Expecter) ConvertListToGraphQL(list interface{}) *ListConverter_ConvertListToGraphQL_Call {
	return &ListConverter_ConvertListToGraphQL_Call{Call: _e.mock.On("ConvertListToGraphQL", list)}
}

func (_c *ListConverter_ConvertListToGraphQL_Call) Run(run func(list models.List)) *ListConverter_ConvertListToGraphQL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(models.List))
	})
	return _c
}

func (_c *ListConverter_ConvertListToGraphQL_Call) Return(_a0 *graphql.List, _a1 error) *ListConverter_ConvertListToGraphQL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ListConverter_ConvertListToGraphQL_Call) RunAndReturn(run func(models.List) (*graphql.List, error)) *ListConverter_ConvertListToGraphQL_Call {
	_c.Call.Return(run)
	return _c
}

// ConvertMultipleListsToGraphQL provides a mock function with given fields: lists
func (_m *ListConverter) ConvertMultipleListsToGraphQL(lists []*models.List) ([]*graphql.List, error) {
	ret := _m.Called(lists)

	if len(ret) == 0 {
		panic("no return value specified for ConvertMultipleListsToGraphQL")
	}

	var r0 []*graphql.List
	var r1 error
	if rf, ok := ret.Get(0).(func([]*models.List) ([]*graphql.List, error)); ok {
		return rf(lists)
	}
	if rf, ok := ret.Get(0).(func([]*models.List) []*graphql.List); ok {
		r0 = rf(lists)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*graphql.List)
		}
	}

	if rf, ok := ret.Get(1).(func([]*models.List) error); ok {
		r1 = rf(lists)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListConverter_ConvertMultipleListsToGraphQL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConvertMultipleListsToGraphQL'
type ListConverter_ConvertMultipleListsToGraphQL_Call struct {
	*mock.Call
}

// ConvertMultipleListsToGraphQL is a helper method to define mock.On call
//   - lists []*models.List
func (_e *ListConverter_Expecter) ConvertMultipleListsToGraphQL(lists interface{}) *ListConverter_ConvertMultipleListsToGraphQL_Call {
	return &ListConverter_ConvertMultipleListsToGraphQL_Call{Call: _e.mock.On("ConvertMultipleListsToGraphQL", lists)}
}

func (_c *ListConverter_ConvertMultipleListsToGraphQL_Call) Run(run func(lists []*models.List)) *ListConverter_ConvertMultipleListsToGraphQL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]*models.List))
	})
	return _c
}

func (_c *ListConverter_ConvertMultipleListsToGraphQL_Call) Return(_a0 []*graphql.List, _a1 error) *ListConverter_ConvertMultipleListsToGraphQL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ListConverter_ConvertMultipleListsToGraphQL_Call) RunAndReturn(run func([]*models.List) ([]*graphql.List, error)) *ListConverter_ConvertMultipleListsToGraphQL_Call {
	_c.Call.Return(run)
	return _c
}

// ConvertUpdateListInput provides a mock function with given fields: input
func (_m *ListConverter) ConvertUpdateListInput(input graphql.UpdateListInput) (models.List, error) {
	ret := _m.Called(input)

	if len(ret) == 0 {
		panic("no return value specified for ConvertUpdateListInput")
	}

	var r0 models.List
	var r1 error
	if rf, ok := ret.Get(0).(func(graphql.UpdateListInput) (models.List, error)); ok {
		return rf(input)
	}
	if rf, ok := ret.Get(0).(func(graphql.UpdateListInput) models.List); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Get(0).(models.List)
	}

	if rf, ok := ret.Get(1).(func(graphql.UpdateListInput) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListConverter_ConvertUpdateListInput_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConvertUpdateListInput'
type ListConverter_ConvertUpdateListInput_Call struct {
	*mock.Call
}

// ConvertUpdateListInput is a helper method to define mock.On call
//   - input graphql.UpdateListInput
func (_e *ListConverter_Expecter) ConvertUpdateListInput(input interface{}) *ListConverter_ConvertUpdateListInput_Call {
	return &ListConverter_ConvertUpdateListInput_Call{Call: _e.mock.On("ConvertUpdateListInput", input)}
}

func (_c *ListConverter_ConvertUpdateListInput_Call) Run(run func(input graphql.UpdateListInput)) *ListConverter_ConvertUpdateListInput_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(graphql.UpdateListInput))
	})
	return _c
}

func (_c *ListConverter_ConvertUpdateListInput_Call) Return(_a0 models.List, _a1 error) *ListConverter_ConvertUpdateListInput_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ListConverter_ConvertUpdateListInput_Call) RunAndReturn(run func(graphql.UpdateListInput) (models.List, error)) *ListConverter_ConvertUpdateListInput_Call {
	_c.Call.Return(run)
	return _c
}

// NewListConverter creates a new instance of ListConverter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewListConverter(t interface {
	mock.TestingT
	Cleanup(func())
}) *ListConverter {
	mock := &ListConverter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}