package token

import (
	"context"

	"github.com/kichiyaki/graphql-starter/backend/models"
)

type Repository interface {
	Store(ctx context.Context, t *models.Token) error
	Get(ctx context.Context, t, value string) (*models.Token, error)
	Delete(ctx context.Context, ids []int) ([]*models.Token, error)
	DeleteByUserID(ctx context.Context, t string, id int) ([]*models.Token, error)
}
