package resolvers

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/client"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg"
)

type UserData struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type Directive struct {
	httpClient client.Client
}

func NewDirective(httpClient client.Client) *Directive {
	return &Directive{
		httpClient: httpClient,
	}
}

func (d *Directive) ValidateDirective(ctx context.Context, obj interface{}, next graphql.Resolver, validationType string) (interface{}, error) {
	resolvedValue, err := next(ctx)
	if err != nil {
		return nil, err
	}

	switch validationType {
	case "email":
		if email, ok := resolvedValue.(string); ok {
			if !pkg.IsValidEmail(email) {
				return nil, fmt.Errorf("invalid email format")
			}
		} else {
			return resolvedValue, nil
		}
	case "name":
		if name, ok := resolvedValue.(string); ok {
			if !pkg.IsValidName(name) {
				return nil, fmt.Errorf("invalid name format")
			}
		} else {
			return resolvedValue, nil
		}
	default:
		return nil, fmt.Errorf("unsupported validation type: %s", validationType)
	}
	return resolvedValue, nil
}
