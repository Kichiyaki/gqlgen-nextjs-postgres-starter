package auth

import (
	"context"

	"github.com/kichiyaki/graphql-starter/backend/models"
)

type Usecase interface {
	Signup(ctx context.Context, input models.UserInput) (*models.User, error)
	Login(ctx context.Context, login, password string) (*models.User, error)
	Activate(ctx context.Context, id int, token string) (*models.User, error)
	GenerateNewActivationToken(ctx context.Context, id int) error

	IsLogged(ctx context.Context) bool
	HasAdministrativePrivileges(ctx context.Context) bool
	CurrentUser(ctx context.Context) *models.User
}
