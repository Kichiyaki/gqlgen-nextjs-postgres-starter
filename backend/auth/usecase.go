package auth

import (
	"context"

	"backend/models"
)

type Usecase interface {
	Signup(ctx context.Context, input models.UserInput) (*models.User, error)
	Signin(ctx context.Context, login, password string) (*models.User, error)
	GenerateNewActivationToken(ctx context.Context, id int) (*models.User, error)
	Activate(ctx context.Context, id int, token string) (*models.User, error)
	GenerateNewResetPasswordToken(ctx context.Context, email string) (*models.User, error)
	ResetPassword(ctx context.Context, id int, token string) (*models.User, string, error)
}
