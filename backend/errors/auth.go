package errors

const (
	ErrMustBeLoggedIn           = "auth.mustBeLoggedInError"
	ErrMustBeLoggedOut          = "auth.mustBeLoggedOutError"
	ErrMustHaveActivatedAccount = "auth.mustHaveActivatedAccountError"
	ErrWrongActivationToken     = "auth.wrongActivationTokenError"
	ErrWrongResetPasswordToken  = "auth.wrongResetPasswordTokenError"
)
