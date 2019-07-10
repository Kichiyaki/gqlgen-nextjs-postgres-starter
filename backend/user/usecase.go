package user

import (
	"context"

	"github.com/kichiyaki/graphql-starter/backend/models"
)

type Usecase interface {
	Store(ctx context.Context, input models.UserInput) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetBySlug(ctx context.Context, slug string) (*models.User, error)
	GetByCredentials(ctx context.Context, login, password string) (*models.User, error)
	Fetch(ctx context.Context,
		p models.Pagination,
		f *models.UserFilter) (*models.List, error)
	Update(ctx context.Context, id int, input models.UserInput) (*models.User, error)
	Delete(ctx context.Context, ids []int) ([]*models.User, error)
}
