package resolvers

import (
	"context"

	"github.com/kichiyaki/graphql-starter/backend/models"
)

func (r *queryResolver) FetchCurrentUser(ctx context.Context) (*models.User, error) {
	return r.AuthUcase.CurrentUser(ctx), nil
}

func (r *queryResolver) ActivateUserAccount(ctx context.Context, id int, token string) (*models.User, error) {
	return r.AuthUcase.Activate(ctx, id, token)
}

func (r *queryResolver) ResetPassword(ctx context.Context, id int, token string) (*string, error) {
	if err := r.AuthUcase.ResetPassword(ctx, id, token); err != nil {
		return nil, err
	}

	msg := "Pomyślnie zresetowano hasło"
	return &msg, nil
}

func (r *mutationResolver) Signup(ctx context.Context, user models.UserInput) (*models.User, error) {
	session := r.AuthUcase.Session(ctx)
	u, err := r.AuthUcase.Signup(ctx, user)
	if err != nil {
		return nil, err
	}
	session.Set("user", u.ID)
	if err := session.Save(); err != nil {
		return nil, err
	}

	return u, nil
}
func (r *mutationResolver) Login(ctx context.Context, login string, password string) (*models.User, error) {
	session := r.AuthUcase.Session(ctx)
	user, err := r.AuthUcase.Login(ctx, login, password)
	if err != nil {
		return nil, err
	}
	session.Set("user", user.ID)
	if err := session.Save(); err != nil {
		return nil, err
	}

	return user, nil
}
func (r *mutationResolver) Logout(ctx context.Context) (*string, error) {
	session := r.AuthUcase.Session(ctx)
	session.Delete("user")
	if err := session.Save(); err != nil {
		return nil, err
	}

	msg := "Pomyślnie wylogowano"
	return &msg, nil
}

func (r *mutationResolver) GenerateNewActivationTokenForCurrentUser(ctx context.Context) (*string, error) {
	if err := r.AuthUcase.GenerateNewActivationTokenForCurrentUser(ctx); err != nil {
		return nil, err
	}

	msg := "Pomyślnie wygenerowano"
	return &msg, nil
}

func (r *mutationResolver) ChangePassword(ctx context.Context, currentPassword string, newPassword string) (*string, error) {
	if err := r.AuthUcase.ChangePassword(ctx, currentPassword, newPassword); err != nil {
		return nil, err
	}

	msg := "Pomyślnie zmieniono hasło"
	return &msg, nil
}

func (r *mutationResolver) ActivateUserAccount(ctx context.Context, id int, token string) (*models.User, error) {
	return r.AuthUcase.Activate(ctx, id, token)
}

func (r *mutationResolver) GenerateNewResetPasswordToken(ctx context.Context, email string) (*string, error) {
	if err := r.AuthUcase.GenerateNewResetPasswordToken(ctx, email); err != nil {
		return nil, err
	}

	msg := "Pomyślnie wygenerowano"
	return &msg, nil
}
