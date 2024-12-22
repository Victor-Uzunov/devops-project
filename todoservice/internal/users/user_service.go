package users

import (
	"context"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/models"
	"time"
)

//go:generate mockery --name=UserService --output=automock --with-expecter=true --outpkg=automock --case=underscore --disable-version-string --with-expecter=true
type UserService interface {
	CreateUser(ctx context.Context, list models.User) (string, error)
	GetUser(ctx context.Context, id string) (models.User, error)
	GetUserByEmail(ctx context.Context, username string) (models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	UpdateUser(ctx context.Context, todo models.User) error
	DeleteUser(ctx context.Context, id string) error
	SaveRefreshToken(ctx context.Context, email string, refreshToken string, expirationTime time.Time) error
	FindByRefreshToken(ctx context.Context, refreshToken string) (models.User, error)
	Logout(ctx context.Context, email string) error
}

//go:generate mockery --name=UUIDService --output=automock --with-expecter=true --outpkg=automock --case=underscore --disable-version-string
type UUIDService interface {
	Generate() string
}

//go:generate mockery --name=TimeService --output=automock --with-expecter=true --outpkg=automock --case=underscore --disable-version-string
type TimeService interface {
	Now() time.Time
}

var _ UserService = &service{}

type service struct {
	repo        UserRepository
	uuidService UUIDService
	timeService TimeService
}

func NewService(repo UserRepository, uuidService UUIDService, timeService TimeService) UserService {
	return &service{repo: repo, uuidService: uuidService, timeService: timeService}
}

func (s *service) CreateUser(ctx context.Context, user models.User) (string, error) {
	log.C(ctx).Infof("creating user: %v", user)
	if err := validateUser(user); err != nil {
		log.C(ctx).Errorf("validation error: %v", err)
		return "invalid data for the user", err
	}

	user.ID = s.uuidService.Generate()
	user.CreatedAt = s.timeService.Now()
	user.UpdatedAt = s.timeService.Now()

	log.C(ctx).Debugf("creating user with id: %v", user.ID)

	return s.repo.Create(ctx, user)
}

func (s *service) GetUser(ctx context.Context, id string) (models.User, error) {
	log.C(ctx).Infof("getting user: %v", id)
	return s.repo.Get(ctx, id)
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	log.C(ctx).Infof("deleting user: %v", id)
	return s.repo.Delete(ctx, id)
}

func (s *service) UpdateUser(ctx context.Context, user models.User) error {
	log.C(ctx).Infof("updating user: %v", user)
	if err := validateUser(user); err != nil {
		return err
	}
	dbUser, err := s.repo.Get(ctx, user.ID)
	if err != nil {
		log.C(ctx).Errorf("user not found: %v", user)
		return err
	}
	log.C(ctx).Debugf("updating user with id: %s", dbUser.ID)

	return s.repo.Update(ctx, user)
}

func (s *service) GetAllUsers(ctx context.Context) ([]models.User, error) {
	log.C(ctx).Infof("getting all users")
	return s.repo.GetAll(ctx)
}

func (s *service) SaveRefreshToken(ctx context.Context, email string, refreshToken string, expirationTime time.Time) error {
	log.C(ctx).Infof("service layer: saving refresh token: %v", refreshToken)
	return s.repo.SaveRefreshToken(ctx, email, refreshToken, expirationTime)
}

func validateUser(user models.User) error {
	return nil
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	log.C(ctx).Infof("getting user by email: %v", email)
	return s.repo.GetByEmail(ctx, email)
}

func (s *service) FindByRefreshToken(ctx context.Context, refreshToken string) (models.User, error) {
	log.C(ctx).Infof("getting user by refresh token: %v", refreshToken)
	return s.repo.FindByRefreshToken(ctx, refreshToken)
}

func (s *service) Logout(ctx context.Context, email string) error {
	log.C(ctx).Infof("logout: %v", email)
	return s.repo.Logout(ctx, email)
}
