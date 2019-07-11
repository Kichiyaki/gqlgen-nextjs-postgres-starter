package usecase

import (
	"context"
	"fmt"

	"github.com/kichiyaki/graphql-starter/backend/auth/errors"
	"github.com/kichiyaki/graphql-starter/backend/middleware"

	"github.com/kichiyaki/graphql-starter/backend/auth"
	"github.com/kichiyaki/graphql-starter/backend/email"
	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/kichiyaki/graphql-starter/backend/token"
	"github.com/kichiyaki/graphql-starter/backend/user"
	"github.com/kichiyaki/graphql-starter/backend/user/validate"
)

type authUsecase struct {
	userRepo  user.Repository
	tokenRepo token.Repository
	email     email.Email
}

func NewAuthUsecase(userRepo user.Repository, tokenRepo token.Repository, e email.Email) auth.Usecase {
	return &authUsecase{
		userRepo,
		tokenRepo,
		e,
	}
}

func (ucase *authUsecase) Signup(ctx context.Context, input models.UserInput) (*models.User, error) {
	if ucase.IsLogged(ctx) {
		return nil, errors.ErrCannotCreateAccountWhileLoggedIn
	}

	user := input.ToUser()
	user.Activated = false
	user.Role = models.DefaultRole

	cfg := validate.UserValidationConfig{
		Login:    true,
		Password: true,
		Role:     true,
		Email:    true,
	}
	if err := cfg.Validate(*user); err != nil {
		return nil, err
	}

	err := ucase.userRepo.Store(ctx, user)
	if err != nil {
		return nil, err
	}

	token := models.NewToken(models.AccountActivationTokenType, user.ID)
	if err := ucase.tokenRepo.Store(ctx, token); err != nil {
		return nil, errors.ErrActivationTokenCannotBeCreated
	}
	go func() {
		ucase.email.Send(email.
			NewEmailConfig().
			SetTo([]string{user.Email}).
			SetBody(fmt.Sprintf("<p>Token: %s</p>", token.Value)).
			SetSubject("Rejestracja w serwisie xyz"))
	}()

	return user, nil
}

func (ucase *authUsecase) Login(ctx context.Context, login, password string) (*models.User, error) {
	if ucase.IsLogged(ctx) {
		return nil, errors.ErrCannotLoginWhileLoggedIn
	}

	return ucase.userRepo.GetByCredentials(ctx, login, password)
}

func (ucase *authUsecase) Logout(ctx context.Context) error {
	if !ucase.IsLogged(ctx) {
		return errors.ErrNotLoggedIn
	}
	return nil
}

func (ucase *authUsecase) Activate(ctx context.Context, id int, token string) (*models.User, error) {
	user, err := ucase.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user.Activated {
		return nil, errors.ErrAccountHasBeenActivated
	}

	t, err := ucase.tokenRepo.Get(ctx, models.AccountActivationTokenType, token)
	if err != nil {
		return nil, err
	}
	if t.UserID != id {
		return nil, errors.ErrInvalidActivationToken
	}

	user.Activated = true
	if err := ucase.userRepo.Update(ctx, user); err != nil {
		return nil, errors.ErrAccountCannotBeActivated
	}

	ucase.tokenRepo.DeleteByUserID(ctx, models.AccountActivationTokenType, user.ID)

	return user, nil
}

func (ucase *authUsecase) GenerateNewActivationToken(ctx context.Context, id int) error {
	user, err := ucase.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user.Activated {
		return errors.ErrAccountHasBeenActivated
	}

	token := models.NewToken(models.AccountActivationTokenType, user.ID)
	if err := ucase.tokenRepo.Store(ctx, token); err != nil {
		return errors.ErrActivationTokenCannotBeCreated
	}
	go func() {
		ucase.email.Send(email.
			NewEmailConfig().
			SetTo([]string{user.Email}).
			SetBody(fmt.Sprintf("<p>Token: %s</p>", token.Value)).
			SetSubject("Nowy token aktywacyjny"))
	}()

	return nil
}

func (ucase *authUsecase) IsLogged(ctx context.Context) bool {
	user, err := middleware.UserFromContext(ctx)
	return user != nil && user.ID > 0 && err == nil
}

func (ucase *authUsecase) HasAdministrativePrivileges(ctx context.Context) bool {
	user, err := middleware.UserFromContext(ctx)
	return user != nil && user.ID > 0 && err == nil && user.Role == models.AdministrativeRole
}

func (ucase *authUsecase) CurrentUser(ctx context.Context) *models.User {
	user, _ := middleware.UserFromContext(ctx)
	return user
}
