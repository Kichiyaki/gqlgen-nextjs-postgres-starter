package user

import (
	"context"

	"github.com/kichiyaki/graphql-starter/backend/models"
	pgfilter "github.com/kichiyaki/pg-filter"
	pgpagination "github.com/kichiyaki/pg-pagination"
)

type Repository interface {
	Store(ctx context.Context, u *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetBySlug(ctx context.Context, slug string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByCredentials(ctx context.Context, login, password string) (*models.User, error)
	Fetch(ctx context.Context,
		p pgpagination.Pagination,
		f *pgfilter.Filter) (*models.List, error)
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, ids []int) ([]*models.User, error)
}
