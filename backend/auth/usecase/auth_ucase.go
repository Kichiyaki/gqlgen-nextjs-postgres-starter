package usecase

import (
	"backend/auth"
	_errors "backend/errors"
	"backend/models"
	"backend/user"
	"backend/user/validation"
	"backend/utils"
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
)

type Config struct {
	UserRepo                        user.Repository
	PasswordGenerator               password.PasswordGenerator
	IntervalBetweenTokensGeneration int
	ResetPasswordTokenExpiresIn     int
}

type usecase struct {
	userRepo                        user.Repository
	generator                       password.PasswordGenerator
	logrus                          *logrus.Entry
	intervalBetweenTokensGeneration int
	resetPasswordTokenExpiresIn     int
}

func NewAuthUsecase(cfg Config) auth.Usecase {
	if cfg.PasswordGenerator == nil {
		cfg.PasswordGenerator, _ = password.NewGenerator(nil)
	}
	return &usecase{
		cfg.UserRepo,
		cfg.PasswordGenerator,
		logrus.WithField("package", "auth/usecase"),
		cfg.IntervalBetweenTokensGeneration,
		cfg.ResetPasswordTokenExpiresIn,
	}
}

func (ucase *usecase) Signup(ctx context.Context, input models.UserInput) (*models.User, error) {
	u := input.ToUser()
	entry := ucase.logrus.WithField("user", u)
	entry.Debug("Signup")
	u.Role = models.UserDefaultRole
	u.Activated = false
	u.ActivationToken = uuid.New().String()
	if u.DisplayName == "" {
		u.DisplayName = u.Login
	}
	cfg := validation.NewConfig()
	if err := cfg.Validate(u); err != nil {
		entry.Debugf("Signup - Cannot create user: %s", err.Error())
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
	now := time.Now()
	entry := ucase.logrus.WithField("id", id)
	entry.Debug("GenerateNewActivationToken")
	u, err := ucase.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	} else if u.Activated {
		entry.Debug("GenerateNewActivationToken - The account is activated.")
		return nil, _errors.Wrap(_errors.ErrAccountIsActivated)
	} else if isProperInterval(now, u.ActivationTokenGeneratedAt, ucase.intervalBetweenTokensGeneration) {
		entry.Debug("GenerateNewActivationToken - Token has been generated recently.")
		return nil, _errors.Wrap(_errors.ErrActivationTokenHasBeenGeneratedRecently)
	}
	u.ActivationToken = uuid.New().String()
	u.ActivationTokenGeneratedAt = now
	return u, ucase.userRepo.Update(ctx, u)
}

func (ucase *usecase) Activate(ctx context.Context, id int, token string) (*models.User, error) {
	entry := ucase.logrus.WithField("id", id).WithField("token", token)
	entry.Debug("Activate")
	u, err := ucase.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	} else if u.Activated {
		entry.Debug("The account is activated.")
		return nil, _errors.Wrap(_errors.ErrUnauthorized)
	}
	if u.ActivationToken == token {
		u.Activated = true
		if err := ucase.userRepo.Update(ctx, u); err != nil {
			return nil, err
		}
		return u, nil
	}
	entry.Debug("Activate - Wrong activation token.")
	return nil, _errors.Wrap(_errors.ErrWrongActivationToken)
}

func (ucase *usecase) GenerateNewResetPasswordToken(ctx context.Context, email string) (*models.User, error) {
	now := time.Now()
	entry := ucase.logrus.WithField("email", email)
	entry.Debug("GenerateNewResetPasswordToken")
	u, err := ucase.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	} else if isProperInterval(now, u.ResetPasswordTokenGeneratedAt, ucase.intervalBetweenTokensGeneration) {
		entry.Debug("GenerateNewResetPasswordToken - Token has been generated recently.")
		return nil, _errors.Wrap(_errors.ErrResetPasswordTokenHasBeenGeneratedRecently)
	}
	u.ResetPasswordToken = uuid.New().String()
	u.ResetPasswordTokenGeneratedAt = time.Now()
	if err := ucase.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (ucase *usecase) ResetPassword(ctx context.Context, id int, token string) (*models.User, string, error) {
	entry := ucase.logrus.WithField("id", id).WithField("token", token)
	entry.Debug("ResetPassword")
	u, err := ucase.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, "", err
	}
	if u.ResetPasswordToken == token {
		if !isProperInterval(time.Now(), u.ResetPasswordTokenGeneratedAt, ucase.resetPasswordTokenExpiresIn) {
			entry.Debug("ResetPassword - Reset password token expired.")
			return nil, "", _errors.Wrap(_errors.ErrTokenExpired)
		}
		pswd := ucase.generator.MustGenerate(16, 4, 4, false, false)
		u.Password = pswd
		u.ResetPasswordToken = uuid.New().String()
		if err := ucase.userRepo.Update(ctx, u); err != nil {
			return nil, "", err
		}
		return u, pswd, nil
	}
	entry.Debug("ResetPassword - Wrong reset password token")
	return u, "", _errors.Wrap(_errors.ErrWrongResetPasswordToken)
}

func isProperInterval(a, b time.Time, interval int) bool {
	year, month, day, hour, min, _ := utils.DateDifference(a, b)
	return year == 0 &&
		month == 0 &&
		day == 0 &&
		hour == 0 &&
		min < interval
}
