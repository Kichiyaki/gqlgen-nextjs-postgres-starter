package errors

const (
	ErrUserNotFound       = "user.notFoundError"
	ErrInvalidCredentials = "user.invalidCredentialsError"
	ErrLoginMustBeUnique  = "user.loginMustBeUniqueError"
	ErrEmailMustBeUnique  = "user.emailMustBeUniqueError"
	ErrLoginPolicy        = "user.loginPolicyError"
	ErrDisplayNamePolicy  = "user.displayNamePolicyError"
	ErrPasswordPolicy     = "user.passwordPolicyError"
	ErrEmailPolicy        = "user.emailPolicyError"
	ErrInvalidUserRole    = "user.invalidUserRoleError"
)
