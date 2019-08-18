package mocks

import (
	"net/http"

	ginSessions "github.com/gin-contrib/sessions"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/mock"
)

type Store struct {
	mock.Mock
}

func (_m *Store) Get(r *http.Request, name string) (*sessions.Session, error) {
	ret := _m.Called(r, name)

	var r0 *sessions.Session
	if rf, ok := ret.Get(0).(func(*http.Request, string) *sessions.Session); ok {
		r0 = rf(r, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sessions.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*http.Request, string) error); ok {
		r1 = rf(r, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *Store) New(r *http.Request, name string) (*sessions.Session, error) {
	ret := _m.Called(r, name)

	var r0 *sessions.Session
	if rf, ok := ret.Get(0).(func(*http.Request, string) *sessions.Session); ok {
		r0 = rf(r, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sessions.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*http.Request, string) error); ok {
		r1 = rf(r, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *Store) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	ret := _m.Called(r, w, session)

	var r0 error
	if rf, ok := ret.Get(0).(func(*http.Request, http.ResponseWriter, *sessions.Session) error); ok {
		r0 = rf(r, w, session)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *Store) Delete(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	ret := _m.Called(r, w, session)

	var r0 error
	if rf, ok := ret.Get(0).(func(*http.Request, http.ResponseWriter, *sessions.Session) error); ok {
		r0 = rf(r, w, session)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *Store) DeleteByID(ids ...string) error {
	ret := _m.Called(ids)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(ids)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *Store) GetAll() ([]*sessions.Session, error) {
	ret := _m.Called()

	var r0 []*sessions.Session
	if rf, ok := ret.Get(0).(func() []*sessions.Session); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*sessions.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *Store) Options(options ginSessions.Options) {

}
