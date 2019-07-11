package resolvers

import (
	"context"
	"fmt"

	"github.com/kichiyaki/graphql-starter/backend/models"
)

type userListResolver struct{ *Resolver }

func (r *userListResolver) Items(ctx context.Context, obj *models.List) ([]*models.User, error) {
	if obj == nil {
		return []*models.User{}, nil
	}

	if users, ok := obj.Items.([]*models.User); ok {
		return users, nil
	}

	return []*models.User{}, nil
}

func (r *queryResolver) FetchUser(ctx context.Context, id *int, slug *string) (*models.User, error) {
	if id != nil {
		user, err := r.UserUcase.GetByID(ctx, *id)
		if err != nil {
			return nil, err
		}
		return user, nil
	} else if slug != nil {
		user, err := r.UserUcase.GetBySlug(ctx, *slug)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	return nil, fmt.Errorf("Brak danych do wyszukania użytkownika")
}

func (r *queryResolver) FetchUsers(ctx context.Context, pagination models.Pagination, filter *models.UserFilter) (*models.List, error) {
	return r.UserUcase.Fetch(ctx, pagination, filter)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id int, user models.UserInput) (*models.User, error) {
	return r.UserUcase.Update(ctx, id, user)
}

func (r *mutationResolver) DeleteUsers(ctx context.Context, ids []int) ([]*models.User, error) {
	return r.UserUcase.Delete(ctx, ids)
}
