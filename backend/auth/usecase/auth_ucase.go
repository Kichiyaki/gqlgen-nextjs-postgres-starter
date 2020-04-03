package usecase

import (
	"backend/auth"
	_errors "backend/errors"
	"backend/models"
	"backend/user"
	"backend/user/validation"
	"context"

	"github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
)

type Config struct {
	UserRepo          user.Repository
	PasswordGenerator password.PasswordGenerator
}

type usecase struct {
	userRepo  user.Repository
	generator password.PasswordGenerator
	logrus    *logrus.Entry
}

func NewAuthUsecase(cfg Config) auth.Usecase {
	if cfg.PasswordGenerator == nil {
		cfg.PasswordGenerator, _ = password.NewGenerator(nil)
	}
	return &usecase{
		cfg.UserRepo,
		cfg.PasswordGenerator,
		logrus.WithField("package", "auth/usecase"),
	}
}

func (ucase *usecase) Signup(ctx context.Context, input models.UserInput) (*models.User, error) {
	u := input.ToUser()
	log := ucase.logrus.WithField("user", u)
	log.Debug("Signup")
	u.Role = models.UserDefaultRole
	u.Activated = false
	u.ActivationToken = uuid.New().String()
	if u.DisplayName == "" {
		u.DisplayName = u.Login
	}
	cfg := validation.NewConfig()
	if err := cfg.Validate(u); err != nil {
		log.Debugf("Cannot create user: %s", err.Error())
		return nil, err
	}
	if err := ucase.userRepo.Store(ctx, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (ucase *usecase) Signin(ctx context.Context, login, password string) (*models.User, error) {
	ucase.logrus.WithField("password", password).WithField("login", login).Debug("Sign in")
	return ucase.userRepo.GetByCredentials(ctx, login, password)
}

func (ucase *usecase) GenerateNewActivationToken(ctx context.Context, id int) (*models.User, error) {
	ucase.logrus.WithField("id", id).Debug("GenerateNewActivationToken")
	u, err := ucase.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	} else if u.Activated {
		ucase.logrus.WithField("id", id).Debug("The account is activated.")
		return nil, _errors.Wrap(_errors.ErrUnauthorized)
	}
	u.ActivationToken = uuid.New().String()
	return u, ucase.userRepo.Update(ctx, u)
}

func (ucase *usecase) Activate(ctx context.Context, id int, token string) (*models.User, error) {
	log := ucase.logrus.WithField("id", id).WithField("token", token)
	log.Debug("Activate")
	u, err := ucase.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	} else if u.Activated {
		log.Debug("The account is activated.")
		return nil, _errors.Wrap(_errors.ErrUnauthorized)
	}
	if u.ActivationToken == token {
		u.Activated = true
		if err := ucase.userRepo.Update(ctx, u); err != nil {
			return nil, err
		}
		return u, nil
	}
	log.Debug("Wrong activation token.")
	return nil, _errors.Wrap(_errors.ErrWrongActivationToken)
}

func (ucase *usecase) GenerateNewResetPasswordToken(ctx context.Context, email string) (*models.User, error) {
	ucase.logrus.WithField("email", email).Debug("GenerateNewResetPasswordToken")
	u, err := ucase.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	u.ResetPasswordToken = uuid.New().String()
	if err := ucase.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (ucase *usecase) ResetPassword(ctx context.Context, id int, token string) (*models.User, string, error) {
	log := ucase.logrus.WithField("id", id).WithField("token", token)
	log.Debug("ResetPassword")
	u, err := ucase.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, "", err
	}
	if u.ResetPasswordToken == token {
		pswd := ucase.generator.MustGenerate(16, 4, 4, false, false)
		u.Password = pswd
		u.ResetPasswordToken = uuid.New().String()
		if err := ucase.userRepo.Update(ctx, u); err != nil {
			return nil, "", err
		}
		return u, pswd, nil
	}
	log.Debug("Wrong reset password token")
	return u, "", _errors.Wrap(_errors.ErrWrongResetPasswordToken)
}
