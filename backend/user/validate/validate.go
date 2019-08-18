package validate

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"regexp"

	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/kichiyaki/graphql-starter/backend/utils"
)

const (
	emailRegex              = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	containUppercaseRegex   = "[A-ZŻŹĆĄŚĘŁÓŃ]+"
	containLowercaseRegex   = "[a-zzżźćńółęąś]+"
	containDigitRegex       = `\d+`
	minimalLengthOfLogin    = 3
	maximalLengthOfLogin    = 72
	minimalLengthOfPassword = 8
	maximalLengthOfPassword = 128
)

type UserValidationConfig struct {
	Login     bool
	Password  bool
	Email     bool
	Role      bool
	localizer *i18n.Localizer
}

func NewConfig(localizer *i18n.Localizer) UserValidationConfig {
	return UserValidationConfig{
		localizer: localizer,
	}
}

func (cfg UserValidationConfig) IsSomethingToValidate() bool {
	return cfg.Login || cfg.Password || cfg.Role || cfg.Email
}

func (cfg UserValidationConfig) Validate(user models.User) error {
	if cfg.Login {
		if user.Login == "" {
			return utils.GetErrorMsg(cfg.localizer, "ErrMustProvideLogin")
		} else if len(user.Login) < minimalLengthOfLogin {
			return utils.GetErrorMsgWithDataAndPluralCount(cfg.localizer, "ErrMinimalLengthOfLogin", map[string]interface{}{
				"Characters": minimalLengthOfLogin,
			}, minimalLengthOfLogin)
		} else if len(user.Login) > maximalLengthOfLogin {
			return utils.GetErrorMsgWithDataAndPluralCount(cfg.localizer, "ErrMaximalLengthOfLogin", map[string]interface{}{
				"Characters": maximalLengthOfLogin,
			}, maximalLengthOfLogin)
		}
	}
	if cfg.Password {
		if user.Password == "" {
			return utils.GetErrorMsg(cfg.localizer, "ErrMustProvidePassword")
		} else if len(user.Password) < minimalLengthOfPassword {
			return utils.GetErrorMsgWithDataAndPluralCount(cfg.localizer, "ErrMinimalLengthOfPassword", map[string]interface{}{
				"Characters": minimalLengthOfPassword,
			}, minimalLengthOfPassword)
		} else if len(user.Password) > maximalLengthOfPassword {
			return utils.GetErrorMsgWithDataAndPluralCount(cfg.localizer, "ErrMaximalLengthOfPassword", map[string]interface{}{
				"Characters": maximalLengthOfPassword,
			}, maximalLengthOfPassword)
		} else if matched, _ := regexp.Match(containUppercaseRegex, []byte(user.Password)); !matched {
			return utils.GetErrorMsg(cfg.localizer, "ErrPasswordMustContainAtLeastOneUppercase")

		} else if matched, _ := regexp.Match(containLowercaseRegex, []byte(user.Password)); !matched {
			return utils.GetErrorMsg(cfg.localizer, "ErrPasswordMustContainAtLeastOneLowercase")

		} else if matched, _ := regexp.Match(containDigitRegex, []byte(user.Password)); !matched {
			return utils.GetErrorMsg(cfg.localizer, "ErrPasswordMustContainAtLeastOneDigit")

		}
	}
	if cfg.Email {
		if user.Email == "" {
			return utils.GetErrorMsg(cfg.localizer, "ErrMustProvideEmailAddress")

		} else if matched, _ := regexp.Match(emailRegex, []byte(user.Email)); !matched {
			return utils.GetErrorMsgWithData(cfg.localizer, "ErrInvalidEmailAddress", map[string]interface{}{
				"Email": user.Email,
			})
		}
	}
	if cfg.Role {
		if user.Role != models.DefaultRole && user.Role != models.AdministrativeRole {
			return utils.GetErrorMsgWithData(cfg.localizer, "ErrInvalidRole", map[string]interface{}{
				"Role": user.Role,
			})
		}
	}

	return nil
}
