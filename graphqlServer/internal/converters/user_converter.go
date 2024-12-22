package converters

import (
	"fmt"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/generated/graphql"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
)

type ConverterUserGraphQL struct{}

//go:generate mockery --name=UserConverter --output=automock --with-expecter=true --outpkg=automock --case=underscore --disable-version-string
type UserConverter interface {
	ConvertUserToGraphQL(user models.User) (*graphql.User, error)
	ConvertCreateUserInput(input graphql.CreateUserInput) (models.User, error)
	ConvertUpdateUserInput(input graphql.UpdateUserInput) (models.User, error)
}

func NewConverterUserGraphQL() UserConverter {
	return &ConverterUserGraphQL{}
}

func (c *ConverterUserGraphQL) ConvertUserToGraphQL(user models.User) (*graphql.User, error) {
	role, err := ConvertRoleToGraphQL(user.Role)
	if err != nil {
		return nil, fmt.Errorf("convert role to graphql: %w", err)
	}
	return &graphql.User{
		ID:        user.ID,
		Email:     user.Email,
		GithubID:  user.GithubID,
		Role:      role,
		UpdatedAt: user.UpdatedAt.Format(constants.DateFormat),
		CreatedAt: user.CreatedAt.Format(constants.DateFormat),
	}, nil
}

func (c *ConverterUserGraphQL) ConvertCreateUserInput(input graphql.CreateUserInput) (models.User, error) {
	role, err := ConvertRoleFromGraphQL(input.Role)
	if err != nil {
		return models.User{}, fmt.Errorf("convertRoleFromGraphQL: %w", err)
	}
	return models.User{
		GithubID: input.GithubID,
		Email:    input.Email,
		Role:     role,
	}, nil
}

func (c *ConverterUserGraphQL) ConvertUpdateUserInput(input graphql.UpdateUserInput) (models.User, error) {
	role, err := ConvertRoleFromGraphQL(*input.Role)
	if err != nil {
		return models.User{}, fmt.Errorf("convertRoleFromGraphQL: %w", err)
	}
	return models.User{
		GithubID: *input.GithubID,
		Email:    *input.Email,
		Role:     role,
	}, nil
}
