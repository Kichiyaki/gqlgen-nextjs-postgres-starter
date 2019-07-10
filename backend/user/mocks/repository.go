package mocks

import (
	"context"

	pgfilter "github.com/kichiyaki/pg-filter"
	pgpagination "github.com/kichiyaki/pg-pagination"

	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func (_m *Repository) Store(ctx context.Context, u *models.User) error {
	ret := _m.Called(ctx, u)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.User) error); ok {
		r0 = rf(ctx, u)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *Repository) GetByID(ctx context.Context, id int) (*models.User, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(context.Context, int) *models.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *Repository) GetBySlug(ctx context.Context, slug string) (*models.User, error) {
	ret := _m.Called(ctx, slug)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.User); ok {
		r0 = rf(ctx, slug)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, slug)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *Repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ret := _m.Called(ctx, email)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *Repository) GetByCredentials(ctx context.Context, login, password string) (*models.User, error) {
	ret := _m.Called(ctx, login, password)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *models.User); ok {
		r0 = rf(ctx, login, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, login, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *Repository) Fetch(ctx context.Context,
	p pgpagination.Pagination,
	f *pgfilter.Filter) (*models.List, error) {
	ret := _m.Called(ctx, p, f)

	var r0 *models.List
	if rf, ok := ret.Get(0).(func(context.Context, pgpagination.Pagination, *pgfilter.Filter) *models.List); ok {
		r0 = rf(ctx, p, f)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.List)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pgpagination.Pagination, *pgfilter.Filter) error); ok {
		r1 = rf(ctx, p, f)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *Repository) Update(ctx context.Context, u *models.User) error {
	ret := _m.Called(ctx, u)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.User) error); ok {
		r0 = rf(ctx, u)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *Repository) Delete(ctx context.Context, ids []int) ([]*models.User, error) {
	ret := _m.Called(ctx, ids)

	var r0 []*models.User
	if rf, ok := ret.Get(0).(func(context.Context, []int) []*models.User); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.User)
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
