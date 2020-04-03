package user

import (
	"context"

	"backend/models"
)

type Repository interface {
	Fetch(ctx context.Context, f *models.UserFilter) (models.UserList, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetBySlug(ctx context.Context, slug string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByCredentials(ctx context.Context, login, password string) (*models.User, error)
	Update(ctx context.Context, u *models.User) error
	Store(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, f *models.UserFilter) ([]*models.User, error)
}
