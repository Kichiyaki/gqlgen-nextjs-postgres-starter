package mocks

import (
	"context"

	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/stretchr/testify/mock"
)

type Usecase struct {
	mock.Mock
}

func (_m *Usecase) Signup(ctx context.Context, input models.UserInput) (*models.User, error) {
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

func (_m *Usecase) Login(ctx context.Context, login, password string) (*models.User, error) {
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

func (_m *Usecase) Logout(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *Usecase) Activate(ctx context.Context, id int, token string) (*models.User, error) {
	ret := _m.Called(ctx, id, token)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(context.Context, int, string) *models.User); ok {
		r0 = rf(ctx, id, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(ctx, id, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *Usecase) GenerateNewActivationToken(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *Usecase) IsLogged(ctx context.Context) bool {
	ret := _m.Called(ctx)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context) bool); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(bool)
		}
	}

	return r0
}

func (_m *Usecase) HasAdministrativePrivileges(ctx context.Context) bool {
	ret := _m.Called(ctx)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context) bool); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(bool)
		}
	}

	return r0
}

func (_m *Usecase) CurrentUser(ctx context.Context) *models.User {
	ret := _m.Called(ctx)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(context.Context) *models.User); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	return r0
}
