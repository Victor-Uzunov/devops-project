package resolvers

import (
	"context"
	graph "github.com/Victor-Uzunov/devops-project/graphqlServer/generated"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/generated/graphql"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/client"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/converters"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/resolvers/list"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/resolvers/todo"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/internal/resolvers/user"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
)

var _ graph.ResolverRoot = &RootResolver{}

type RootResolver struct {
	list *list.Resolver
	user *user.Resolver
	todo *todo.Resolver
}

func NewRootResolver(todoService client.Client) *RootResolver {
	listConverter := converters.NewConverterListGraphQL()
	todoConverter := converters.NewConverterTodoGraphQL()
	userConverter := converters.NewConverterUserGraphQL()

	return &RootResolver{
		list: list.NewResolver(todoService, listConverter, userConverter),
		user: user.NewResolver(todoService, userConverter, listConverter),
		todo: todo.NewResolver(todoService, todoConverter, listConverter, userConverter),
	}
}

func (r *RootResolver) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}

func (r *RootResolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}

func (r *RootResolver) List() graph.ListResolver {
	return &listResolver{r}
}

func (r *RootResolver) Todo() graph.TodoResolver {
	return &todoResolver{r}
}

type todoResolver struct {
	*RootResolver
}

func (r *todoResolver) AssignedTo(ctx context.Context, obj *graphql.Todo) (*graphql.User, error) {
	log.C(ctx).Info("todoResolver.AssignedTo")
	return r.todo.AssignedTo(ctx, obj)
}

func (r *todoResolver) List(ctx context.Context, obj *graphql.Todo) (*graphql.List, error) {
	log.C(ctx).Info("todoResolver.List")
	return r.todo.List(ctx, obj)
}

type listResolver struct {
	*RootResolver
}

func (l *listResolver) Owner(ctx context.Context, obj *graphql.List) (*graphql.User, error) {
	log.C(ctx).Info("listResolver.Owner")
	return l.list.Owner(ctx, obj)
}

func (l *listResolver) Todos(ctx context.Context, obj *graphql.List) ([]*graphql.Todo, error) {
	log.C(ctx).Info("listResolver.Todos")
	return l.list.Todos(ctx, obj)
}

func (l *listResolver) Collaborators(ctx context.Context, obj *graphql.List) ([]*graphql.ListAccess, error) {
	log.C(ctx).Info("listResolver.Collaborators")
	return l.list.Collaborators(ctx, obj)
}
