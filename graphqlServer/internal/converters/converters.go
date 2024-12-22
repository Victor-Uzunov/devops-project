package converters

import (
	"fmt"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/generated/graphql"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
)

func ConvertPriorityToGraphQL(priority constants.PriorityLevel) (graphql.Priority, error) {
	switch priority {
	case constants.PriorityLow:
		return graphql.PriorityLow, nil
	case constants.PriorityMedium:
		return graphql.PriorityMedium, nil
	case constants.PriorityHigh:
		return graphql.PriorityHigh, nil
	default:
		return graphql.Priority(constants.PriorityLow), fmt.Errorf("invalid priority level: %v", priority)
	}
}

func ConvertPriorityFromGraphQL(priority graphql.Priority) (constants.PriorityLevel, error) {
	switch priority {
	case graphql.PriorityLow:
		return constants.PriorityLow, nil
	case graphql.PriorityMedium:
		return constants.PriorityMedium, nil
	case graphql.PriorityHigh:
		return constants.PriorityHigh, nil
	default:
		return constants.PriorityLow, fmt.Errorf("invalid priority level: %v", priority)
	}
}

func ConvertVisibilityToGraphQL(visibility constants.Visibility) (graphql.Visibility, error) {
	switch visibility {
	case constants.VisibilityShared:
		return graphql.VisibilityShared, nil
	case constants.VisibilityPrivate:
		return graphql.VisibilityPrivate, nil
	case constants.VisibilityPublic:
		return graphql.VisibilityPublic, nil
	default:
		return graphql.Visibility(constants.VisibilityPublic), fmt.Errorf("invalid priority level: %v", visibility)
	}
}

func ConvertVisibilityFromGraphQL(visibility graphql.Visibility) (constants.Visibility, error) {
	switch visibility {
	case graphql.VisibilityShared:
		return constants.VisibilityShared, nil
	case graphql.VisibilityPrivate:
		return constants.VisibilityPrivate, nil
	case graphql.VisibilityPublic:
		return constants.VisibilityPublic, nil
	default:
		return constants.VisibilityPublic, fmt.Errorf("invalid visibility level: %v", visibility)
	}
}

func ConvertRoleFromGraphQL(role graphql.UserRole) (constants.Role, error) {
	switch role {
	case graphql.UserRoleAdmin:
		return constants.Admin, nil
	case graphql.UserRoleReader:
		return constants.Reader, nil
	case graphql.UserRoleWriter:
		return constants.Writer, nil
	default:
		return constants.Reader, fmt.Errorf("invalid role level: %v", role)
	}
}

func ConvertRoleToGraphQL(role constants.Role) (graphql.UserRole, error) {
	switch role {
	case constants.Admin:
		return graphql.UserRoleAdmin, nil
	case constants.Reader:
		return graphql.UserRoleReader, nil
	case constants.Writer:
		return graphql.UserRoleWriter, nil
	default:
		return graphql.UserRoleReader, fmt.Errorf("invalid role: %v", role)
	}
}
