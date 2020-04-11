package resolvers

import (
	"backend/errors"
	"backend/models"
	"backend/utils"
	"context"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
	user, err := r.UserUcase.Store(ctx, input)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, err)
	}
	return user, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id int, input models.UserInput) (*models.User, error) {
	user, err := r.UserUcase.Update(ctx, id, input)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, err)
	}
	return user, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, ids []int) ([]*models.User, error) {
	users, err := r.UserUcase.Delete(ctx, ids...)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, err)
	}
	return users, nil
}

func (r *queryResolver) Users(ctx context.Context, filter *models.UserFilter) (*models.UserList, error) {
	list, err := r.UserUcase.Fetch(ctx, filter)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, err)
	}
	return &list, nil
}

func (r *queryResolver) User(ctx context.Context, id *int, slug *string) (*models.User, error) {
	var user *models.User
	var err error
	if id != nil {
		user, err = r.UserUcase.GetByID(ctx, *id)
	} else if slug != nil {
		user, err = r.UserUcase.GetBySlug(ctx, *slug)
	} else {
		err = errors.Wrap(errors.ErrInvalidPayload)
	}
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, err)
	}
	return user, nil
}
