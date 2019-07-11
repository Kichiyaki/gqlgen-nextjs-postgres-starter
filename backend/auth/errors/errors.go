package errors

import "fmt"

var (
	ErrActivationTokenCannotBeCreated   = fmt.Errorf("Nie udało się utworzyć tokenu aktywacyjnego")
	ErrAccountHasBeenActivated          = fmt.Errorf("Konto zostało już aktywowane")
	ErrInvalidActivationToken           = fmt.Errorf("Niepoprawny token aktywacyjny")
	ErrAccountCannotBeActivated         = fmt.Errorf("Wystąpił błąd podczas aktywacji konta")
	ErrNotLoggedIn                      = fmt.Errorf("Nie jesteś zalogowany")
	ErrUnauthorized                     = fmt.Errorf("Brak uprawnień")
	ErrCannotCreateAccountWhileLoggedIn = fmt.Errorf("Nie możesz utworzyć nowego konta, będąc zalogowanym")
	ErrCannotLoginWhileLoggedIn         = fmt.Errorf("Jesteś już zalogowany")
)
