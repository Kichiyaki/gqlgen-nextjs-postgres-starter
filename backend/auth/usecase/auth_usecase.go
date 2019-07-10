package usecase

import (
	"context"
	"fmt"
	"strings"

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
		if strings.Contains(err.Error(), "duplicate key value violates unique") {
			if strings.Contains(err.Error(), "login") {
				return nil, fmt.Errorf("Podany login jest zajęty")
			} else if strings.Contains(err.Error(), "email") {
				return nil, fmt.Errorf("Podany email jest zajęty")
			}
		}
		return nil, fmt.Errorf("Wystąpił błąd podczas tworzenia użytkownika")
	}

	token := models.NewToken(models.AccountActivationTokenType, user.ID)
	if err := ucase.tokenRepo.Store(ctx, token); err != nil {
		return nil, fmt.Errorf("Nie udało się utworzyć tokenu aktywacyjnego")
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
	return ucase.userRepo.GetByCredentials(ctx, login, password)
}

func (ucase *authUsecase) Activate(ctx context.Context, id int, token string) (*models.User, error) {
	user, err := ucase.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user.Activated {
		return nil, fmt.Errorf("Konto zostało już aktywowane")
	}

	t, err := ucase.tokenRepo.Get(ctx, models.AccountActivationTokenType, token)
	if err != nil {
		return nil, err
	}
	if t.UserID != id {
		return nil, fmt.Errorf("Niepoprawny token aktywacyjny")
	}

	user.Activated = true
	if err := ucase.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("Wystąpił błąd podczas aktywacji konta")
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
		return fmt.Errorf("Konto zostało już aktywowane")
	}

	token := models.NewToken(models.AccountActivationTokenType, user.ID)
	if err := ucase.tokenRepo.Store(ctx, token); err != nil {
		return fmt.Errorf("Nie udało się utworzyć tokenu aktywacyjnego")
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
