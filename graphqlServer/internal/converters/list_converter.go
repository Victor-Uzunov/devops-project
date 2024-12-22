package converters

import (
	"encoding/json"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/generated/graphql"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
)

//go:generate mockery --name=ListConverter --output=automock --with-expecter=true --outpkg=automock --case=underscore --disable-version-string
type ListConverter interface {
	ConvertListToGraphQL(list models.List) (*graphql.List, error)
	ConvertCreateListInput(input graphql.CreateListInput, userID string) (models.List, error)
	ConvertUpdateListInput(input graphql.UpdateListInput) (models.List, error)
	ConvertAccessLevelToGraphQL(role constants.Role) (graphql.AccessLevel, error)
	ConvertAccessLevelFromGraphQL(role graphql.AccessLevel) (constants.Role, error)
	ConvertGrantListAccessInputToModel(input graphql.GrantListAccessInput) (models.Access, error)
	ConvertMultipleListsToGraphQL(lists []*models.List) ([]*graphql.List, error)
}

type ConverterListGraphQL struct{}

func NewConverterListGraphQL() ListConverter {
	return &ConverterListGraphQL{}
}

func (c *ConverterListGraphQL) ConvertListToGraphQL(list models.List) (*graphql.List, error) {
	visibility, err := convertVisibilityToGraphQL(list.Visibility)
	if err != nil {
		return &graphql.List{}, fmt.Errorf("converting visibility to graphQL: %w", err)
	}
	var tags []string

	if list.Tags != nil {
		err := json.Unmarshal(list.Tags, &tags)
		if err != nil {
			fmt.Println("Error unmarshalling tags:", err)
			tags = make([]string, 0)
		}
	}
	return &graphql.List{
		ID:            list.ID,
		Name:          list.Name,
		Description:   &list.Description,
		Owner:         nil,
		Visibility:    visibility,
		Tags:          tags,
		CreatedAt:     list.CreatedAt.Format(constants.DateFormat),
		UpdatedAt:     list.UpdatedAt.Format(constants.DateFormat),
		Todos:         make([]*graphql.Todo, 0),
		Collaborators: make([]*graphql.ListAccess, 0),
	}, nil
}

func convertVisibilityToGraphQL(visibility constants.Visibility) (graphql.Visibility, error) {
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

func (c *ConverterListGraphQL) ConvertCreateListInput(input graphql.CreateListInput, userID string) (models.List, error) {
	jsonBytes, err := json.Marshal(input.Tags)
	if err != nil {
		return models.List{}, fmt.Errorf("convertCreateListInput: %w", err)
	}
	jsonRawMessage := json.RawMessage(jsonBytes)

	visibility, err := ConvertVisibilityFromGraphQL(input.Visibility)
	if err != nil {
		return models.List{}, fmt.Errorf("convertVisibilityFromGraphQL: %w", err)
	}

	shared := make([]string, 0)
	for _, collaborator := range input.Shared {
		shared = append(shared, collaborator)
	}
	return models.List{
		Name:        input.Name,
		Description: *input.Description,
		OwnerID:     userID,
		Visibility:  visibility,
		Tags:        jsonRawMessage,
		SharedWith:  shared,
	}, nil
}

func (c *ConverterListGraphQL) ConvertUpdateListInput(input graphql.UpdateListInput) (models.List, error) {
	jsonBytes, err := json.Marshal(input.Tags)
	if err != nil {
		return models.List{}, fmt.Errorf("convertCreateListInput: %w", err)
	}
	jsonRawMessage := json.RawMessage(jsonBytes)

	visibility, err := ConvertVisibilityFromGraphQL(*input.Visibility)
	if err != nil {
		return models.List{}, fmt.Errorf("convertVisibilityFromGraphQL: %w", err)
	}
	return models.List{
		Name:        *input.Name,
		Description: *input.Description,
		Visibility:  visibility,
		Tags:        jsonRawMessage,
	}, nil
}

func (c *ConverterListGraphQL) ConvertAccessLevelToGraphQL(role constants.Role) (graphql.AccessLevel, error) {
	switch role {
	case constants.Admin:
		return graphql.AccessLevelAdmin, nil
	case constants.Reader:
		return graphql.AccessLevelReader, nil
	case constants.Writer:
		return graphql.AccessLevelWriter, nil
	default:
		return graphql.AccessLevelReader, fmt.Errorf("invalid role level: %v", role)
	}
}

func (c *ConverterListGraphQL) ConvertAccessLevelFromGraphQL(role graphql.AccessLevel) (constants.Role, error) {
	switch role {
	case graphql.AccessLevelAdmin:
		return constants.Admin, nil
	case graphql.AccessLevelReader:
		return constants.Reader, nil
	case graphql.AccessLevelWriter:
		return constants.Writer, nil
	default:
		return constants.Reader, fmt.Errorf("invalid role level: %v", role)
	}
}

func (c *ConverterListGraphQL) ConvertGrantListAccessInputToModel(input graphql.GrantListAccessInput) (models.Access, error) {
	role, err := c.ConvertAccessLevelFromGraphQL(input.AccessLevel)
	if err != nil {
		return models.Access{}, fmt.Errorf("invalid access level: %v", input.AccessLevel)
	}
	return models.Access{
		ListID: input.ListID,
		UserID: input.UserID,
		Role:   role,
	}, nil
}

func (c *ConverterListGraphQL) ConvertMultipleListsToGraphQL(lists []*models.List) ([]*graphql.List, error) {
	result := make([]*graphql.List, 0)
	for _, list := range lists {
		l, err := c.ConvertListToGraphQL(*list)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling response: %w", err)
		}
		result = append(result, l)
	}
	return result, nil
}
