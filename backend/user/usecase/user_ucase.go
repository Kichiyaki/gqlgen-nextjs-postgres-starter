package usecase

import (
	"backend/models"
	"backend/user"
	"backend/user/validation"
	"context"

	"github.com/sirupsen/logrus"
)

type Config struct {
	UserRepo user.Repository
}

type usecase struct {
	userRepo user.Repository
	logrus   *logrus.Entry
}

func NewUserUsecase(cfg Config) user.Usecase {
	return &usecase{
		cfg.UserRepo,
		logrus.WithField("package", "user/usecase"),
	}
}

func (ucase *usecase) Fetch(ctx context.Context, f *models.UserFilter) (models.UserList, error) {
	ucase.logrus.WithField("filter", f).Debug("Fetch")
	if f == nil {
		f = &models.UserFilter{
			Limit: 100,
		}
	}
	return ucase.userRepo.Fetch(ctx, f)
}

func (ucase *usecase) GetByID(ctx context.Context, id int) (*models.User, error) {
	ucase.logrus.WithField("id", id).Debug("GetByID")
	return ucase.userRepo.GetByID(ctx, id)
}

func (ucase *usecase) GetBySlug(ctx context.Context, slug string) (*models.User, error) {
	ucase.logrus.WithField("slug", slug).Debug("GetBySlug")
	return ucase.userRepo.GetBySlug(ctx, slug)
}

func (ucase *usecase) Update(ctx context.Context, id int, input models.UserInput) (*models.User, error) {
	entry := ucase.logrus.WithField("id", id).WithField("input", input)
	entry.Debug("Update")
	user := input.ToUser()
	user.ID = id
	cfg := validation.NewConfig()
	if user.Login == "" {
		cfg.Login = false
	}
	if user.Password == "" {
		cfg.Password = false
	}
	if user.Role == 0 {
		cfg.Role = false
	}
	if user.Email == "" {
		cfg.Email = false
	}
	if err := cfg.Validate(user); err != nil {
		entry.Debugf("Update - Validation error: %s", err.Error())
		return nil, err
	}
	if err := ucase.userRepo.Update(ctx, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ucase *usecase) Store(ctx context.Context, input models.UserInput) (*models.User, error) {
	entry := ucase.logrus.WithField("input", input)
	entry.Debug("Store")
	user := input.ToUser()
	if err := validation.NewConfig().Validate(user); err != nil {
		entry.Debugf("Store - Validation error: %s", err.Error())
		return nil, err
	}
	if err := ucase.userRepo.Update(ctx, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ucase *usecase) Delete(ctx context.Context, ids ...int) ([]*models.User, error) {
	entry := ucase.logrus.WithField("ids", ids)
	entry.Debug("Delete")
	f := &models.UserFilter{
		ID: ids,
	}
	users, err := ucase.userRepo.Delete(ctx, f)
	if err != nil {
		return nil, err
	}
	return users, nil
}
