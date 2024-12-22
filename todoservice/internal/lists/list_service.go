package lists

import (
	"context"
	"errors"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"time"
)

//go:generate mockery --name=ListService --output=automock --with-expecter=true --outpkg=automock --case=underscore --disable-version-string --with-expecter=true
type ListService interface {
	CreateList(ctx context.Context, list models.List) (string, error)
	CreateAccess(ctx context.Context, list models.Access) (models.Access, error)
	GetList(ctx context.Context, id string) (models.List, error)
	GetAccess(ctx context.Context, listID string, userID string) (models.Access, error)
	GetAllLists(ctx context.Context) ([]models.List, error)
	GetUsersByListID(ctx context.Context, listID string) ([]models.Access, error)
	GetListOwnerID(ctx context.Context, listID string) (string, error)
	UpdateList(ctx context.Context, list models.List) error
	DeleteList(ctx context.Context, id string) error
	DeleteAccess(ctx context.Context, listID string, userID string) error
	ListAllByUserID(ctx context.Context, useID string) ([]models.Access, error)
	UpdateListDescription(ctx context.Context, id, description string) (models.List, error)
	UpdateListName(ctx context.Context, id string, name string) (models.List, error)
	GetAllTodosForList(ctx context.Context, listID string) ([]models.Todo, error)
	GetPendingLists(ctx context.Context, userID string) ([]models.Access, error)
	AcceptList(ctx context.Context, listID string, userID string) error
	GetAccessesByListID(ctx context.Context, listID string) ([]models.Access, error)
	GetAcceptedLists(ctx context.Context, userID string) ([]models.Access, error)
}

//go:generate mockery --name=UUIDService --output=automock --with-expecter=true --outpkg=automock --case=underscore --disable-version-string
type UUIDService interface {
	Generate() string
}

//go:generate mockery --name=TimeService --output=automock --with-expecter=true --outpkg=automock --case=underscore --disable-version-string
type TimeService interface {
	Now() time.Time
}

var _ ListService = &service{}

type service struct {
	repo        ListRepository
	uuidService UUIDService
	timeService TimeService
}

func NewService(repo ListRepository, uuidService UUIDService, timeService TimeService) ListService {
	return &service{repo: repo, uuidService: uuidService, timeService: timeService}
}

func (s *service) CreateList(ctx context.Context, list models.List) (string, error) {
	log.C(ctx).Info("creating list service")
	if err := validateList(ctx, list); err != nil {
		return "", err
	}

	list.ID = s.uuidService.Generate()
	list.CreatedAt = s.timeService.Now()
	list.UpdatedAt = s.timeService.Now()

	return s.repo.Create(ctx, list)
}

func (s *service) GetList(ctx context.Context, id string) (models.List, error) {
	log.C(ctx).Info("getting list service")
	return s.repo.Get(ctx, id)
}

func (s *service) DeleteList(ctx context.Context, id string) error {
	log.C(ctx).Info("deleting list service")
	return s.repo.Delete(ctx, id)
}

func (s *service) DeleteAccess(ctx context.Context, listId string, userID string) error {
	log.C(ctx).Info("deleting list access service")
	return s.repo.DeleteAccess(ctx, listId, userID)
}

func (s *service) UpdateList(ctx context.Context, list models.List) error {
	log.C(ctx).Info("updating list service")
	if err := validateList(ctx, list); err != nil {
		return err
	}
	_, err := s.repo.Get(ctx, list.ID)
	log.C(ctx).Debugf("updating list service with id %s", list.ID)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, list)
}

func (s *service) ListAllByUserID(ctx context.Context, userID string) ([]models.Access, error) {
	log.C(ctx).Info("listing all access service")
	return s.repo.ListAllByUserID(ctx, userID)
}

func (s *service) GetAcceptedLists(ctx context.Context, userID string) ([]models.Access, error) {
	log.C(ctx).Info("get accepted lists service")
	return s.repo.GetAcceptedLists(ctx, userID)
}

func (s *service) GetPendingLists(ctx context.Context, userID string) ([]models.Access, error) {
	log.C(ctx).Info("getting pending list service")
	return s.repo.GetPendingLists(ctx, userID)
}

func (s *service) GetAllLists(ctx context.Context) ([]models.List, error) {
	log.C(ctx).Info("getting all list service")
	return s.repo.GetAll(ctx)
}

func (s *service) GetUsersByListID(ctx context.Context, listID string) ([]models.Access, error) {
	log.C(ctx).Info("getting all access service")
	return s.repo.GetUsersByListID(ctx, listID)
}

func (s *service) GetListOwnerID(ctx context.Context, listID string) (string, error) {
	log.C(ctx).Info("getting list owner service")
	return s.repo.GetListOwnerID(ctx, listID)
}

func (s *service) GetAccess(ctx context.Context, listID, userID string) (models.Access, error) {
	log.C(ctx).Info("getting access service")
	return s.repo.GetAccess(ctx, listID, userID)
}

func (s *service) CreateAccess(ctx context.Context, access models.Access) (models.Access, error) {
	log.C(ctx).Info("creating access service")
	return s.repo.CreateAccess(ctx, access)
}

func (s *service) UpdateListDescription(ctx context.Context, id, description string) (models.List, error) {
	log.C(ctx).Info("updating list description service")
	return s.repo.UpdateListDescription(ctx, id, description)
}

func (s *service) UpdateListName(ctx context.Context, id, name string) (models.List, error) {
	log.C(ctx).Info("updating list name service")
	return s.repo.UpdateListName(ctx, id, name)
}

func (s *service) GetAllTodosForList(ctx context.Context, listID string) ([]models.Todo, error) {
	log.C(ctx).Info("listing todos by list id")
	return s.repo.GetAllTodosForList(ctx, listID)
}

func (s *service) AcceptList(ctx context.Context, listID string, userID string) error {
	log.C(ctx).Info("accepting list service")
	return s.repo.AcceptList(ctx, listID, userID)
}

func validateList(ctx context.Context, list models.List) error {
	log.C(ctx).Info("validating list service")
	if len(list.SharedWith) != 0 && list.Visibility == constants.VisibilityPrivate {
		log.C(ctx).Error("cannot have private visibility and be shared at the same time")
		return errors.New("cannot have private visibility and be shared at the same time")
	}
	return nil
}

func (s *service) GetAccessesByListID(ctx context.Context, listID string) ([]models.Access, error) {
	log.C(ctx).Infof("getting accesses by list id: %s", listID)
	return s.repo.GetAccessesByListID(ctx, listID)
}
