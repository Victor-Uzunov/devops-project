package resolvers

import (
	"context"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/generated/graphql"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
)

type mutationResolver struct {
	*RootResolver
}

func (r *mutationResolver) CreateUser(ctx context.Context, input graphql.CreateUserInput) (*graphql.User, error) {
	log.C(ctx).Info("creating user mutation resolver")
	return r.user.CreateUser(ctx, input)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input graphql.UpdateUserInput) (*graphql.User, error) {
	log.C(ctx).Info("updating user mutation resolver")
	return r.user.UpdateUser(ctx, id, input)
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*graphql.User, error) {
	log.C(ctx).Info("deleting user mutation resolver")
	return r.user.DeleteUser(ctx, id)
}

func (r *mutationResolver) CreateList(ctx context.Context, input graphql.CreateListInput) (*graphql.List, error) {
	log.C(ctx).Info("creating list mutation resolver")
	return r.list.CreateList(ctx, input)
}

func (r *mutationResolver) UpdateList(ctx context.Context, id string, input graphql.UpdateListInput) (*graphql.List, error) {
	log.C(ctx).Info("updating list mutation resolver")
	return r.list.UpdateList(ctx, id, input)
}

func (r *mutationResolver) DeleteList(ctx context.Context, id string) (*graphql.List, error) {
	log.C(ctx).Info("deleting list mutation resolver")
	return r.list.DeleteList(ctx, id)
}

func (r *mutationResolver) CreateTodo(ctx context.Context, input graphql.CreateTodoInput) (*graphql.Todo, error) {
	log.C(ctx).Info("creating todo mutation resolver")
	return r.todo.CreateTodo(ctx, input)
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, id string, input graphql.UpdateTodoInput) (*graphql.Todo, error) {
	log.C(ctx).Info("updating todo mutation resolver")
	return r.todo.UpdateTodo(ctx, id, input)
}

func (r *mutationResolver) DeleteTodo(ctx context.Context, id string) (*graphql.Todo, error) {
	log.C(ctx).Info("deleting todo mutation resolver")
	return r.todo.DeleteTodo(ctx, id)
}

func (r *mutationResolver) AddListAccess(ctx context.Context, input graphql.GrantListAccessInput) (*graphql.ListAccess, error) {
	log.C(ctx).Info("adding list access mutation resolver")
	return r.list.AddListAccess(ctx, input)
}

func (r *mutationResolver) RemoveListAccess(ctx context.Context, listID string) (*graphql.ListAccess, error) {
	log.C(ctx).Info("removing list access mutation resolver")
	return r.list.RemoveListAccess(ctx, listID)
}

func (r *mutationResolver) UpdateListName(ctx context.Context, id string, name string) (*graphql.List, error) {
	log.C(ctx).Info("updating list name mutation resolver")
	return r.list.UpdateListName(ctx, id, name)
}

func (r *mutationResolver) UpdateListDescription(ctx context.Context, id string, description string) (*graphql.List, error) {
	log.C(ctx).Info("updating list description mutation resolver")
	return r.list.UpdateListDescription(ctx, id, description)
}

func (r *mutationResolver) UpdateTodoTitle(ctx context.Context, id string, title string) (*graphql.Todo, error) {
	log.C(ctx).Info("updating todo title mutation resolver")
	return r.todo.UpdateTodoTitle(ctx, id, title)
}

func (r *mutationResolver) UpdateTodoDescription(ctx context.Context, id string, description string) (*graphql.Todo, error) {
	log.C(ctx).Info("updating todo description mutation resolver")
	return r.todo.UpdateTodoDescription(ctx, id, description)
}

func (r *mutationResolver) UpdateTodoPriority(ctx context.Context, id string, priority graphql.Priority) (*graphql.Todo, error) {
	log.C(ctx).Info("updating todo priority mutation resolver")
	return r.todo.UpdateTodoPriority(ctx, id, priority)
}

func (r *mutationResolver) UpdateTodoAssignTo(ctx context.Context, id string, userID string) (*graphql.Todo, error) {
	log.C(ctx).Info("updating todo assignment mutation resolver")
	return r.todo.UpdateTodoAssignTo(ctx, id, userID)
}

func (r *mutationResolver) CompleteTodo(ctx context.Context, id string) (*graphql.Todo, error) {
	log.C(ctx).Info("updating todo completion mutation resolver")
	return r.todo.CompleteTodo(ctx, id)
}

func (r *mutationResolver) AcceptList(ctx context.Context, listID string) (*bool, error) {
	log.C(ctx).Info("accepting list access mutation resolver")
	return r.list.AcceptList(ctx, listID)
}

func (r *mutationResolver) RemoveCollaborator(ctx context.Context, listID string, userID string) (*graphql.ListAccess, error) {
	log.C(ctx).Info("removing collaborator mutation resolver")
	return r.list.RemoveCollaborator(ctx, listID, userID)
}
