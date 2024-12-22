package list

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/generated/graphql"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/client"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/converters"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/jwt"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"net/http"
)

type Resolver struct {
	httpClient client.Client
	listConv   converters.ListConverter
	userConv   converters.UserConverter
}

func NewResolver(client client.Client, listConverter converters.ListConverter, userConverter converters.UserConverter) *Resolver {
	return &Resolver{
		httpClient: client,
		listConv:   listConverter,
		userConv:   userConverter,
	}
}

func (r *Resolver) ListsGlobal(ctx context.Context) ([]*graphql.List, error) {
	log.C(ctx).Info("list resolver for getting all lists")
	response, err := r.httpClient.Do(ctx, http.MethodGet, "/lists/all", nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch all lists: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	var lists []*models.List
	if err = json.Unmarshal(response, &lists); err != nil {
		log.C(ctx).Errorf("failed to unmarshal list response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	result, err := r.listConv.ConvertMultipleListsToGraphQL(lists)
	if err != nil {
		log.C(ctx).Errorf("failed to convert multiple lists to graphql: %v", err)
		return nil, fmt.Errorf("error while converting multiple lists to graphql: %w", err)
	}
	log.C(ctx).Debugf("converted lists: %v", result)
	return result, nil
}

func (r *Resolver) Lists(ctx context.Context) ([]*graphql.List, error) {
	log.C(ctx).Info("list resolver for ListByUser")
	url := "/lists/user/all"

	response, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch list response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("list response: %v", string(response))

	var lists []*models.List
	if err = json.Unmarshal(response, &lists); err != nil {
		log.C(ctx).Errorf("failed to unmarshal list response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	result, err := r.listConv.ConvertMultipleListsToGraphQL(lists)
	if err != nil {
		log.C(ctx).Errorf("failed to convert multiple lists to graphql: %v", err)
		return nil, fmt.Errorf("error converting multiple lists t graphql: %w", err)
	}
	return result, nil
}

func (r *Resolver) ListsAccepted(ctx context.Context) ([]*graphql.List, error) {
	log.C(ctx).Info("list resolver for accepted lists")
	url := "/lists/user/accepted"

	response, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch list response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("list response: %v", string(response))

	var lists []*models.List
	if err = json.Unmarshal(response, &lists); err != nil {
		log.C(ctx).Errorf("failed to unmarshal list response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	result, err := r.listConv.ConvertMultipleListsToGraphQL(lists)
	if err != nil {
		log.C(ctx).Errorf("failed to convert multiple lists to graphql: %v", err)
		return nil, fmt.Errorf("error converting multiple lists t graphql: %w", err)
	}
	return result, nil
}

func (r *Resolver) ListsPending(ctx context.Context) ([]*graphql.List, error) {
	log.C(ctx).Info("list resolver for pending list for a user")
	url := "/lists/pending/all"

	response, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch pending list response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("pending list response: %v", string(response))

	var lists []*models.List
	if err = json.Unmarshal(response, &lists); err != nil {
		log.C(ctx).Errorf("failed to unmarshal pending list response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	result, err := r.listConv.ConvertMultipleListsToGraphQL(lists)
	if err != nil {
		log.C(ctx).Errorf("failed to convert multiple pending lists to graphql: %v", err)
		return nil, fmt.Errorf("error converting multiple pending lists t graphql: %w", err)
	}
	return result, nil
}

func (r *Resolver) List(ctx context.Context, id string) (*graphql.List, error) {
	log.C(ctx).Info("list resolver for getting list by id")
	url := fmt.Sprintf("/lists/%s", id)
	response, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch list by id: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	var l models.List
	if err = json.Unmarshal(response, &l); err != nil {
		log.C(ctx).Errorf("failed to unmarshal list response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	return r.listConv.ConvertListToGraphQL(l)
}

func (r *Resolver) Owner(ctx context.Context, obj *graphql.List) (*graphql.User, error) {
	log.C(ctx).Info("list resolver for getting user by id")
	if obj == nil {
		log.C(ctx).Error("list is nil")
		return nil, nil
	}
	var ownerID string
	url := fmt.Sprintf("/lists/%s/owner", obj.ID)
	body, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch user by id: %v", err)
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	log.C(ctx).Debugf("user: %v", string(body))

	if err = json.Unmarshal(body, &ownerID); err != nil {
		log.C(ctx).Errorf("failed to unmarshal user response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return r.getUser(ctx, ownerID)
}

func (r *Resolver) Todos(ctx context.Context, obj *graphql.List) ([]*graphql.Todo, error) {
	log.C(ctx).Info("list resolver for getting todos")
	if obj == nil {
		return nil, nil
	}
	url := fmt.Sprintf("/lists/%s/todos", obj.ID)
	body, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch todos: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("todos: %v", string(body))
	var todos []models.Todo
	if err = json.Unmarshal(body, &todos); err != nil {
		log.C(ctx).Errorf("failed to unmarshal todos: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	var result []*graphql.Todo
	for _, el := range todos {
		t, err := converters.NewConverterTodoGraphQL().ConvertTodoToGraphQL(el)

		if err != nil {
			log.C(ctx).Errorf("failed to convert todo: %v", err)
			return nil, fmt.Errorf("error converting todo: %w", err)
		}
		result = append(result, t)
	}
	return result, err
}

func (r *Resolver) Collaborators(ctx context.Context, obj *graphql.List) ([]*graphql.ListAccess, error) {
	log.C(ctx).Info("list resolver for getting collaborators")
	if obj == nil {
		return nil, nil
	}
	url := fmt.Sprintf("/lists/%s/users", obj.ID)
	body, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch collaborators: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	var access []models.Access
	if err = json.Unmarshal(body, &access); err != nil {
		log.C(ctx).Errorf("failed to unmarshal collaborators: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	var result []*graphql.ListAccess
	for _, el := range access {
		list, err := r.getList(ctx, el.ListID)
		if err != nil {
			log.C(ctx).Errorf("failed to fetch collaborators: %v", err)
			return nil, fmt.Errorf("error converting list: %w", err)
		}
		log.C(ctx).Debugf("collaborators list: %v", list)
		user, err := r.getUser(ctx, el.UserID)
		if err != nil {
			log.C(ctx).Errorf("failed to fetch collaborators: %v", err)
			return nil, fmt.Errorf("error converting user: %w", err)
		}
		log.C(ctx).Debugf("collaborators user: %v", user)
		role, err := r.listConv.ConvertAccessLevelToGraphQL(el.Role)
		if err != nil {
			log.C(ctx).Errorf("failed to convert access level: %v", err)
			return nil, fmt.Errorf("error converting role: %w", err)
		}
		log.C(ctx).Debugf("collaborators role: %v", role)
		result = append(result, &graphql.ListAccess{
			List:        list,
			User:        user,
			AccessLevel: role,
			Status:      &el.Status,
		})
		log.C(ctx).Debugf("collaborators result: %v", result)
	}
	return result, err
}

func (r *Resolver) CreateList(ctx context.Context, input graphql.CreateListInput) (*graphql.List, error) {
	log.C(ctx).Info("create list resolver")
	claim, ok := ctx.Value("user").(*jwt.Claims)
	if !ok {
		log.C(ctx).Error("failed to get claim from the context")
		return nil, fmt.Errorf("failed to get claim from the context")
	}
	httpInput, err := r.listConv.ConvertCreateListInput(input, claim.ID)

	if err != nil {
		log.C(ctx).Errorf("failed to convert create list: %v", err)
		return nil, fmt.Errorf("error converting create list input to struct: %w", err)
	}
	body, err := json.Marshal(httpInput)
	if err != nil {
		log.C(ctx).Errorf("failed to marshal create list input: %v", err)
		return nil, fmt.Errorf("error marshalling user: %v", err)
	}
	responsePostBody, err := r.httpClient.Do(ctx, http.MethodPost, "/lists/create", body)
	log.C(ctx).Debugf("create list response: %v", responsePostBody)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch create list response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	var id string

	if err = json.Unmarshal(responsePostBody, &id); err != nil {
		log.C(ctx).Errorf("failed to unmarshal create list response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	log.C(ctx).Debugf("create list id: %v", id)

	url := fmt.Sprintf("/lists/%s", id)

	responseGet, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch create list response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	var list models.List

	if err = json.Unmarshal(responseGet, &list); err != nil {
		log.C(ctx).Errorf("failed to unmarshal create list response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	log.C(ctx).Debugf("create list id: %v", list.ID)

	return r.listConv.ConvertListToGraphQL(list)
}

func (r *Resolver) UpdateList(ctx context.Context, id string, input graphql.UpdateListInput) (*graphql.List, error) {
	log.C(ctx).Info("update list resolver")
	httpInput, err := r.listConv.ConvertUpdateListInput(input)
	if err != nil {
		log.C(ctx).Errorf("failed to convert update list: %v", err)
		return nil, fmt.Errorf("error converting update list input to struct: %w", err)
	}
	body, err := json.Marshal(httpInput)
	if err != nil {
		log.C(ctx).Errorf("failed to marshal update list input: %v", err)
		return nil, fmt.Errorf("error marshalling user: %v", err)
	}

	url := fmt.Sprintf("/lists/%s", id)
	log.C(ctx).Debugf("update list response: %v", string(body))

	_, err = r.httpClient.Do(ctx, http.MethodPut, url, body)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch update list response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	responseGet, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch update list response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("update list id: %v", id)

	var l models.List

	if err = json.Unmarshal(responseGet, &l); err != nil {
		log.C(ctx).Errorf("failed to unmarshal update list response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	log.C(ctx).Debugf("update list: %v", l)

	return r.listConv.ConvertListToGraphQL(l)
}

func (r *Resolver) AcceptList(ctx context.Context, listID string) (*bool, error) {
	log.C(ctx).Info("accept list resolver")
	claim, ok := ctx.Value("user").(*jwt.Claims)
	if !ok {
		log.C(ctx).Error("failed to get claim from the context")
		return nil, fmt.Errorf("failed to get claim from the context")
	}
	url := fmt.Sprintf("/lists_access/%s/%s", listID, claim.ID)

	_, err := r.httpClient.Do(ctx, http.MethodPost, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to accept list response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	success := true
	return &success, nil
}

func (r *Resolver) DeleteList(ctx context.Context, id string) (*graphql.List, error) {
	log.C(ctx).Info("delete list resolver")
	url := fmt.Sprintf("/lists/%s", id)

	response, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch get list response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("delete list response: %v", string(response))

	var l models.List
	if err = json.Unmarshal(response, &l); err != nil {
		log.C(ctx).Errorf("failed to unmarshal delete list response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	_, err = r.httpClient.Do(ctx, http.MethodDelete, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch delete list response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return r.listConv.ConvertListToGraphQL(l)
}

func (r *Resolver) ListsByUser(ctx context.Context, id string) ([]*graphql.List, error) {
	log.C(ctx).Info("list resolver for ListByUser")
	url := fmt.Sprintf("/users/%s/lists", id)

	response, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch list response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("list response: %v", string(response))

	var lists []*models.List
	if err = json.Unmarshal(response, &lists); err != nil {
		log.C(ctx).Errorf("failed to unmarshal list response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	result, err := r.listConv.ConvertMultipleListsToGraphQL(lists)
	if err != nil {
		log.C(ctx).Errorf("failed to convert multiple lists to graphql response: %v", err)
		return nil, fmt.Errorf("error converting lists response: %w", err)
	}
	return result, nil
}

func (r *Resolver) getUser(ctx context.Context, id string) (*graphql.User, error) {
	log.C(ctx).Info("get user resolver")
	url := fmt.Sprintf("/users/%s", id)
	body, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch user response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	var u models.User
	if err = json.Unmarshal(body, &u); err != nil {
		log.C(ctx).Errorf("failed to unmarshal user response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	graphqlUser, err := r.userConv.ConvertUserToGraphQL(u)
	log.C(ctx).Debugf("user grapjql: %v", graphqlUser)
	if err != nil {
		log.C(ctx).Errorf("failed to convert user response: %v", err)
		return nil, fmt.Errorf("error converting user: %w", err)
	}
	return graphqlUser, nil
}

func (r *Resolver) AddListAccess(ctx context.Context, input graphql.GrantListAccessInput) (*graphql.ListAccess, error) {
	log.C(ctx).Info("add list access resolver")
	httpInput, err := r.listConv.ConvertGrantListAccessInputToModel(input)
	if err != nil {
		log.C(ctx).Errorf("failed to convert grant list access input to model: %v", err)
		return nil, fmt.Errorf("error converting create list input to struct: %w", err)
	}
	body, err := json.Marshal(httpInput)
	if err != nil {
		log.C(ctx).Errorf("failed to marshal add list access input: %v", err)
		return nil, fmt.Errorf("error marshalling user: %v", err)
	}
	log.C(ctx).Debugf("add list access response body: %v", string(body))
	url := fmt.Sprintf("/lists_access/create/%s/%s", input.ListID, input.UserID)

	response, err := r.httpClient.Do(ctx, http.MethodPost, url, body)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch add list access response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	log.C(ctx).Debugf("add list access response: %v", string(response))

	var access models.Access

	if err = json.Unmarshal(response, &access); err != nil {
		log.C(ctx).Errorf("failed to unmarshal add list access response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	list, err := r.getList(ctx, access.ListID)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch add list access list: %v", err)
		return nil, fmt.Errorf("error converting list: %w", err)
	}
	log.C(ctx).Debugf("add list access list: %v", list)
	user, err := r.getUser(ctx, access.UserID)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch add list access user: %v", err)
		return nil, fmt.Errorf("error converting user: %w", err)
	}
	log.C(ctx).Debugf("add list access user: %v", user)
	role, err := r.listConv.ConvertAccessLevelToGraphQL(access.Role)
	if err != nil {
		log.C(ctx).Errorf("failed to convert access level: %v", err)
		return nil, fmt.Errorf("error converting role: %w", err)
	}
	log.C(ctx).Debugf("add list access role: %v", role)

	return &graphql.ListAccess{
		List:        list,
		User:        user,
		AccessLevel: role,
		Status:      &access.Status,
	}, nil
}

func (r *Resolver) RemoveListAccess(ctx context.Context, listID string) (*graphql.ListAccess, error) {
	log.C(ctx).Info("remove list access resolver")
	claim, ok := ctx.Value("user").(*jwt.Claims)
	if !ok {
		log.C(ctx).Error("failed to get claim from the context")
		return nil, fmt.Errorf("failed to get claim from the context")
	}
	url := fmt.Sprintf("/lists_access/%s/%s", listID, claim.ID)

	response, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch remove list access response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	var access models.Access
	if err = json.Unmarshal(response, &access); err != nil {
		log.C(ctx).Errorf("failed to unmarshal remove list access response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	list, err := r.getList(ctx, access.ListID)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch remove list access list: %v", err)
		return nil, fmt.Errorf("error converting list: %w", err)
	}

	_, err = r.httpClient.Do(ctx, http.MethodDelete, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch remove list access response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	log.C(ctx).Debugf("remove list access list: %v", list)
	user, err := r.getUser(ctx, access.UserID)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch remove list access user: %v", err)
		return nil, fmt.Errorf("error converting user: %w", err)
	}
	log.C(ctx).Debugf("remove list access user: %v", user)
	role, err := r.listConv.ConvertAccessLevelToGraphQL(access.Role)
	if err != nil {
		log.C(ctx).Errorf("failed to convert access level: %v", err)
		return nil, fmt.Errorf("error converting role: %w", err)
	}
	log.C(ctx).Debugf("remove list access role: %v", role)

	return &graphql.ListAccess{
		List:        list,
		User:        user,
		AccessLevel: role,
		Status:      &access.Status,
	}, nil
}

func (r *Resolver) getList(ctx context.Context, id string) (*graphql.List, error) {
	log.C(ctx).Info("get list resolver")
	url := fmt.Sprintf("/lists/%s", id)

	body, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch list response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("get list response: %v", string(body))
	var l models.List
	if err = json.Unmarshal(body, &l); err != nil {
		log.C(ctx).Errorf("failed to unmarshal list response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	return r.listConv.ConvertListToGraphQL(l)
}

func (r *Resolver) UpdateListName(ctx context.Context, id string, name string) (*graphql.List, error) {
	log.C(ctx).Info("update list name resolver")
	url := fmt.Sprintf("/lists/%s/name", id)

	var updateData struct {
		Name string `json:"name"`
	}

	updateData.Name = name

	body, err := json.Marshal(updateData)
	if err != nil {
		log.C(ctx).Errorf("failed to marshal name for the list: %v", err)
		return nil, fmt.Errorf("error marshalling name: %v", err)
	}

	response, err := r.httpClient.Do(ctx, http.MethodPatch, url, body)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch list response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("list response: %v", string(response))

	var l models.List
	if err = json.Unmarshal(response, &l); err != nil {
		log.C(ctx).Errorf("failed to unmarshal list response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	return r.listConv.ConvertListToGraphQL(l)
}

func (r *Resolver) UpdateListDescription(ctx context.Context, id string, description string) (*graphql.List, error) {
	log.C(ctx).Info("update list description resolver")
	url := fmt.Sprintf("/lists/%s/description", id)

	var updateData struct {
		Description string `json:"description"`
	}

	updateData.Description = description

	body, err := json.Marshal(updateData)
	if err != nil {
		log.C(ctx).Errorf("failed to marshal description for the list: %v", err)
		return nil, fmt.Errorf("error marshalling description: %v", err)
	}

	response, err := r.httpClient.Do(ctx, http.MethodPatch, url, body)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch list response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("list response: %v", string(response))

	var l models.List
	if err = json.Unmarshal(response, &l); err != nil {
		log.C(ctx).Errorf("failed to unmarshal list response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	return r.listConv.ConvertListToGraphQL(l)
}

func (r *Resolver) RemoveCollaborator(ctx context.Context, listID string, userID string) (*graphql.ListAccess, error) {
	log.C(ctx).Info("remove collaborator resolver")
	url := fmt.Sprintf("/lists_access/%s/%s", listID, userID)

	response, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch remove collaborator response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	var access models.Access
	if err = json.Unmarshal(response, &access); err != nil {
		log.C(ctx).Errorf("failed to unmarshal remove collaborator response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	_, err = r.httpClient.Do(ctx, http.MethodDelete, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch remove collaborator response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	list, err := r.getList(ctx, access.ListID)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch remove collaborator list: %v", err)
		return nil, fmt.Errorf("error converting list: %w", err)
	}
	log.C(ctx).Debugf("remove collaborator list: %v", list)
	user, err := r.getUser(ctx, access.UserID)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch remove collaborator user: %v", err)
		return nil, fmt.Errorf("error converting user: %w", err)
	}
	log.C(ctx).Debugf("remove collaborator user: %v", user)
	role, err := r.listConv.ConvertAccessLevelToGraphQL(access.Role)
	if err != nil {
		log.C(ctx).Errorf("failed to convert access level: %v", err)
		return nil, fmt.Errorf("error converting role: %w", err)
	}
	log.C(ctx).Debugf("remove list access role: %v", role)

	return &graphql.ListAccess{
		List:        list,
		User:        user,
		AccessLevel: role,
		Status:      &access.Status,
	}, nil
}

func (r *Resolver) GetListAccesses(ctx context.Context, listID string) ([]*graphql.ListAccess, error) {
	log.C(ctx).Info("get list accesses byt list id resolver")

	url := fmt.Sprintf("/lists_access/list/%s", listID)

	response, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("failed to fetch remove list access response: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	var accesses []models.Access
	if err = json.Unmarshal(response, &accesses); err != nil {
		log.C(ctx).Errorf("failed to unmarshal remove list access response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	var listAccesses []*graphql.ListAccess
	for _, access := range accesses {
		list, err := r.getList(ctx, access.ListID)
		if err != nil {
			log.C(ctx).Errorf("failed to fetch remove list access list: %v", err)
			return nil, fmt.Errorf("error converting list: %w", err)
		}
		user, err := r.getUser(ctx, access.UserID)
		if err != nil {
			log.C(ctx).Errorf("failed to fetch remove list access user: %v", err)
			return nil, fmt.Errorf("error converting user: %w", err)
		}
		role, err := r.listConv.ConvertAccessLevelToGraphQL(access.Role)
		if err != nil {
			log.C(ctx).Errorf("failed to convert access level: %v", err)
			return nil, fmt.Errorf("error converting role: %w", err)
		}
		listAccesses = append(listAccesses, &graphql.ListAccess{
			List:        list,
			User:        user,
			AccessLevel: role,
			Status:      &access.Status,
		})
	}

	return listAccesses, nil
}
