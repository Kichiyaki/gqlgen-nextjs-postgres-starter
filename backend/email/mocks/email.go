package mocks

import (
	"github.com/kichiyaki/graphql-starter/backend/email"
	"github.com/stretchr/testify/mock"
)

type Email struct {
	mock.Mock
}

func (_m *Email) Send(cfg *email.EmailConfig) error {
	ret := _m.Called(cfg)

	var r0 error
	if rf, ok := ret.Get(0).(func(*email.EmailConfig) error); ok {
		r0 = rf(cfg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
