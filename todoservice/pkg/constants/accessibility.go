package constants

type Accessibility string

const (
	IsOwner       Accessibility = "owner"
	HasAccessList Accessibility = "has_access_list"
	HasAccessTodo Accessibility = "has_access_todo"
	NoRestriction Accessibility = "no_restriction"
)
