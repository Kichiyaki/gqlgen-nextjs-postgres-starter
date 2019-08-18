package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	AccountActivationTokenType = "activate_account"
	ResetPasswordTokenType     = "reset_password"
)

type Token struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UserID    int       `json:"userID" gqlgen:"user"`
	User      *User     `json:"user" gqlgen:"-"`
}

func NewToken(t string, uid int) *Token {
	return &Token{
		Type:      t,
		Value:     uuid.New().String(),
		CreatedAt: time.Now(),
		UserID:    uid,
	}
}

func (t *Token) BeforeInsert(c context.Context) (context.Context, error) {
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}

	return c, nil
}
