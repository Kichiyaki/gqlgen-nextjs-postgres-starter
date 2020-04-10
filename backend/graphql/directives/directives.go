package directives

import (
	"backend/errors"
	"backend/middleware"
	"backend/utils"
	"context"

	"github.com/99designs/gqlgen/graphql"
)

type Handler struct {
}

func (h *Handler) Activated(ctx context.Context, obj interface{}, next graphql.Resolver, yes bool) (interface{}, error) {
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, errors.Wrap(errors.ErrMustBeLoggedIn, err))
	}
	if yes && !user.Activated {
		return nil, utils.FormatErrorMsg(ctx, errors.Wrap(errors.ErrMustHaveActivatedAccount))
	} else if !yes && user.Activated {
		return nil, utils.FormatErrorMsg(ctx, errors.Wrap(errors.ErrMustHaveDeactivatedAccount))
	}

	return next(ctx)
}

func (h *Handler) Authenticated(ctx context.Context, obj interface{}, next graphql.Resolver, yes bool) (interface{}, error) {
	_, err := middleware.UserFromContext(ctx)
	if yes && err != nil {
		return nil, utils.FormatErrorMsg(ctx, errors.Wrap(errors.ErrMustBeLoggedIn, err))
	} else if !yes && err == nil {
		return nil, utils.FormatErrorMsg(ctx, errors.Wrap(errors.ErrMustBeLoggedOut))
	}

	return next(ctx)
}

func (h *Handler) HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, role int) (interface{}, error) {
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, errors.Wrap(errors.ErrMustBeLoggedIn, err))
	}
	if user.Role != role {
		return nil, utils.FormatErrorMsg(ctx, errors.Wrap(errors.ErrUnauthorized))
	}

	return next(ctx)
}
