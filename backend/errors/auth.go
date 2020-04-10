package errors

const (
	ErrMustBeLoggedIn                             = "auth.mustBeLoggedInError"
	ErrMustBeLoggedOut                            = "auth.mustBeLoggedOutError"
	ErrMustHaveActivatedAccount                   = "auth.mustHaveActivatedAccountError"
	ErrMustHaveDeactivatedAccount                 = "auth.mustHaveDeactivatedAccountError"
	ErrAccountIsActivated                         = "auth.accountIsActivatedError"
	ErrWrongActivationToken                       = "auth.wrongActivationTokenError"
	ErrWrongResetPasswordToken                    = "auth.wrongResetPasswordTokenError"
	ErrActivationTokenHasBeenGeneratedRecently    = "auth.activationTokenHasBeenGeneratedRecentlyError"
	ErrResetPasswordTokenHasBeenGeneratedRecently = "auth.resetPasswordTokenHasBeenGeneratedRecentlyError"
)
