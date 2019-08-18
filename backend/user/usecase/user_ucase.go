package usecase

import (
	"context"
	"fmt"

	"github.com/kichiyaki/graphql-starter/backend/utils"

	"github.com/kichiyaki/graphql-starter/backend/middleware"

	"github.com/kichiyaki/graphql-starter/backend/auth"

	pgfilter "github.com/kichiyaki/pg-filter"
	pgpagination "github.com/kichiyaki/pg-pagination"
	"github.com/kichiyaki/structs"
	"golang.org/x/crypto/bcrypt"

	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/kichiyaki/graphql-starter/backend/user"
	"github.com/kichiyaki/graphql-starter/backend/user/validate"
)

type userUsecase struct {
	userRepo  user.Repository
	authUcase auth.Usecase
}

func NewUserUsecase(userRepo user.Repository, authUcase auth.Usecase) user.Usecase {
	return &userUsecase{
		userRepo,
		authUcase,
	}
}

func (ucase *userUsecase) Store(ctx context.Context, input models.UserInput) (*models.User, error) {
	localizer, _ := middleware.LocalizerFromContext(ctx)
	if !ucase.authUcase.IsLogged(ctx) {
		return nil, utils.GetErrorMsg(localizer, "ErrNotLoggedIn")
	}
	if !ucase.authUcase.HasAdministrativePrivileges(ctx) {
		return nil, utils.GetErrorMsg(localizer, "ErrUnauthorized")
	}

	user := input.ToUser()

	cfg := validate.NewConfig(localizer)
	cfg.Login = true
	cfg.Password = true
	cfg.Role = true
	cfg.Email = true
	if err := cfg.Validate(*user); err != nil {
		return nil, err
	}

	err := ucase.userRepo.Store(ctx, user)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ucase *userUsecase) GetByID(ctx context.Context, id int) (*models.User, error) {
	return ucase.userRepo.GetByID(ctx, id)
}

func (ucase *userUsecase) GetBySlug(ctx context.Context, slug string) (*models.User, error) {
	return ucase.userRepo.GetBySlug(ctx, slug)
}

func (ucase *userUsecase) GetByCredentials(ctx context.Context, login, password string) (*models.User, error) {
	return ucase.userRepo.GetByCredentials(ctx, login, password)
}

func (ucase *userUsecase) Fetch(ctx context.Context,
	p models.Pagination,
	f *models.UserFilter) (*models.List, error) {
	m := make(map[string]string)
	if f != nil {
		fi := structs.Map(f)
		for key, value := range fi {
			strKey := fmt.Sprintf("%v", key)
			strValue := fmt.Sprintf("%v", value)

			m[strKey] = strValue
		}

		if f.OnlyActivated {
			m["activated"] = "true"
		}
	}

	return ucase.userRepo.Fetch(ctx,
		pgpagination.Pagination{Page: p.Page, Limit: p.Limit},
		pgfilter.New(m))
}

func (ucase *userUsecase) Update(ctx context.Context, id int, input models.UserInput) (*models.User, error) {
	localizer, _ := middleware.LocalizerFromContext(ctx)
	if !ucase.authUcase.IsLogged(ctx) {
		return nil, utils.GetErrorMsg(localizer, "ErrNotLoggedIn")
	}
	if !ucase.authUcase.HasAdministrativePrivileges(ctx) {
		return nil, utils.GetErrorMsg(localizer, "ErrUnauthorized")
	}

	user, err := ucase.GetByID(ctx, id)
	if err != nil || user == nil || user.ID <= 0 {
		return nil, err
	}

	cfg := validate.UserValidationConfig{}
	if input.Login != nil && *input.Login != "" {
		user.Login = *input.Login
		cfg.Login = true
	}
	if input.Password != nil && *input.Password != "" {
		user.Password = *input.Password
		cfg.Password = true
	}
	if input.Email != nil && *input.Email != "" {
		user.Email = *input.Email
		cfg.Email = true
	}
	if input.Role != nil && *input.Role != "" {
		user.Role = *input.Role
		cfg.Role = true
	}
	if input.Activated != nil && *input.Activated != user.Activated {
		user.Activated = *input.Activated
	}
	if somethingToValidate := cfg.IsSomethingToValidate(); !somethingToValidate {
		return nil, utils.GetErrorMsg(localizer, "ErrNothingChanged")
	}
	if err := cfg.Validate(*user); err != nil {
		return nil, err
	}

	if cfg.Password {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}

	if err := ucase.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (ucase *userUsecase) Delete(ctx context.Context, ids []int) ([]*models.User, error) {
	localizer, _ := middleware.LocalizerFromContext(ctx)
	if !ucase.authUcase.IsLogged(ctx) {
		return nil, utils.GetErrorMsg(localizer, "ErrNotLoggedIn")
	}
	if !ucase.authUcase.HasAdministrativePrivileges(ctx) {
		return nil, utils.GetErrorMsg(localizer, "ErrUnauthorized")
	}

	user := ucase.authUcase.CurrentUser(ctx)
	fmt.Println(user.ID, ids)
	for _, id := range ids {
		if id == user.ID {
			return nil, utils.GetErrorMsg(localizer, "ErrUserCannotDeleteHisAccountByHimself")
		}
	}

	return ucase.userRepo.Delete(ctx, ids)
}
