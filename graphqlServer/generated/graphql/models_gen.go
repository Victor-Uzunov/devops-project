// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphql

import (
	"fmt"
	"io"
	"strconv"
)

type CreateListInput struct {
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Visibility  Visibility `json:"visibility"`
	Tags        []string   `json:"tags,omitempty"`
	Shared      []string   `json:"shared,omitempty"`
}

type CreateTodoInput struct {
	ListID      string    `json:"listId"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	DueDate     *string   `json:"dueDate,omitempty"`
	StartDate   *string   `json:"startDate,omitempty"`
	Priority    *Priority `json:"priority,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Completed   *bool     `json:"completed,omitempty"`
	AssignedTo  *string   `json:"assignedTo,omitempty"`
}

type CreateUserInput struct {
	Email    string   `json:"email"`
	GithubID string   `json:"githubId"`
	Role     UserRole `json:"role"`
}

type GrantListAccessInput struct {
	ListID      string      `json:"listId"`
	UserID      string      `json:"userId"`
	AccessLevel AccessLevel `json:"accessLevel"`
	Status      *string     `json:"status,omitempty"`
}

type List struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Description   *string       `json:"description,omitempty"`
	Owner         *User         `json:"owner"`
	Visibility    Visibility    `json:"visibility"`
	Tags          []string      `json:"tags,omitempty"`
	CreatedAt     string        `json:"createdAt"`
	UpdatedAt     string        `json:"updatedAt"`
	Todos         []*Todo       `json:"todos"`
	Collaborators []*ListAccess `json:"collaborators"`
}

type ListAccess struct {
	List        *List       `json:"list"`
	User        *User       `json:"user"`
	AccessLevel AccessLevel `json:"accessLevel"`
	Status      *string     `json:"status,omitempty"`
}

type Mutation struct {
}

type Query struct {
}

type Todo struct {
	ID          string    `json:"id"`
	List        *List     `json:"list"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	DueDate     *string   `json:"dueDate,omitempty"`
	StartDate   *string   `json:"startDate,omitempty"`
	Priority    *Priority `json:"priority,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
	AssignedTo  *User     `json:"assignedTo,omitempty"`
}

type UpdateListInput struct {
	Name        *string     `json:"name,omitempty"`
	Description *string     `json:"description,omitempty"`
	Visibility  *Visibility `json:"visibility,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
}

type UpdateTodoInput struct {
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	Completed   *bool     `json:"completed,omitempty"`
	DueDate     *string   `json:"dueDate,omitempty"`
	StartDate   *string   `json:"startDate,omitempty"`
	Priority    *Priority `json:"priority,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	AssignedTo  *string   `json:"assignedTo,omitempty"`
}

type UpdateUserInput struct {
	GithubID *string   `json:"githubID,omitempty"`
	Email    *string   `json:"email,omitempty"`
	Role     *UserRole `json:"role,omitempty"`
}

type User struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	GithubID  string   `json:"githubID"`
	Role      UserRole `json:"role"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}

type AccessLevel string

const (
	AccessLevelReader AccessLevel = "READER"
	AccessLevelWriter AccessLevel = "WRITER"
	AccessLevelAdmin  AccessLevel = "ADMIN"
)

var AllAccessLevel = []AccessLevel{
	AccessLevelReader,
	AccessLevelWriter,
	AccessLevelAdmin,
}

func (e AccessLevel) IsValid() bool {
	switch e {
	case AccessLevelReader, AccessLevelWriter, AccessLevelAdmin:
		return true
	}
	return false
}

func (e AccessLevel) String() string {
	return string(e)
}

func (e *AccessLevel) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AccessLevel(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AccessLevel", str)
	}
	return nil
}

func (e AccessLevel) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Priority string

const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
)

var AllPriority = []Priority{
	PriorityLow,
	PriorityMedium,
	PriorityHigh,
}

func (e Priority) IsValid() bool {
	switch e {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	}
	return false
}

func (e Priority) String() string {
	return string(e)
}

func (e *Priority) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Priority(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Priority", str)
	}
	return nil
}

func (e Priority) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type UserRole string

const (
	UserRoleReader UserRole = "READER"
	UserRoleWriter UserRole = "WRITER"
	UserRoleAdmin  UserRole = "ADMIN"
)

var AllUserRole = []UserRole{
	UserRoleReader,
	UserRoleWriter,
	UserRoleAdmin,
}

func (e UserRole) IsValid() bool {
	switch e {
	case UserRoleReader, UserRoleWriter, UserRoleAdmin:
		return true
	}
	return false
}

func (e UserRole) String() string {
	return string(e)
}

func (e *UserRole) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UserRole(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UserRole", str)
	}
	return nil
}

func (e UserRole) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Visibility string

const (
	VisibilityPrivate Visibility = "PRIVATE"
	VisibilityShared  Visibility = "SHARED"
	VisibilityPublic  Visibility = "PUBLIC"
)

var AllVisibility = []Visibility{
	VisibilityPrivate,
	VisibilityShared,
	VisibilityPublic,
}

func (e Visibility) IsValid() bool {
	switch e {
	case VisibilityPrivate, VisibilityShared, VisibilityPublic:
		return true
	}
	return false
}

func (e Visibility) String() string {
	return string(e)
}

func (e *Visibility) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Visibility(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Visibility", str)
	}
	return nil
}

func (e Visibility) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
