package errors

import "fmt"

var (
	ErrActivationTokenCannotBeCreated    = fmt.Errorf("Nie udało się utworzyć tokenu aktywacyjnego")
	ErrResetPasswordTokenCannotBeCreated = fmt.Errorf("Nie udało się utworzyć tokenu do zmiany hasła")
	ErrAccountHasBeenActivated           = fmt.Errorf("Konto zostało już aktywowane")
	ErrInvalidActivationToken            = fmt.Errorf("Niepoprawny token aktywacyjny")
	ErrInvalidResetPasswordToken         = fmt.Errorf("Niepoprawny token do zmiany hasła")
	ErrAccountCannotBeActivated          = fmt.Errorf("Wystąpił błąd podczas aktywacji konta")
	ErrNotLoggedIn                       = fmt.Errorf("Nie jesteś zalogowany")
	ErrUnauthorized                      = fmt.Errorf("Brak uprawnień")
	ErrCannotCreateAccountWhileLoggedIn  = fmt.Errorf("Nie możesz utworzyć nowego konta, będąc zalogowanym")
	ErrCannotLoginWhileLoggedIn          = fmt.Errorf("Jesteś już zalogowany")
	ErrReachedLimitOfActivationTokens    = fmt.Errorf("Osiągnąłeś limit wygenerowanych tokenów aktywacyjnych, prosimy spróbować później")
	ErrReachedLimitOfResetPasswordTokens = fmt.Errorf("Osiągnąłeś limit wygenerowanych tokenów do zmiany hasła, prosimy spróbować później")
	ErrCannotGeneratePassword            = fmt.Errorf("Wystąpił błąd podczas generowania nowego hasła")
)
