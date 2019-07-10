package mocks

import (
	"context"

	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/stretchr/testify/mock"
)

type Usecase struct {
	mock.Mock
}

func (_m *Usecase) Store(ctx context.Context, input models.UserInput) (*models.User, error) {
	ret := _m.Called(ctx, input)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(context.Context, models.UserInput) *models.User); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.UserInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *Usecase) GetByID(ctx context.Context, id int) (*models.User, error) {
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

func (_m *Usecase) GetBySlug(ctx context.Context, slug string) (*models.User, error) {
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

func (_m *Usecase) Update(ctx context.Context, id int, input models.UserInput) (*models.User, error) {
	ret := _m.Called(ctx, id, input)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(context.Context, int, models.UserInput) *models.User); ok {
		r0 = rf(ctx, id, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, models.UserInput) error); ok {
		r1 = rf(ctx, id, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *Usecase) Delete(ctx context.Context, ids []int) ([]*models.User, error) {
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
