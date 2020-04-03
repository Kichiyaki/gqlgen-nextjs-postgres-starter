package validation

import (
	"strings"
	"testing"

	_errors "backend/errors"
	"backend/utils/seed"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	cfg := NewConfig()
	u := seed.Users(1)[0]

	t.Run("passed validation", func(t *testing.T) {
		err := cfg.Validate(u)
		require.Equal(t, nil, err)
	})

	t.Run("login policy test", func(t *testing.T) {
		copy := u
		t.Run("login is too short", func(t *testing.T) {
			copy.Login = "a"
			err := cfg.Validate(copy)
			require.Equal(t, true, strings.Contains(err.Error(), _errors.ErrLoginPolicy))
		})
		t.Run("login is too long", func(t *testing.T) {
			for i := 0; i < MaximumLoginLength+1; i++ {
				copy.Login += "s"
			}
			err := cfg.Validate(copy)
			require.Equal(t, true, strings.Contains(err.Error(), _errors.ErrLoginPolicy))
		})
	})

	t.Run("password policy test", func(t *testing.T) {
		copy := u
		t.Run("password is too short", func(t *testing.T) {
			copy.Password = "a"
			err := cfg.Validate(copy)
			require.Equal(t, true, strings.Contains(err.Error(), _errors.ErrPasswordPolicy))
		})
		t.Run("password is too long", func(t *testing.T) {
			for i := 0; i < MaximumPasswordLength+1; i++ {
				copy.Password += "s"
			}
			err := cfg.Validate(copy)
			require.Equal(t, true, strings.Contains(err.Error(), _errors.ErrPasswordPolicy))
		})
		t.Run("the password does not contain numbers", func(t *testing.T) {
			copy.Password = "asdasdadsaASDASD"
			err := cfg.Validate(copy)
			require.Equal(t, true, strings.Contains(err.Error(), _errors.ErrPasswordPolicy))
		})
		t.Run("the password does not contain uppercase", func(t *testing.T) {
			copy.Password = "asdasdadsa123"
			err := cfg.Validate(copy)
			require.Equal(t, true, strings.Contains(err.Error(), _errors.ErrPasswordPolicy))
		})
		t.Run("the password does not contain lowercase", func(t *testing.T) {
			copy.Password = "ASDASDADSAX123"
			err := cfg.Validate(copy)
			require.Equal(t, true, strings.Contains(err.Error(), _errors.ErrPasswordPolicy))
		})
	})

	t.Run("email policy test", func(t *testing.T) {
		copy := u
		t.Run("email string is empty", func(t *testing.T) {
			copy.Email = ""
			err := cfg.Validate(copy)
			require.Equal(t, true, strings.Contains(err.Error(), _errors.ErrEmailPolicy))
		})
		t.Run("email is invalid", func(t *testing.T) {
			copy.Email = "asdsd"
			err := cfg.Validate(copy)
			require.Equal(t, true, strings.Contains(err.Error(), _errors.ErrEmailPolicy))
		})
	})

	t.Run("role is invalid", func(t *testing.T) {
		copy := u
		copy.Role = 125
		err := cfg.Validate(copy)
		require.Equal(t, true, strings.Contains(err.Error(), _errors.ErrInvalidUserRole))
	})
}
