package validation

import (
	"regexp"

	_errors "backend/errors"
	"backend/models"
)

const (
	MinimumPasswordLength = 6
	MaximumPasswordLength = 64
	MinimumLoginLength    = 2
	MaximumLoginLength    = 128
	emailRegex            = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	containUppercaseRegex = "[A-ZŻŹĆĄŚĘŁÓŃ]+"
	containLowercaseRegex = "[a-zzżźćńółęąś]+"
	containDigitRegex     = `\d+`
)

type Config struct {
	Login    bool
	Password bool
	Email    bool
	Role     bool
}

func NewConfig() Config {
	return Config{
		true, true, true, true,
	}
}

func (c Config) Validate(u models.User) error {
	if c.Login && (len(u.Login) < MinimumLoginLength || len(u.Login) > MaximumLoginLength) {
		return _errors.Wrap(_errors.ErrLoginPolicy)
	}

	if c.Password {
		if len(u.Password) < MinimumPasswordLength || len(u.Password) > MaximumPasswordLength {
			return _errors.Wrap(_errors.ErrPasswordPolicy)
		} else if matched, _ := regexp.Match(containUppercaseRegex, []byte(u.Password)); !matched {
			return _errors.Wrap(_errors.ErrPasswordPolicy)

		} else if matched, _ := regexp.Match(containLowercaseRegex, []byte(u.Password)); !matched {
			return _errors.Wrap(_errors.ErrPasswordPolicy)

		} else if matched, _ := regexp.Match(containDigitRegex, []byte(u.Password)); !matched {
			return _errors.Wrap(_errors.ErrPasswordPolicy)
		}
	}

	if c.Email {
		if u.Email == "" {
			return _errors.Wrap(_errors.ErrEmailPolicy)
		} else if matched, _ := regexp.Match(emailRegex, []byte(u.Email)); !matched {
			return _errors.Wrap(_errors.ErrEmailPolicy)
		}
	}

	if c.Role {
		if u.Role != models.UserDefaultRole && u.Role != models.UserAdminRole {
			return _errors.Wrap(_errors.ErrInvalidUserRole)
		}
	}

	return nil
}
