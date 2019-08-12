package sessions

import (
	ginSessions "github.com/gin-contrib/sessions"
	"github.com/gorilla/sessions"
)

type Store interface {
	ginSessions.Store
	GetAll() ([]*sessions.Session, error)
	DeleteByID(ids ...string) error
}
