package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/generated/graphql"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/client"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/converters"
	jwts "github.com/Victor-Uzunov/devops-project/todoservice/pkg/jwt"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"net/http"
)

type Resolver struct {
	httpClient client.Client
	userConv   converters.UserConverter
	listConv   converters.ListConverter
}

func NewResolver(client client.Client, converter converters.UserConverter, listConverter converters.ListConverter) *Resolver {
	return &Resolver{
		httpClient: client,
		userConv:   converter,
		listConv:   listConverter,
	}
}

func (r *Resolver) Users(ctx context.Context) ([]*graphql.User, error) {
	log.C(ctx).Info("users resolver users")
	response, err := r.httpClient.Do(ctx, http.MethodGet, "/users/all", nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch users: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	var users []*models.User
	if err = json.Unmarshal(response, &users); err != nil {
		log.C(ctx).Errorf("failed to unmarshal users: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	var result []*graphql.User
	for _, el := range users {
		u, err := r.userConv.ConvertUserToGraphQL(*el)
		log.C(ctx).Debugf("user: %+v", u)
		if err != nil {
			log.C(ctx).Errorf("failed to convert user: %v", err)
			return nil, fmt.Errorf("error converting user: %w", err)
		}
		result = append(result, u)
		log.C(ctx).Debugf("user result: %+v", u)
	}
	return result, nil
}

func (r *Resolver) User(ctx context.Context, id string) (*graphql.User, error) {
	log.C(ctx).Info("users resolver user", id)
	url := fmt.Sprintf("/users/%s", id)

	response, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch user: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("user response: %+v", response)
	var u models.User
	if err = json.Unmarshal(response, &u); err != nil {
		log.C(ctx).Errorf("failed to unmarshal user: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	graphqlUser, err := r.userConv.ConvertUserToGraphQL(u)
	if err != nil {
		log.C(ctx).Errorf("failed to convert user: %v", err)
		return nil, fmt.Errorf("error converting user: %w", err)
	}
	log.C(ctx).Debugf("user graphql: %+v", graphqlUser)

	return graphqlUser, nil
}

func (r *Resolver) UserByEmail(ctx context.Context) (*graphql.User, error) {
	log.C(ctx).Info("users resolver user by email")
	claims, ok := ctx.Value("user").(*jwts.Claims)
	if !ok {
		return nil, fmt.Errorf("unable to extract claims from context")
	}

	url := fmt.Sprintf("/users/email/%s", claims.Email)

	response, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch user: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("user response: %+v", response)
	var u models.User
	if err = json.Unmarshal(response, &u); err != nil {
		log.C(ctx).Errorf("failed to unmarshal user: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	graphqlUser, err := r.userConv.ConvertUserToGraphQL(u)
	if err != nil {
		log.C(ctx).Errorf("failed to convert user: %v", err)
		return nil, fmt.Errorf("error converting user: %w", err)
	}
	log.C(ctx).Debugf("user graphql: %+v", graphqlUser)

	return graphqlUser, nil
}

func (r *Resolver) CreateUser(ctx context.Context, input graphql.CreateUserInput) (*graphql.User, error) {
	log.C(ctx).Info("users resolver create user")
	httpInput, err := r.userConv.ConvertCreateUserInput(input)
	if err != nil {
		log.C(ctx).Errorf("failed to convert create user: %v", err)
		return nil, fmt.Errorf("convert create user input to struct: %w", err)
	}
	body, err := json.Marshal(httpInput)
	if err != nil {
		log.C(ctx).Errorf("failed to marshal create user: %v", err)
		return nil, fmt.Errorf("error marshalling user: %v", err)
	}

	responseID, err := r.httpClient.Do(ctx, http.MethodPost, "/users/create", body)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch user: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("user response id: %+v", responseID)

	var id string

	if err = json.Unmarshal(responseID, &id); err != nil {
		log.C(ctx).Errorf("failed to unmarshal user: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	url := fmt.Sprintf("/users/%s", id)

	responseGet, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch user: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("user response get: %+v", responseGet)

	var user models.User
	if err = json.Unmarshal(responseGet, &user); err != nil {
		log.C(ctx).Errorf("failed to unmarshal user: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return r.userConv.ConvertUserToGraphQL(user)
}

func (r *Resolver) UpdateUser(ctx context.Context, id string, input graphql.UpdateUserInput) (*graphql.User, error) {
	log.C(ctx).Info("users resolver update user", id)
	httpInput, err := r.userConv.ConvertUpdateUserInput(input)
	if err != nil {
		log.C(ctx).Errorf("failed to convert update user: %v", err)
		return nil, fmt.Errorf("convert update user input to struct: %w", err)
	}
	body, err := json.Marshal(httpInput)
	if err != nil {
		log.C(ctx).Errorf("failed to marshal update user: %v", err)
		return nil, fmt.Errorf("error marshalling user: %v", err)
	}
	log.C(ctx).Debugf("body input: %+v", body)

	url := fmt.Sprintf("/users/%s", id)
	_, err = r.httpClient.Do(ctx, http.MethodPut, url, body)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch user: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	responseGet, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch user: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	log.C(ctx).Debugf("response get: %+v", responseGet)

	var user models.User
	if err = json.Unmarshal(responseGet, &user); err != nil {
		log.C(ctx).Errorf("failed to unmarshal user: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	log.C(ctx).Debugf("user response get: %+v", user)

	return r.userConv.ConvertUserToGraphQL(user)
}

func (r *Resolver) DeleteUser(ctx context.Context, id string) (*graphql.User, error) {
	log.C(ctx).Infof("users resolver delete user with id %s", id)
	url := fmt.Sprintf("/users/%s", id)

	response, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch user: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	var u models.User
	if err = json.Unmarshal(response, &u); err != nil {
		log.C(ctx).Errorf("failed to unmarshal user: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	log.C(ctx).Infof("user response user: %+v", u)
	_, err = r.httpClient.Do(ctx, http.MethodDelete, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch user: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return r.userConv.ConvertUserToGraphQL(u)
}

func (r *Resolver) GetList(ctx context.Context, id string) (*graphql.List, error) {
	log.C(ctx).Infof("lists resolver get list with id %s", id)
	url := fmt.Sprintf("/lists/%s", id)

	body, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch list: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	var l models.List

	if err = json.Unmarshal(body, &l); err != nil {
		log.C(ctx).Errorf("failed to unmarshal list: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	return r.listConv.ConvertListToGraphQL(l)
}

func (r *Resolver) GetUser(ctx context.Context, id string) (*graphql.User, error) {
	log.C(ctx).Infof("users resolver get user with id %s", id)
	url := fmt.Sprintf("/users/%s", id)

	body, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch user: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	var u models.User

	if err = json.Unmarshal(body, &u); err != nil {
		log.C(ctx).Errorf("failed to unmarshal user: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	log.C(ctx).Debugf("user response user: %+v", u)
	return r.userConv.ConvertUserToGraphQL(u)
}

func (r *Resolver) UsersByList(ctx context.Context, id string) ([]*graphql.User, error) {
	log.C(ctx).Infof("users resolver get user for a list with id %s", id)
	url := fmt.Sprintf("/lists/%s/users", id)

	body, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch user: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	var access []models.Access
	if err = json.Unmarshal(body, &access); err != nil {
		log.C(ctx).Errorf("failed to unmarshal collaborators: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	var result []*graphql.User
	for _, el := range access {
		user, err := r.GetUser(ctx, el.UserID)
		if err != nil {
			log.C(ctx).Errorf("failed to fetch collaborators: %v", err)
			return nil, fmt.Errorf("error converting user: %w", err)
		}
		log.C(ctx).Debugf("collaborators user: %v", user)

		result = append(result, user)
	}
	return result, err
}
