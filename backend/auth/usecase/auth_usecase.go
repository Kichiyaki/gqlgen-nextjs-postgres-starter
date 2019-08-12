package usecase

import (
	"context"
	"fmt"

	"github.com/kichiyaki/graphql-starter/backend/sessions"

	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"

	pgfilter "github.com/kichiyaki/pg-filter"

	"github.com/kichiyaki/graphql-starter/backend/auth/errors"
	"github.com/kichiyaki/graphql-starter/backend/middleware"

	"github.com/kichiyaki/graphql-starter/backend/auth"
	"github.com/kichiyaki/graphql-starter/backend/email"
	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/kichiyaki/graphql-starter/backend/token"
	"github.com/kichiyaki/graphql-starter/backend/user"
	_userErrors "github.com/kichiyaki/graphql-starter/backend/user/errors"
	"github.com/kichiyaki/graphql-starter/backend/user/validate"
)

const (
	limitOfActivationTokens    = 3
	limitOfResetPasswordTokens = 3
)

type authUsecase struct {
	userRepo  user.Repository
	tokenRepo token.Repository
	email     email.Email
	sessStore sessions.Store
}

func NewAuthUsecase(userRepo user.Repository, tokenRepo token.Repository, e email.Email, sessStore sessions.Store) auth.Usecase {
	return &authUsecase{
		userRepo,
		tokenRepo,
		e,
		sessStore,
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

	filter := &models.TokenFilter{
		Type:   models.AccountActivationTokenType,
		Value:  token,
		UserID: fmt.Sprint(id),
	}
	tokens, err := ucase.tokenRepo.Fetch(ctx, pgfilter.New(filter.ToMap()))
	if err != nil {
		return nil, err
	}
	if len(tokens) == 0 {
		return nil, errors.ErrInvalidActivationToken
	}

	user.Activated = true
	if err := ucase.userRepo.Update(ctx, user); err != nil {
		return nil, errors.ErrAccountCannotBeActivated
	}

	go func() {
		filter := &models.TokenFilter{
			Type:   models.AccountActivationTokenType,
			UserID: fmt.Sprint(id),
		}
		ctx := context.Background()
		tokens, err := ucase.tokenRepo.Fetch(ctx, pgfilter.New(filter.ToMap()))
		if len(tokens) > 0 && err == nil {
			ids := []int{}
			for _, token := range tokens {
				ids = append(ids, token.ID)
			}
			ucase.tokenRepo.Delete(ctx, ids)
		}
	}()

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
	filter := &models.TokenFilter{
		Type:   models.AccountActivationTokenType,
		UserID: fmt.Sprint(id),
	}
	tokens, err := ucase.tokenRepo.Fetch(ctx, pgfilter.New(filter.ToMap()))
	if err != nil {
		return err
	}
	if len(tokens) > limitOfActivationTokens {
		return errors.ErrReachedLimitOfActivationTokens
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

func (ucase *authUsecase) GenerateNewActivationTokenForCurrentUser(ctx context.Context) error {
	if !ucase.IsLogged(ctx) {
		return errors.ErrNotLoggedIn
	}
	return ucase.GenerateNewActivationToken(ctx, ucase.CurrentUser(ctx).ID)
}

func (ucase *authUsecase) GenerateNewResetPasswordToken(ctx context.Context, emailAddress string) error {
	user, err := ucase.userRepo.GetByEmail(ctx, emailAddress)
	if err != nil {
		return err
	}
	filter := &models.TokenFilter{
		Type:   models.ResetPasswordTokenType,
		UserID: fmt.Sprint(user.ID),
	}
	tokens, err := ucase.tokenRepo.Fetch(ctx, pgfilter.New(filter.ToMap()))
	if err != nil {
		return err
	}
	if len(tokens) > limitOfResetPasswordTokens {
		return errors.ErrReachedLimitOfResetPasswordTokens
	}

	token := models.NewToken(models.ResetPasswordTokenType, user.ID)
	if err := ucase.tokenRepo.Store(ctx, token); err != nil {
		return errors.ErrResetPasswordTokenCannotBeCreated
	}
	go func() {
		ucase.email.Send(email.
			NewEmailConfig().
			SetTo([]string{user.Email}).
			SetBody(fmt.Sprintf("<p>Token: %s</p>", token.Value)).
			SetSubject("Token do zmiany hasła"))
	}()

	return nil
}

func (ucase *authUsecase) ResetPassword(ctx context.Context, id int, token string) error {
	user, err := ucase.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	filter := &models.TokenFilter{
		Type:   models.ResetPasswordTokenType,
		Value:  token,
		UserID: fmt.Sprint(id),
	}
	tokens, err := ucase.tokenRepo.Fetch(ctx, pgfilter.New(filter.ToMap()))
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		return errors.ErrInvalidResetPasswordToken
	}
	password, err := password.Generate(32, 6, 6, false, true)
	if err != nil {
		return errors.ErrCannotGeneratePassword
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.ErrCannotGeneratePassword
	}
	user.Password = string(hashedPassword)
	if err := ucase.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf(_userErrors.UserCannotBeUpdatedErrFormatWithLogin, user.Login)
	}
	go func() {
		ucase.tokenRepo.Delete(ctx, []int{tokens[0].ID})
	}()
	go func() {
		ucase.email.Send(email.
			NewEmailConfig().
			SetTo([]string{user.Email}).
			SetBody(fmt.Sprintf("<p>Nowe hasło: %s</p>", password)).
			SetSubject("Nowe hasło"))
	}()
	go func() {
		sess, _ := ucase.sessStore.GetAll()
		ids := []string{}
		for _, session := range sess {
			v := session.Values["user"]
			if v != nil {
				id, ok := v.(float64)
				userID := int(id)
				if ok {
					if userID == user.ID {
						ids = append(ids, session.ID)
					}
				}
			}
		}
		ucase.sessStore.DeleteByID(ids...)
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
