package resolvers

import (
	"backend/auth"
	"backend/errors"
	"backend/middleware"
	"backend/models"
	"backend/utils"
	"context"
	"fmt"

	"github.com/labstack/echo-contrib/session"
)

const (
	activateAccountEmailTitle   = "activateAccountEmailTitle"
	activateAccountEmailContent = "activateAccountEmailContent"
	resetPasswordEmailTitle     = "resetPasswordEmailTitle"
	resetPasswordEmailContent   = "resetPasswordEmailContent"
	passwordChangedEmailTitle   = "passwordChangedEmailTitle"
	passwordChangedEmailContent = "passwordChangedEmailContent"
)

func (r *mutationResolver) Signup(ctx context.Context, input models.UserInput) (*models.User, error) {
	user, err := r.AuthUcase.Signup(ctx, input)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, err)
	}

	echoCtx, err := middleware.EchoContextFromContext(ctx)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, errors.Wrap(errors.ErrInternalServerError, err))
	}
	sess, err := session.Get(auth.SessionName, echoCtx)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, errors.Wrap(errors.ErrInternalServerError, err))
	}
	sess.Values["login"] = user.Login
	sess.Values["password"] = input.Password
	sess.Save(echoCtx.Request(), echoCtx.Response())
	go func() {
		sendEmail(ctx,
			activateAccountEmailTitle,
			activateAccountEmailContent,
			user.Email,
			map[string]interface{}{
				"Login": user.Login,
				"Href":  fmt.Sprintf("%s/%d/activate/%s", r.FrontendURL, user.ID, user.ActivationToken),
			})
	}()

	return user, nil
}

func (r *mutationResolver) Signin(ctx context.Context, login string, password string) (*models.User, error) {
	user, err := r.AuthUcase.Signin(ctx, login, password)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, err)
	}

	echoCtx, err := middleware.EchoContextFromContext(ctx)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, errors.Wrap(errors.ErrInternalServerError, err))
	}
	sess, err := session.Get(auth.SessionName, echoCtx)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, errors.Wrap(errors.ErrInternalServerError, err))
	}
	sess.Values["login"] = login
	sess.Values["password"] = password
	sess.Save(echoCtx.Request(), echoCtx.Response())
	return user, nil
}

func (r *mutationResolver) Signout(ctx context.Context) (*string, error) {
	echoCtx, err := middleware.EchoContextFromContext(ctx)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, errors.Wrap(errors.ErrInternalServerError, err))
	}
	sess, err := session.Get(auth.SessionName, echoCtx)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, errors.Wrap(errors.ErrInternalServerError, err))
	}
	delete(sess.Values, "login")
	delete(sess.Values, "password")
	sess.Save(echoCtx.Request(), echoCtx.Response())
	msg := "Success"
	return &msg, nil
}

func (r *mutationResolver) GenerateNewActivationTokenForMe(ctx context.Context) (*string, error) {
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, errors.Wrap(errors.ErrMustBeLoggedIn, err))
	}
	user, err = r.AuthUcase.GenerateNewActivationToken(ctx, user.ID)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, err)
	}
	go func() {
		sendEmail(ctx,
			activateAccountEmailTitle,
			activateAccountEmailContent,
			user.Email,
			map[string]interface{}{
				"Login": user.Login,
				"Href":  fmt.Sprintf("%s/%d/activate/%s", r.FrontendURL, user.ID, user.ActivationToken),
			})
	}()
	msg := "Success"
	return &msg, nil
}

func (r *mutationResolver) GenerateNewResetPasswordToken(ctx context.Context, email string) (*string, error) {
	user, err := r.AuthUcase.GenerateNewResetPasswordToken(ctx, email)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, err)
	}
	go func() {
		sendEmail(ctx,
			resetPasswordEmailTitle,
			resetPasswordEmailContent,
			user.Email,
			map[string]interface{}{
				"Login": user.Login,
				"Href":  fmt.Sprintf("%s/%d/reset-password/%s", r.FrontendURL, user.ID, user.ResetPasswordToken),
			})
	}()
	msg := "Success"
	return &msg, nil
}

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	user, _ := middleware.UserFromContext(ctx)
	return user, nil
}

func (r *queryResolver) ActivateUserAccount(ctx context.Context, id int, token string) (*models.User, error) {
	user, err := r.AuthUcase.Activate(ctx, id, token)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, err)
	}
	return user, nil
}

func (r *queryResolver) ResetUserPassword(ctx context.Context, id int, token string) (*string, error) {
	user, password, err := r.AuthUcase.ResetPassword(ctx, id, token)
	if err != nil {
		return nil, utils.FormatErrorMsg(ctx, err)
	}
	go func() {
		sendEmail(ctx,
			passwordChangedEmailTitle,
			passwordChangedEmailContent,
			user.Email,
			map[string]interface{}{
				"Login":    user.Login,
				"Password": password,
			})
	}()
	msg := "Success"
	return &msg, nil
}
