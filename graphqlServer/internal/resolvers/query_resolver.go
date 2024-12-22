package resolvers

import (
	"context"
	"github.com/Victor-Uzunov/devops-project/graphqlServer/generated/graphql"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
)

type queryResolver struct {
	*RootResolver
}

func NewQueryResolver(resolver RootResolver) *queryResolver {
	return &queryResolver{
		RootResolver: &resolver,
	}
}

func (r *queryResolver) ListsGlobal(ctx context.Context) ([]*graphql.List, error) {
	log.C(ctx).Infof("queryResolve ListsGlobal")
	return r.list.ListsGlobal(ctx)
}

func (r *queryResolver) TodosGlobal(ctx context.Context) ([]*graphql.Todo, error) {
	log.C(ctx).Infof("queryResolve TodosGlobal")
	return r.todo.TodosGlobal(ctx)
}

func (r *queryResolver) Users(ctx context.Context) ([]*graphql.User, error) {
	log.C(ctx).Info("queryResolver users")
	return r.user.Users(ctx)
}

func (r *queryResolver) User(ctx context.Context, id string) (*graphql.User, error) {
	log.C(ctx).Info("queryResolver user", id)
	return r.user.User(ctx, id)
}

func (r *queryResolver) UserByEmail(ctx context.Context) (*graphql.User, error) {
	log.C(ctx).Info("queryResolver user by email")
	return r.user.UserByEmail(ctx)
}

func (r *queryResolver) Lists(ctx context.Context) ([]*graphql.List, error) {
	log.C(ctx).Info("queryResolver lists")
	return r.list.Lists(ctx)
}

func (r *queryResolver) List(ctx context.Context, id string) (*graphql.List, error) {
	log.C(ctx).Infof("queryResolver list with id: %s", id)
	return r.list.List(ctx, id)
}

func (r *queryResolver) Todos(ctx context.Context) ([]*graphql.Todo, error) {
	log.C(ctx).Info("queryResolver todos")
	return r.todo.Todos(ctx)
}

func (r *queryResolver) Todo(ctx context.Context, id string) (*graphql.Todo, error) {
	log.C(ctx).Infof("queryResolver todo with id %s", id)
	return r.todo.Todo(ctx, id)
}

func (r *queryResolver) TodosByList(ctx context.Context, id string) ([]*graphql.Todo, error) {
	log.C(ctx).Infof("queryResolver todos by list with id %s", id)
	return r.todo.TodosByList(ctx, id)
}

func (r *queryResolver) ListsPending(ctx context.Context) ([]*graphql.List, error) {
	log.C(ctx).Info("queryResolve ListsPending")
	return r.list.ListsPending(ctx)
}

func (r *queryResolver) UsersByList(ctx context.Context, id string) ([]*graphql.User, error) {
	log.C(ctx).Infof("queryResolve UsersByList with id %s", id)
	return r.user.UsersByList(ctx, id)
}

func (r *queryResolver) GetListAccesses(ctx context.Context, listID string) ([]*graphql.ListAccess, error) {
	log.C(ctx).Infof("queryResolve GetListAccesses for id %s", listID)
	return r.list.GetListAccesses(ctx, listID)
}

func (r *queryResolver) ListsAccepted(ctx context.Context) ([]*graphql.List, error) {
	log.C(ctx).Info("queryResolve ListsAccepted")
	return r.list.ListsAccepted(ctx)
}
