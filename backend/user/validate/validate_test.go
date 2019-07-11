package validate

import (
	"fmt"
	"testing"

	"github.com/kichiyaki/graphql-starter/backend/seed"

	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/kichiyaki/graphql-starter/backend/user/errors"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	cfg := UserValidationConfig{true, true, true, true}

	t.Run("Login", func(t *testing.T) {
		t.Run("must provide login", func(t *testing.T) {
			err := cfg.Validate(models.User{})
			require.Equal(t, errors.ErrMustProvideLogin, err)
		})

		t.Run(fmt.Sprintf("login must be at least %d characters", minimalLengthOfLogin), func(t *testing.T) {
			err := cfg.Validate(models.User{Login: "a"})
			require.Equal(t,
				fmt.Errorf(errors.MinimalLengthOfLoginErrFormat, minimalLengthOfLogin),
				err)
		})

		t.Run(fmt.Sprintf("login must be shorter than %d characters", maximalLengthOfLogin), func(t *testing.T) {
			l := "a"
			for i := 1; i <= maximalLengthOfLogin+5; i++ {
				l += "b"
			}
			err := cfg.Validate(models.User{Login: l})
			require.Equal(t,
				fmt.Errorf(errors.MaximalLengthOfLoginErrFormat, maximalLengthOfLogin),
				err)
		})
	})

	t.Run("Password", func(t *testing.T) {
		t.Run("must provide password", func(t *testing.T) {
			err := cfg.Validate(models.User{Login: "test12"})
			require.Equal(t, errors.ErrMustProvidePassword, err)
		})

		t.Run(fmt.Sprintf("password must be at least %d characters", minimalLengthOfPassword), func(t *testing.T) {
			err := cfg.Validate(models.User{Login: "abcd", Password: "12"})
			require.Equal(t,
				fmt.Errorf(errors.MinimalLengthOfPasswordErrFormat, minimalLengthOfPassword),
				err)
		})

		t.Run(fmt.Sprintf("password must be shorter than %d characters", maximalLengthOfPassword), func(t *testing.T) {
			l := "a"
			for i := 1; i <= maximalLengthOfPassword+5; i++ {
				l += "b"
			}
			err := cfg.Validate(models.User{Login: "asds", Password: l})
			require.Equal(t,
				fmt.Errorf(errors.MaximalLengthOfPasswordErrFormat, maximalLengthOfPassword),
				err)
		})

		t.Run("password must contain at least one uppercase letter", func(t *testing.T) {
			err := cfg.Validate(models.User{Login: "asds", Password: "asdasdadasda"})
			require.Equal(t, errors.ErrPasswordMustContainAtLeastOneUppercase, err)
		})

		t.Run("password must contain at least one lowercase letter", func(t *testing.T) {
			err := cfg.Validate(models.User{Login: "asds", Password: "ASDASDADSADA"})
			require.Equal(t, errors.ErrPasswordMustContainsAtLeastOneLowercase, err)
		})

		t.Run("password must contain at least one digit", func(t *testing.T) {
			err := cfg.Validate(models.User{Login: "asds", Password: "ASDASDADSADAasd"})
			require.Equal(t, errors.ErrPasswordMustContainsAtLeastOneDigit, err)
		})
	})

	t.Run("Email", func(t *testing.T) {
		t.Run("must provide email address", func(t *testing.T) {
			u := seed.Users()[0]
			u.Email = ""
			err := cfg.Validate(u)
			require.Equal(t, errors.ErrMustProvideEmailAddress, err)
		})

		t.Run("invalid email address", func(t *testing.T) {
			u := seed.Users()[0]
			u.Email = "tesasd2sd..as"
			err := cfg.Validate(u)
			require.Equal(t, errors.ErrInvalidEmailAddress, err)
		})
	})

	t.Run("Role", func(t *testing.T) {
		t.Run("invalid role", func(t *testing.T) {
			u := seed.Users()[0]
			u.Role = "eloszka"
			err := cfg.Validate(u)
			require.Equal(t, fmt.Errorf(errors.InvalidRoleErrFormat, u.Role), err)
		})
	})

	t.Run("success", func(t *testing.T) {
		u := seed.Users()[0]
		err := cfg.Validate(u)
		require.Equal(t, nil, err)
	})
}
