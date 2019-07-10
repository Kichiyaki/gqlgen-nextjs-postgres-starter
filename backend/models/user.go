package models

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultRole        = "Użytkownik"
	AdministrativeRole = "Administrator"
)

type User struct {
	ID        int
	Slug      string `sql:",unique"`
	Login     string `sql:",unique"`
	Email     string `sql:",unique"`
	Password  string `json:"-"`
	Role      string
	Activated bool
	CreatedAt time.Time
}

func (u *User) BeforeInsert(c context.Context) (context.Context, error) {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)

	return c, nil
}
