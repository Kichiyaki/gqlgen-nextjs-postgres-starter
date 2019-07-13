package mocks

import (
	"context"

	"github.com/kichiyaki/graphql-starter/backend/models"
	pgfilter "github.com/kichiyaki/pg-filter"
	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func (_m *Repository) Store(ctx context.Context, u *models.Token) error {
	ret := _m.Called(ctx, u)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Token) error); ok {
		r0 = rf(ctx, u)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *Repository) Fetch(ctx context.Context, f *pgfilter.Filter) ([]*models.Token, error) {
	ret := _m.Called(ctx, f)

	var r0 []*models.Token
	if rf, ok := ret.Get(0).(func(context.Context, *pgfilter.Filter) []*models.Token); ok {
		r0 = rf(ctx, f)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pgfilter.Filter) error); ok {
		r1 = rf(ctx, f)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *Repository) Delete(ctx context.Context, ids []int) ([]*models.Token, error) {
	ret := _m.Called(ctx, ids)

	var r0 []*models.Token
	if rf, ok := ret.Get(0).(func(context.Context, []int) []*models.Token); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []int) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *Repository) DeleteByUserID(ctx context.Context, t string, id int) ([]*models.Token, error) {
	ret := _m.Called(ctx, t, id)

	var r0 []*models.Token
	if rf, ok := ret.Get(0).(func(context.Context, string, int) []*models.Token); ok {
		r0 = rf(ctx, t, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, int) error); ok {
		r1 = rf(ctx, t, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
