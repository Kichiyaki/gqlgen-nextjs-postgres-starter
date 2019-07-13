package token

import (
	"context"

	"github.com/kichiyaki/graphql-starter/backend/models"
	pgfilter "github.com/kichiyaki/pg-filter"
)

type Repository interface {
	Store(ctx context.Context, t *models.Token) error
	Fetch(ctx context.Context, f *pgfilter.Filter) ([]*models.Token, error)
	Delete(ctx context.Context, ids []int) ([]*models.Token, error)
	DeleteByUserID(ctx context.Context, t string, id int) ([]*models.Token, error)
}
