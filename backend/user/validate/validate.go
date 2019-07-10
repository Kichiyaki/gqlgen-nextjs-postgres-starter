package validate

import (
	"fmt"

	"regexp"

	"github.com/kichiyaki/graphql-starter/backend/models"
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
			return fmt.Errorf("Musisz wprowadzić login")
		} else if len(user.Login) < 3 {
			return fmt.Errorf("Minimalna długość loginu to 3 znaki")
		} else if len(user.Login) > 72 {
			return fmt.Errorf("Maksymalna długość login to 72 znaki")
		}
	}
	if cfg.Password {
		if user.Password == "" {
			return fmt.Errorf("Musisz wprowadzić hasło")
		} else if len(user.Password) < 8 {
			return fmt.Errorf("Minimalna długość hasła to 8 znaków")
		} else if len(user.Password) > 128 {
			return fmt.Errorf("Maksymalna długość hasła to 128 znaków")
		} else if matched, _ := regexp.Match("[A-ZŻŹĆĄŚĘŁÓŃ]+", []byte(user.Password)); !matched {
			return fmt.Errorf("Hasło musi zawierać przynajmniej jedną wielką literę")
		} else if matched, _ := regexp.Match("[a-zzżźćńółęąś]+", []byte(user.Password)); !matched {
			return fmt.Errorf("Hasło musi zawierać przynajmniej jedną małą literę")
		} else if matched, _ := regexp.Match(`\d+`, []byte(user.Password)); !matched {
			return fmt.Errorf("Hasło musi zawierać przynajmniej jedną cyfrę")
		}
	}
	if cfg.Email {
		if user.Email == "" {
			return fmt.Errorf("Musisz wprowadzić adres email")
		} else if matched, _ := regexp.Match("^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$", []byte(user.Email)); !matched {
			return fmt.Errorf("Wprowadziłeś niepoprawny adres email")
		}
	}
	if cfg.Role {
		if user.Role != models.DefaultRole && user.Role != models.AdministrativeRole {
			return fmt.Errorf("%s jest niepoprawną rolą", user.Role)
		}
	}

	return nil
}
