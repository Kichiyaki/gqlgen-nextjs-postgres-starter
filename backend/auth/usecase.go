package auth

import (
	"context"

	"github.com/kichiyaki/graphql-starter/backend/models"
	ginSessions "github.com/kichiyaki/sessions/gin-sessions"
)

type Usecase interface {
	Signup(ctx context.Context, input models.UserInput) (*models.User, error)
	Login(ctx context.Context, login, password string) (*models.User, error)
	Logout(ctx context.Context) error
	Activate(ctx context.Context, id int, token string) (*models.User, error)
	GenerateNewActivationToken(ctx context.Context, id int) error
	GenerateNewActivationTokenForCurrentUser(ctx context.Context) error
	GenerateNewResetPasswordToken(ctx context.Context, emailAddress string) error
	ResetPassword(ctx context.Context, id int, token string) error
	ChangePassword(ctx context.Context, currentPassword, newPassword string) error
	IsLogged(ctx context.Context) bool
	HasAdministrativePrivileges(ctx context.Context) bool
	CurrentUser(ctx context.Context) *models.User
	Session(ctx context.Context) ginSessions.Session
}
