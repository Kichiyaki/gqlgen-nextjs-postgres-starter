package usecase

import (
	"context"
	"fmt"

	"github.com/kichiyaki/graphql-starter/backend/utils"

	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/kichiyaki/graphql-starter/backend/sessions"

	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"

	pgfilter "github.com/kichiyaki/pg-filter"

	"github.com/kichiyaki/graphql-starter/backend/middleware"

	"github.com/kichiyaki/graphql-starter/backend/auth"
	"github.com/kichiyaki/graphql-starter/backend/email"
	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/kichiyaki/graphql-starter/backend/token"
	"github.com/kichiyaki/graphql-starter/backend/user"
	"github.com/kichiyaki/graphql-starter/backend/user/validate"
)

const (
	limitOfActivationTokens    = 3
	limitOfResetPasswordTokens = 3
)

type Config struct {
	ApplicationName string
	FrontendURL     string
	UserRepo        user.Repository
	TokenRepo       token.Repository
	Email           email.Email
	SessStore       sessions.Store
}

type authUsecase struct {
	cfg *Config
}

func NewAuthUsecase(cfg *Config) auth.Usecase {
	return &authUsecase{
		cfg,
	}
}

func (ucase *authUsecase) Signup(ctx context.Context, input models.UserInput) (*models.User, error) {
	localizer, _ := middleware.LocalizerFromContext(ctx)
	if ucase.IsLogged(ctx) {
		utils.GetErrorMsg(localizer, "ErrCannotCreateAccountWhileLoggedIn")
		return nil, utils.GetErrorMsg(localizer, "ErrCannotCreateAccountWhileLoggedIn")
	}

	user := input.ToUser()
	user.Activated = false
	user.Role = models.DefaultRole

	cfg := validate.NewConfig(localizer)
	cfg.Login = true
	cfg.Password = true
	cfg.Role = true
	cfg.Email = true

	if err := cfg.Validate(*user); err != nil {
		return nil, err
	}

	err := ucase.cfg.UserRepo.Store(ctx, user)
	if err != nil {
		return nil, err
	}

	token := models.NewToken(models.AccountActivationTokenType, user.ID)
	if err := ucase.cfg.TokenRepo.Store(ctx, token); err != nil {
		utils.GetErrorMsg(localizer, "ErrActivationTokenCannotBeCreated")
		return nil, utils.GetErrorMsg(localizer, "ErrActivationTokenCannotBeCreated")
	}
	go func() {
		ucase.cfg.Email.Send(email.
			NewEmailConfig().
			SetTo([]string{user.Email}).
			SetBody(fmt.Sprintf("<p>Token: %s / URL: %s/%d/activate/%s</p>",
				token.Value,
				ucase.cfg.FrontendURL,
				user.ID,
				token.Value)).
			SetSubject(localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "SignupEmailSubject",
				TemplateData: map[string]interface{}{
					"ApplicationName": ucase.cfg.ApplicationName,
				},
			})))
	}()

	return user, nil
}

func (ucase *authUsecase) Login(ctx context.Context, login, password string) (*models.User, error) {
	localizer, _ := middleware.LocalizerFromContext(ctx)
	if ucase.IsLogged(ctx) {
		return nil, utils.GetErrorMsg(localizer, "ErrCannotLoginWhileLoggedIn")
	}

	return ucase.cfg.UserRepo.GetByCredentials(ctx, login, password)
}

func (ucase *authUsecase) Logout(ctx context.Context) error {
	localizer, _ := middleware.LocalizerFromContext(ctx)
	if !ucase.IsLogged(ctx) {
		return utils.GetErrorMsg(localizer, "ErrNotLoggedIn")
	}
	return nil
}

func (ucase *authUsecase) Activate(ctx context.Context, id int, token string) (*models.User, error) {
	localizer, _ := middleware.LocalizerFromContext(ctx)
	user, err := ucase.cfg.UserRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user.Activated {
		return nil, utils.GetErrorMsg(localizer, "ErrAccountHasBeenActivated")
	}

	filter := &models.TokenFilter{
		Type:   models.AccountActivationTokenType,
		Value:  token,
		UserID: fmt.Sprint(id),
	}
	tokens, err := ucase.cfg.TokenRepo.Fetch(ctx, pgfilter.New(filter.ToMap()))
	if err != nil {
		return nil, err
	}
	if len(tokens) == 0 {
		return nil, utils.GetErrorMsg(localizer, "ErrInvalidActivationToken")
	}

	user.Activated = true
	if err := ucase.cfg.UserRepo.Update(ctx, user); err != nil {
		return nil, utils.GetErrorMsg(localizer, "ErrAccountCannotBeActivated")
	}

	go func() {
		filter := &models.TokenFilter{
			Type:   models.AccountActivationTokenType,
			UserID: fmt.Sprint(id),
		}
		ctx := context.Background()
		tokens, err := ucase.cfg.TokenRepo.Fetch(ctx, pgfilter.New(filter.ToMap()))
		if len(tokens) > 0 && err == nil {
			ids := []int{}
			for _, token := range tokens {
				ids = append(ids, token.ID)
			}
			ucase.cfg.TokenRepo.Delete(ctx, ids)
		}
	}()

	return user, nil
}

func (ucase *authUsecase) GenerateNewActivationToken(ctx context.Context, id int) error {
	localizer, _ := middleware.LocalizerFromContext(ctx)
	user, err := ucase.cfg.UserRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user.Activated {
		return utils.GetErrorMsg(localizer, "ErrAccountHasBeenActivated")
	}
	filter := &models.TokenFilter{
		Type:   models.AccountActivationTokenType,
		UserID: fmt.Sprint(id),
	}
	tokens, err := ucase.cfg.TokenRepo.Fetch(ctx, pgfilter.New(filter.ToMap()))
	if err != nil {
		return err
	}
	if len(tokens) > limitOfActivationTokens {
		utils.GetErrorMsg(localizer, "ErrReachedLimitOfActivationTokens")
		return utils.GetErrorMsg(localizer, "ErrReachedLimitOfActivationTokens")
	}

	token := models.NewToken(models.AccountActivationTokenType, user.ID)
	if err := ucase.cfg.TokenRepo.Store(ctx, token); err != nil {
		return utils.GetErrorMsg(localizer, "ErrActivationTokenCannotBeCreated")
	}
	go func() {
		ucase.cfg.Email.Send(email.
			NewEmailConfig().
			SetTo([]string{user.Email}).
			SetBody(fmt.Sprintf("<p>Token: %s</p>, URL: %s/%d/reset-password/%s",
				token.Value,
				ucase.cfg.FrontendURL,
				user.ID,
				token.Value)).
			SetSubject(localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "NewActivationTokenEmailSubject",
			})))
	}()

	return nil
}

func (ucase *authUsecase) GenerateNewActivationTokenForCurrentUser(ctx context.Context) error {
	localizer, _ := middleware.LocalizerFromContext(ctx)
	if !ucase.IsLogged(ctx) {
		return utils.GetErrorMsg(localizer, "ErrNotLoggedIn")
	}
	return ucase.GenerateNewActivationToken(ctx, ucase.CurrentUser(ctx).ID)
}

func (ucase *authUsecase) GenerateNewResetPasswordToken(ctx context.Context, emailAddress string) error {
	localizer, _ := middleware.LocalizerFromContext(ctx)
	user, err := ucase.cfg.UserRepo.GetByEmail(ctx, emailAddress)
	if err != nil {
		return err
	}
	filter := &models.TokenFilter{
		Type:   models.ResetPasswordTokenType,
		UserID: fmt.Sprint(user.ID),
	}
	tokens, err := ucase.cfg.TokenRepo.Fetch(ctx, pgfilter.New(filter.ToMap()))
	if err != nil {
		return err
	}
	if len(tokens) > limitOfResetPasswordTokens {
		return utils.GetErrorMsg(localizer, "ErrReachedLimitOfResetPasswordTokens")

	}

	token := models.NewToken(models.ResetPasswordTokenType, user.ID)
	if err := ucase.cfg.TokenRepo.Store(ctx, token); err != nil {
		return utils.GetErrorMsg(localizer, "ErrResetPasswordTokenCannotBeCreated")

	}
	go func() {
		ucase.cfg.Email.Send(email.
			NewEmailConfig().
			SetTo([]string{user.Email}).
			SetBody(fmt.Sprintf("<p>Token: %s, URL: %s/%d/reset-password/%s",
				token.Value,
				ucase.cfg.FrontendURL,
				user.ID,
				token.Value)).
			SetSubject(localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "ResetPasswordEmailSubject",
			})))
	}()

	return nil
}

func (ucase *authUsecase) ResetPassword(ctx context.Context, id int, token string) error {
	localizer, _ := middleware.LocalizerFromContext(ctx)
	user, err := ucase.cfg.UserRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	filter := &models.TokenFilter{
		Type:   models.ResetPasswordTokenType,
		Value:  token,
		UserID: fmt.Sprint(id),
	}
	tokens, err := ucase.cfg.TokenRepo.Fetch(ctx, pgfilter.New(filter.ToMap()))
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		return utils.GetErrorMsg(localizer, "ErrInvalidResetPasswordToken")

	}
	password, err := password.Generate(32, 6, 6, false, true)
	if err != nil {
		return utils.GetErrorMsg(localizer, "ErrCannotGeneratePassword")

	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return utils.GetErrorMsg(localizer, "ErrCannotGeneratePassword")

	}
	user.Password = string(hashedPassword)
	if err := ucase.cfg.UserRepo.Update(ctx, user); err != nil {
		return utils.GetErrorMsg(localizer, "ErrUserCannotBeUpdated")

	}
	go func() {
		ucase.cfg.TokenRepo.Delete(ctx, []int{tokens[0].ID})
	}()
	go func() {
		ucase.cfg.Email.Send(email.
			NewEmailConfig().
			SetTo([]string{user.Email}).
			SetBody(fmt.Sprintf("<p>Nowe hasło: %s</p>", password)).
			SetSubject(localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "NewPasswordEmailSubject",
			})))
	}()
	go func() {
		sess, _ := ucase.cfg.SessStore.GetAll()
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
		ucase.cfg.SessStore.DeleteByID(ids...)
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
