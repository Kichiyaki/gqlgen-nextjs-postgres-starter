package models

import (
	"context"
	"time"

	_errors "backend/errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	UserAdminRole   = 2
	UserDefaultRole = 1
)

type User struct {
	tableName struct{} `pg:"alias:u"`

	ID                            int       `json:"id" pg:",pk"`
	Slug                          string    `json:"slug" pg:",unique,use_zero"`
	Login                         string    `json:"login" pg:",unique,use_zero"`
	Password                      string    `json:"-" gqlgen:"-"`
	DisplayName                   string    `json:"displayName" pg:",use_zero"`
	Email                         string    `json:"email" pg:",unique"`
	CreatedAt                     time.Time `json:"createdAt" pg:"default:now()"`
	UpdatedAt                     time.Time `json:"updatedAt" pg:"default:now()"`
	Role                          int       `json:"role"`
	Activated                     bool      `json:"activated" pg:"default:false,use_zero"`
	ActivationToken               string    `json:"-" gqlgen:"-"`
	ActivationTokenGeneratedAt    time.Time `json:"-" gqlgen:"-" pg:"default:now()"`
	ResetPasswordToken            string    `json:"-" gqlgen:"-"`
	ResetPasswordTokenGeneratedAt time.Time `json:"-" gqlgen:"-" pg:"default:now()"`
}

func (u *User) CompareHashAndPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil && password != u.Password {
		return _errors.Wrap(_errors.ErrInvalidCredentials)
	}
	return nil
}

func (u *User) MergeInput(input UserInput) {
	if input.Login != "" {
		u.Login = input.Login
	}
	if input.Password != "" {
		u.Password = input.Password
	}
	if input.DisplayName != "" {
		u.DisplayName = input.DisplayName
	}
	if input.Role > 0 {
		u.Role = input.Role
	}
	if input.Activated != nil {
		u.Activated = *input.Activated
	}
}

func (u *User) BeforeInsert(ctx context.Context) (context.Context, error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx, err
	}
	u.Password = string(hashedPassword)

	return ctx, nil
}

func (u *User) BeforeUpdate(ctx context.Context) (context.Context, error) {
	u.UpdatedAt = time.Now()

	if cost, _ := bcrypt.Cost([]byte(u.Password)); u.Password != "" && cost == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return ctx, err
		}
		u.Password = string(hashedPassword)
	}

	return ctx, nil
}

type UserInput struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Role        int    `json:"role"`
	Activated   *bool  `json:"activated"`
}

func (input UserInput) ToUser() User {
	u := User{
		Login:       input.Login,
		Password:    input.Password,
		DisplayName: input.DisplayName,
		Email:       input.Email,
		Role:        input.Role,
	}
	if input.Activated != nil {
		u.Activated = *input.Activated
	}
	return u
}

type UserFilter struct {
	tableName struct{} `urlstruct:"u"`

	ID               []int     `gqlgen:"id"`
	IdNEQ            []int     `gqlgen:"idNeq"`
	Login            []string  `gqlgen:"login"`
	LoginNEQ         []string  `gqlgen:"loginNeq"`
	LoginMATCH       string    `gqlgen:"loginMatch"`
	DisplayName      []string  `gqlgen:"displayName"`
	DisplayNameNEQ   []string  `gqlgen:"displayNameNeq"`
	DisplayNameMATCH string    `gqlgen:"displayNameMatch"`
	Email            []string  `gqlgen:"email"`
	EmailNEQ         []string  `gqlgen:"emailNeq"`
	EmailMATCH       string    `gqlgen:"emailMatch"`
	Role             []int     `gqlgen:"role"`
	Activated        string    `urlstruct:",nowhere" gqlgen:"activated"`
	CreatedAt        time.Time `gqlgen:"createdAt"`
	CreatedAtGT      time.Time `gqlgen:"createdAtGt"`
	CreatedAtLT      time.Time `gqlgen:"createdAtLt"`
	UpdatedAt        time.Time `gqlgen:"updatedAt"`
	UpdatedAtGT      time.Time `gqlgen:"updatedAtGt"`
	UpdatedAtLT      time.Time `gqlgen:"updatedAtLt"`
	Offset           int       `urlstruct:",nowhere"`
	Limit            int       `urlstruct:",nowhere"`
	Order            []string  `urlstruct:",nowhere"`
}
