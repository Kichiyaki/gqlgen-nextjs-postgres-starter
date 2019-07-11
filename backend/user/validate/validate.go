package validate

import (
	"fmt"

	"regexp"

	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/kichiyaki/graphql-starter/backend/user/errors"
)

const (
	emailRegex              = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	containsUppercaseRegex  = "[A-ZŻŹĆĄŚĘŁÓŃ]+"
	containsLowercaseRegex  = "[a-zzżźćńółęąś]+"
	containsDigitRegex      = `\d+`
	minimalLengthOfLogin    = 3
	maximalLengthOfLogin    = 72
	minimalLengthOfPassword = 8
	maximalLengthOfPassword = 128
)

type UserValidationConfig struct {
	Login    bool
	Password bool
	Email    bool
	Role     bool
}

func (cfg UserValidationConfig) IsSomethingToValidate() bool {
	return cfg.Login || cfg.Password || cfg.Role || cfg.Email
}

func (cfg UserValidationConfig) Validate(user models.User) error {
	if cfg.Login {
		if user.Login == "" {
			return errors.ErrMustProvideLogin
		} else if len(user.Login) < minimalLengthOfLogin {
			return fmt.Errorf(errors.MinimalLengthOfLoginErrFormat, minimalLengthOfLogin)
		} else if len(user.Login) > maximalLengthOfLogin {
			return fmt.Errorf(errors.MaximalLengthOfLoginErrFormat, maximalLengthOfLogin)
		}
	}
	if cfg.Password {
		if user.Password == "" {
			return errors.ErrMustProvidePassword
		} else if len(user.Password) < minimalLengthOfPassword {
			return fmt.Errorf(errors.MinimalLengthOfPasswordErrFormat, minimalLengthOfPassword)
		} else if len(user.Password) > maximalLengthOfPassword {
			return fmt.Errorf(errors.MaximalLengthOfPasswordErrFormat, maximalLengthOfPassword)
		} else if matched, _ := regexp.Match(containsUppercaseRegex, []byte(user.Password)); !matched {
			return errors.ErrPasswordMustContainAtLeastOneUppercase
		} else if matched, _ := regexp.Match(containsLowercaseRegex, []byte(user.Password)); !matched {
			return errors.ErrPasswordMustContainsAtLeastOneLowercase
		} else if matched, _ := regexp.Match(containsDigitRegex, []byte(user.Password)); !matched {
			return errors.ErrPasswordMustContainsAtLeastOneDigit
		}
	}
	if cfg.Email {
		if user.Email == "" {
			return errors.ErrMustProvideEmailAddress
		} else if matched, _ := regexp.Match(emailRegex, []byte(user.Email)); !matched {
			return errors.ErrInvalidEmailAddress
		}
	}
	if cfg.Role {
		if user.Role != models.DefaultRole && user.Role != models.AdministrativeRole {
			return fmt.Errorf(errors.InvalidRoleErrFormat, user.Role)
		}
	}

	return nil
}
