package todos

import (
	"context"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"time"
)

//go:generate mockery --name=TodoService --output=automock --with-expecter=true --outpkg=automock --case=underscore --disable-version-string
type TodoService interface {
	CreateTodo(ctx context.Context, todo models.Todo) (string, error)
	GetTodo(ctx context.Context, id string) (models.Todo, error)
	GetAllTodos(ctx context.Context) ([]models.Todo, error)
	UpdateTodo(ctx context.Context, todo models.Todo) error
	DeleteTodo(ctx context.Context, id string) error
	ListTodosByListID(ctx context.Context, listID string) ([]models.Todo, error)
	CompleteTodo(ctx context.Context, id string) (models.Todo, error)
	UpdateTodoTitle(ctx context.Context, id, name string) (models.Todo, error)
	UpdateTodoDescription(ctx context.Context, id, description string) (models.Todo, error)
	UpdateTodoPriority(ctx context.Context, id string, priority constants.PriorityLevel) (models.Todo, error)
	UpdateAssignedTo(ctx context.Context, id, userID string) (models.Todo, error)
}

//go:generate mockery --name=UUIDService --output=automock --with-expecter=true --outpkg=automock --case=underscore --disable-version-string
type UUIDService interface {
	Generate() string
}

//go:generate mockery --name=TimeService --output=automock --with-expecter=true --outpkg=automock --case=underscore --disable-version-string
type TimeService interface {
	Now() time.Time
}

var _ TodoService = &service{}

type service struct {
	repo        TodoRepository
	uuidService UUIDService
	timeService TimeService
}

func NewService(repo TodoRepository, uuidService UUIDService, timeService TimeService) TodoService {
	return &service{repo: repo, uuidService: uuidService, timeService: timeService}
}

func (s *service) CreateTodo(ctx context.Context, todo models.Todo) (string, error) {
	log.C(ctx).Info("creating todo service")
	if err := validateTodo(todo); err != nil {
		return "", err
	}

	todo.ID = s.uuidService.Generate()
	todo.CreatedAt = s.timeService.Now()
	todo.UpdatedAt = s.timeService.Now()

	log.C(ctx).Debugf("creating todo with id %s", todo.ID)

	return s.repo.Create(ctx, todo)
}

func (s *service) GetTodo(ctx context.Context, id string) (models.Todo, error) {
	log.C(ctx).Info("getting todo service")
	return s.repo.Get(ctx, id)
}

func (s *service) UpdateTodo(ctx context.Context, todo models.Todo) error {
	log.C(ctx).Info("updating todo service")
	if err := validateTodo(todo); err != nil {
		return err
	}
	dbTodo, err := s.repo.Get(ctx, todo.ID)
	log.C(ctx).Debugf("updating todo in the database with id %s", dbTodo.ID)
	if err != nil {
		log.C(ctx).Errorf("getting todo with id %s failed", todo.ID)
		return err
	}
	todo.UpdatedAt = s.timeService.Now()
	return s.repo.Update(ctx, todo)
}

func (s *service) DeleteTodo(ctx context.Context, id string) error {
	log.C(ctx).Info("deleting todo service")
	return s.repo.Delete(ctx, id)
}

func (s *service) ListTodosByListID(ctx context.Context, listID string) ([]models.Todo, error) {
	log.C(ctx).Info("listing todos by list id")
	return s.repo.GetAllByListID(ctx, listID)
}

func (s *service) GetAllTodos(ctx context.Context) ([]models.Todo, error) {
	log.C(ctx).Info("getting all todos service")
	return s.repo.GetAll(ctx)
}

func (s *service) CompleteTodo(ctx context.Context, id string) (models.Todo, error) {
	log.C(ctx).Info("completing todo service")
	return s.repo.CompleteTodo(ctx, id)
}

func (s *service) UpdateTodoTitle(ctx context.Context, id, title string) (models.Todo, error) {
	log.C(ctx).Info("updating todo title service")
	return s.repo.UpdateTodoTitle(ctx, id, title)
}

func (s *service) UpdateTodoDescription(ctx context.Context, id, description string) (models.Todo, error) {
	log.C(ctx).Info("updating todo description service")
	return s.repo.UpdateTodoDescription(ctx, id, description)
}

func (s *service) UpdateTodoPriority(ctx context.Context, id string, priority constants.PriorityLevel) (models.Todo, error) {
	log.C(ctx).Info("updating todo priority service")
	return s.repo.UpdateTodoPriority(ctx, id, priority)
}

func (s *service) UpdateAssignedTo(ctx context.Context, id, userID string) (models.Todo, error) {
	log.C(ctx).Info("updating todo assigned service")
	return s.repo.UpdateAssignedTo(ctx, id, userID)
}

func validateTodo(todo models.Todo) error {
	return nil
}
