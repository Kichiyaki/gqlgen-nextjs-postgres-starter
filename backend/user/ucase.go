package user

import (
	"context"

	"backend/models"
)

type Usecase interface {
	Fetch(ctx context.Context, f *models.UserFilter) (models.UserList, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetBySlug(ctx context.Context, slug string) (*models.User, error)
	Update(ctx context.Context, id int, input models.UserInput) (*models.User, error)
	Store(ctx context.Context, input models.UserInput) (*models.User, error)
	Delete(ctx context.Context, ids ...int) ([]*models.User, error)
}
