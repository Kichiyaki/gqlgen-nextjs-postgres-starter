package errors

import "fmt"

var (
	NotFoundUserByIDErrFormat                  = "Nie znaleziono użytkownika o ID: %d"
	NotFoundUserBySlugErrFormat                = "Nie znaleziono użytkownika o slugu: %s"
	NotFoundUserByEmailErrFormat               = "Nie znaleziono użytkownika o adresie email: %s"
	UserCannotBeUpdatedErrFormatWithLogin      = "Nie udało się zaaktualizować użytkownika %s"
	UserCannotBeUpdatedErrFormatWithID         = "Nie udało się zaaktualizować użytkownika o ID: %d"
	MinimalLengthOfLoginErrFormat              = "Minimalna długość loginu to %d znaki"
	MaximalLengthOfLoginErrFormat              = "Maksymalna długość login to %d znaki"
	MinimalLengthOfPasswordErrFormat           = "Minimalna długość hasła to %d znaki"
	MaximalLengthOfPasswordErrFormat           = "Maksymalna długość hasła to %d znaki"
	InvalidRoleErrFormat                       = "%s jest niepoprawną rolą"
	ErrInvalidLoginOrPassword                  = fmt.Errorf("Niepoprawny login/hasło")
	ErrLoginIsOccupied                         = fmt.Errorf("Podany login jest zajęty")
	ErrEmailIsOccupied                         = fmt.Errorf("Podany email jest zajęty")
	ErrUserCannotBeCreated                     = fmt.Errorf("Wystąpił błąd podczas tworzenia użytkownika")
	ErrListOfUsersCannotBeGenerated            = fmt.Errorf("Nie udało się wygenerować listy użytkowników")
	ErrUsersCannotBeDeleted                    = fmt.Errorf("Wystąpił błąd podczas usuwania użytkowników")
	ErrNothingChanged                          = fmt.Errorf("Nie wprowadziłeś żadnych zmian w konfiguracji użytkownika")
	ErrUserCannotDeleteHisAccountByHimself     = fmt.Errorf("Nie możesz sam usunąć swojego konta")
	ErrMustProvideLogin                        = fmt.Errorf("Musisz wprowadzić login")
	ErrMustProvidePassword                     = fmt.Errorf("Musisz wprowadzić hasło")
	ErrMustProvideEmailAddress                 = fmt.Errorf("Musisz wprowadzić adres email")
	ErrInvalidEmailAddress                     = fmt.Errorf("Wprowadziłeś niepoprawny adres email")
	ErrPasswordMustContainAtLeastOneUppercase  = fmt.Errorf("Hasło musi zawierać przynajmniej jedną wielką literę")
	ErrPasswordMustContainsAtLeastOneLowercase = fmt.Errorf("Hasło musi zawierać przynajmniej jedną małą literę")
	ErrPasswordMustContainsAtLeastOneDigit     = fmt.Errorf("Hasło musi zawierać przynajmniej jedną cyfrę")
)
