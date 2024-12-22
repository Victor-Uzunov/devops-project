package todo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/generated/graphql"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/client"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/converters"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"net/http"
)

type Resolver struct {
	httpClient client.Client
	todoConv   converters.TodoConverter
	listConv   converters.ListConverter
	userConv   converters.UserConverter
}

func NewResolver(
	client client.Client,
	converter converters.TodoConverter,
	listConverter converters.ListConverter,
	userConverter converters.UserConverter,
) *Resolver {
	return &Resolver{
		httpClient: client,
		todoConv:   converter,
		listConv:   listConverter,
		userConv:   userConverter,
	}
}

func (r *Resolver) TodosGlobal(ctx context.Context) ([]*graphql.Todo, error) {
	log.C(ctx).Info("todoResolver called todos global")
	response, err := r.httpClient.Do(ctx, http.MethodGet, "/todos/all", nil)
	if err != nil {
		log.C(ctx).Errorf("error getting todos: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	var todos []*models.Todo
	if err = json.Unmarshal(response, &todos); err != nil {
		log.C(ctx).Errorf("error unmarshalling todos: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	var result, errConverting = r.todoConv.ConvertMultipleTodoToGraphQL(todos)
	if errConverting != nil {
		log.C(ctx).Errorf("failed converting multiple todos to graphql: %v", err)
		return nil, fmt.Errorf("error while converting multiple todos to graphql: %w", err)
	}
	return result, nil
}

func (r *Resolver) Todos(ctx context.Context) ([]*graphql.Todo, error) {
	log.C(ctx).Info("todoResolver called todos")
	response, err := r.httpClient.Do(ctx, http.MethodGet, "/todos/user/all", nil)
	if err != nil {
		log.C(ctx).Errorf("error getting todos: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	var todos []*models.Todo
	if err = json.Unmarshal(response, &todos); err != nil {
		log.C(ctx).Errorf("error unmarshalling todos: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	var result, errConverting = r.todoConv.ConvertMultipleTodoToGraphQL(todos)
	if errConverting != nil {
		log.C(ctx).Errorf("failed converting multiple todos to graphql: %v", err)
		return nil, fmt.Errorf("error while converting multiple todos to graphql: %w", err)
	}
	return result, nil
}

func (r *Resolver) Todo(ctx context.Context, id string) (*graphql.Todo, error) {
	log.C(ctx).Info("todoResolver called todo")
	url := fmt.Sprintf("/todos/%s", id)

	response, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("error getting todo: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	var todo models.Todo
	if err = json.Unmarshal(response, &todo); err != nil {
		log.C(ctx).Errorf("error unmarshalling response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	graphTodo, err := r.todoConv.ConvertTodoToGraphQL(todo)
	log.C(ctx).Debugf("converted todo: %v", graphTodo)
	if err != nil {
		log.C(ctx).Errorf("error converting todos: %v", err)
		return nil, fmt.Errorf("error converting todo: %w", err)
	}
	return graphTodo, nil
}

func (r *Resolver) List(ctx context.Context, obj *graphql.Todo) (*graphql.List, error) {
	log.C(ctx).Info("todoResolver called list")
	if obj == nil {
		return nil, nil
	}
	todo, err := r.getTodo(ctx, obj.ID)
	if err != nil {
		log.C(ctx).Errorf("error getting todo: %v", err)
		return nil, err
	}
	url := fmt.Sprintf("/lists/%s", todo.ListID)
	body, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	log.C(ctx).Debugf("body: %v", string(body))
	if err != nil {
		log.C(ctx).Errorf("error getting todo: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	var l models.List

	if err = json.Unmarshal(body, &l); err != nil {
		log.C(ctx).Errorf("error unmarshalling response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	return r.listConv.ConvertListToGraphQL(l)
}

func (r *Resolver) AssignedTo(ctx context.Context, obj *graphql.Todo) (*graphql.User, error) {
	log.C(ctx).Info("todoResolver called for assigned to")
	if obj == nil {
		return nil, nil
	}
	todo, err := r.getTodo(ctx, obj.ID)
	if err != nil {
		log.C(ctx).Errorf("error getting todo: %v", err)
		return nil, err
	}
	if todo.AssignedTo == nil {
		return nil, nil
	}
	log.C(ctx).Debugf("todo: %v", todo)
	url := fmt.Sprintf("/users/%s", *todo.AssignedTo)
	body, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("error getting todo: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("body: %v", string(body))

	var u models.User
	if err = json.Unmarshal(body, &u); err != nil {
		log.C(ctx).Errorf("error unmarshalling response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	graphqlUser, err := r.userConv.ConvertUserToGraphQL(u)
	if err != nil {
		log.C(ctx).Errorf("error converting users: %v", err)
		return nil, fmt.Errorf("error converting user: %w", err)
	}
	log.C(ctx).Debugf("converted user to graphql model: %v", graphqlUser)

	return graphqlUser, nil
}

func (r *Resolver) CreateTodo(ctx context.Context, input graphql.CreateTodoInput) (*graphql.Todo, error) {
	log.C(ctx).Info("todoResolver called")
	httpInput, err := r.todoConv.ConvertCreateTodoInput(input)
	if err != nil {
		log.C(ctx).Errorf("error converting todos: %v", err)
		return nil, fmt.Errorf("error converting create todo input to struct: %w", err)
	}
	body, err := json.Marshal(httpInput)
	if err != nil {
		log.C(ctx).Errorf("error marshalling response: %v", err)
		return nil, fmt.Errorf("error marshalling user: %v", err)
	}
	log.C(ctx).Debugf("body: %v", string(body))

	response, err := r.httpClient.Do(ctx, http.MethodPost, "/todos/create", body)
	if err != nil {
		log.C(ctx).Errorf("error executing request: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("todo: %v", response)

	var id string

	if err = json.Unmarshal(response, &id); err != nil {
		log.C(ctx).Errorf("error unmarshalling response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	url := fmt.Sprintf("/todos/%s", id)
	log.C(ctx).Debugf("id of the todo: %s", id)

	responseGet, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("error getting todo: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("todo response get: %v", responseGet)

	var todo models.Todo
	if err = json.Unmarshal(responseGet, &todo); err != nil {
		log.C(ctx).Errorf("error unmarshalling response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return r.todoConv.ConvertTodoToGraphQL(todo)
}

func (r *Resolver) UpdateTodo(ctx context.Context, id string, input graphql.UpdateTodoInput) (*graphql.Todo, error) {
	log.C(ctx).Info("todoResolver called update")
	httpInput, err := r.todoConv.ConvertUpdateTodoInput(input)
	if err != nil {
		log.C(ctx).Errorf("error converting todos: %v", err)
		return nil, fmt.Errorf("error converting update todo input to struct: %w", err)
	}
	body, err := json.Marshal(httpInput)
	if err != nil {
		log.C(ctx).Errorf("error marshalling response: %v", err)
		return nil, fmt.Errorf("error marshalling user: %v", err)
	}
	log.C(ctx).Debugf("body: %v", string(body))

	url := fmt.Sprintf("/todos/%s", id)

	_, err = r.httpClient.Do(ctx, http.MethodPut, url, body)
	if err != nil {
		log.C(ctx).Errorf("error executing request: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	responseGet, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("error getting todo: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("todo response get: %v", responseGet)

	var t models.Todo

	if err = json.Unmarshal(responseGet, &t); err != nil {
		log.C(ctx).Errorf("error unmarshalling response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return r.todoConv.ConvertTodoToGraphQL(t)
}

func (r *Resolver) DeleteTodo(ctx context.Context, id string) (*graphql.Todo, error) {
	log.C(ctx).Info("todoResolver called delete")
	url := fmt.Sprintf("/todos/%s", id)

	response, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("error getting todo: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	log.C(ctx).Debugf("todo response: %v", response)

	var t models.Todo
	if err = json.Unmarshal(response, &t); err != nil {
		log.C(ctx).Errorf("error unmarshalling response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	_, err = r.httpClient.Do(ctx, http.MethodDelete, url, nil)
	if err != nil {
		log.C(ctx).Errorf("error executing request: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return r.todoConv.ConvertTodoToGraphQL(t)
}

func (r *Resolver) TodosByList(ctx context.Context, id string) ([]*graphql.Todo, error) {
	log.C(ctx).Infof("todoResolver for TodoByList is called")
	url := fmt.Sprintf("/lists/%s/todos", id)

	response, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("error getting todos for listID %s: %v", id, err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	var todos []*models.Todo
	if err = json.Unmarshal(response, &todos); err != nil {
		log.C(ctx).Errorf("error unmarshalling response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	result, err := r.todoConv.ConvertMultipleTodoToGraphQL(todos)
	if err != nil {
		log.C(ctx).Errorf("failed converting multiple todos to graphql: %v", err)
		return nil, fmt.Errorf("error while converting multiple todos to graphql: %w", err)
	}
	return result, nil
}

func (r *Resolver) UpdateTodoTitle(ctx context.Context, id string, title string) (*graphql.Todo, error) {
	log.C(ctx).Info("todoResolver called update todo title")
	url := fmt.Sprintf("/todos/%s/title", id)

	var updateData struct {
		Title string `json:"title"`
	}

	updateData.Title = title

	body, err := json.Marshal(updateData)
	if err != nil {
		log.C(ctx).Errorf("failed to marshal title for the todo: %v", err)
		return nil, fmt.Errorf("error marshalling title: %v", err)
	}

	response, err := r.httpClient.Do(ctx, http.MethodPatch, url, body)
	if err != nil {
		log.C(ctx).Errorf("error getting todo: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	var todo models.Todo
	if err = json.Unmarshal(response, &todo); err != nil {
		log.C(ctx).Errorf("error unmarshalling response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	graphTodo, err := r.todoConv.ConvertTodoToGraphQL(todo)
	if err != nil {
		log.C(ctx).Errorf("error converting todos: %v", err)
		return nil, fmt.Errorf("error converting todo: %w", err)
	}
	log.C(ctx).Debugf("converted todo: %v", graphTodo)
	return graphTodo, nil
}

func (r *Resolver) UpdateTodoDescription(ctx context.Context, id string, description string) (*graphql.Todo, error) {
	log.C(ctx).Info("todoResolver called update todo description")
	url := fmt.Sprintf("/todos/%s/description", id)

	var updateData struct {
		Description string `json:"description"`
	}

	updateData.Description = description

	body, err := json.Marshal(updateData)
	if err != nil {
		log.C(ctx).Errorf("failed to marshal description for the todo: %v", err)
		return nil, fmt.Errorf("error marshalling description: %v", err)
	}

	response, err := r.httpClient.Do(ctx, http.MethodPatch, url, body)
	if err != nil {
		log.C(ctx).Errorf("error getting todo: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	var todo models.Todo
	if err = json.Unmarshal(response, &todo); err != nil {
		log.C(ctx).Errorf("error unmarshalling response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	graphTodo, err := r.todoConv.ConvertTodoToGraphQL(todo)
	if err != nil {
		log.C(ctx).Errorf("error converting todos: %v", err)
		return nil, fmt.Errorf("error converting todo: %w", err)
	}
	log.C(ctx).Debugf("converted todo: %v", graphTodo)
	return graphTodo, nil
}

func (r *Resolver) UpdateTodoPriority(ctx context.Context, id string, priority graphql.Priority) (*graphql.Todo, error) {
	log.C(ctx).Info("todoResolver called update todo priority")
	url := fmt.Sprintf("/todos/%s/priority", id)

	var updateData struct {
		Priority graphql.Priority `json:"priority"`
	}

	updateData.Priority = priority

	body, err := json.Marshal(updateData)
	if err != nil {
		log.C(ctx).Errorf("failed to marshal priority for the todo: %v", err)
		return nil, fmt.Errorf("error marshalling priority: %v", err)
	}

	response, err := r.httpClient.Do(ctx, http.MethodPatch, url, body)
	if err != nil {
		log.C(ctx).Errorf("error getting todo: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	var todo models.Todo
	if err = json.Unmarshal(response, &todo); err != nil {
		log.C(ctx).Errorf("error unmarshalling response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	graphTodo, err := r.todoConv.ConvertTodoToGraphQL(todo)
	if err != nil {
		log.C(ctx).Errorf("error converting todos: %v", err)
		return nil, fmt.Errorf("error converting todo: %w", err)
	}
	log.C(ctx).Debugf("converted todo: %v", graphTodo)
	return graphTodo, nil
}

func (r *Resolver) UpdateTodoAssignTo(ctx context.Context, id string, userID string) (*graphql.Todo, error) {
	log.C(ctx).Info("todoResolver called update todo assigned to user")
	url := fmt.Sprintf("/todos/%s/assign_to", id)

	var updateData struct {
		UserID string `json:"user_id"`
	}

	updateData.UserID = userID

	body, err := json.Marshal(updateData)
	if err != nil {
		log.C(ctx).Errorf("failed to marshal userID for the todo: %v", err)
		return nil, fmt.Errorf("error marshalling userID: %v", err)
	}

	response, err := r.httpClient.Do(ctx, http.MethodPatch, url, body)
	if err != nil {
		log.C(ctx).Errorf("error getting todo: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	var todo models.Todo
	if err = json.Unmarshal(response, &todo); err != nil {
		log.C(ctx).Errorf("error unmarshalling response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	graphTodo, err := r.todoConv.ConvertTodoToGraphQL(todo)
	if err != nil {
		log.C(ctx).Errorf("error converting todos: %v", err)
		return nil, fmt.Errorf("error converting todo: %w", err)
	}
	log.C(ctx).Debugf("converted todo: %v", graphTodo)
	return graphTodo, nil
}

func (r *Resolver) CompleteTodo(ctx context.Context, id string) (*graphql.Todo, error) {
	log.C(ctx).Info("todoResolver called complete todo")
	url := fmt.Sprintf("/todos/%s/complete", id)

	response, err := r.httpClient.Do(ctx, http.MethodPatch, url, nil)
	if err != nil {
		log.C(ctx).Errorf("error getting todo: %v", err)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	var todo models.Todo
	if err = json.Unmarshal(response, &todo); err != nil {
		log.C(ctx).Errorf("error unmarshalling response: %v", err)
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	graphTodo, err := r.todoConv.ConvertTodoToGraphQL(todo)
	if err != nil {
		log.C(ctx).Errorf("error converting todos: %v", err)
		return nil, fmt.Errorf("error converting todo: %w", err)
	}
	log.C(ctx).Debugf("converted todo: %v", graphTodo)
	return graphTodo, nil
}

func (r *Resolver) getTodo(ctx context.Context, todoID string) (models.Todo, error) {
	log.C(ctx).Info("todoResolver called get")
	url := fmt.Sprintf("/todos/%s", todoID)
	body, err := r.httpClient.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.C(ctx).Errorf("error getting todo: %v", err)
		return models.Todo{}, fmt.Errorf("error executing request: %w", err)
	}
	var todo models.Todo
	if err = json.Unmarshal(body, &todo); err != nil {
		log.C(ctx).Errorf("error unmarshalling response: %v", err)
		return models.Todo{}, fmt.Errorf("error unmarshalling response: %w", err)
	}
	return todo, nil
}
