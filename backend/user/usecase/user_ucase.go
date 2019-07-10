package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/kichiyaki/graphql-starter/backend/middleware"

	pgfilter "github.com/kichiyaki/pg-filter"
	pgpagination "github.com/kichiyaki/pg-pagination"
	"github.com/kichiyaki/structs"
	"golang.org/x/crypto/bcrypt"

	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/kichiyaki/graphql-starter/backend/user"
	"github.com/kichiyaki/graphql-starter/backend/user/validate"
)

type userUsecase struct {
	userRepo user.Repository
}

func NewUserUsecase(userRepo user.Repository) user.Usecase {
	return &userUsecase{
		userRepo,
	}
}

func (ucase *userUsecase) Store(ctx context.Context, input models.UserInput) (*models.User, error) {
	user := input.ToUser()

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

	list, err := ucase.userRepo.Fetch(ctx,
		pgpagination.Pagination{Page: p.Page, Limit: p.Limit},
		pgfilter.New(m))
	if err != nil {
		return nil, fmt.Errorf("Nie udało się wygenerować listy użytkowników")
	}
	return list, nil
}

func (ucase *userUsecase) Update(ctx context.Context, id int, input models.UserInput) (*models.User, error) {
	user, err := ucase.GetByID(ctx, id)
	if err != nil || user == nil || user.ID <= 0 {
		return nil, fmt.Errorf("Nie znaleziono użytkownika o ID: %d", id)
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
		return nil, fmt.Errorf("Nie wprowadziłeś żadnych zmian w konfiguracji użytkownika")
	}
	if err := cfg.Validate(*user); err != nil {
		return nil, err
	}

	if cfg.Password {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}

	if err := ucase.userRepo.Update(ctx, user); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique") {
			if strings.Contains(err.Error(), "login") {
				return nil, fmt.Errorf("Podany login jest zajęty")
			} else if strings.Contains(err.Error(), "email") {
				return nil, fmt.Errorf("Podany email jest zajęty")
			}
		}
		return nil, fmt.Errorf("Nie udało się zaaktualizować użytkownika")
	}
	return user, nil
}

func (ucase *userUsecase) Delete(ctx context.Context, ids []int) ([]*models.User, error) {
	user, _ := middleware.UserFromContext(ctx)
	for _, id := range ids {
		if id == user.ID {
			return nil, fmt.Errorf("Nie możesz sam usunąć swojego konta")
		}
	}

	users, err := ucase.userRepo.Delete(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("Nie udało się usunąć kont użytkowników")
	}
	return users, nil
}
