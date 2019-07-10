package email

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

var (
	NoConfigurationError = fmt.Errorf("Nie wprowadziłeś konfiguracji")
)

type Email interface {
	Send(cfg *EmailConfig) error
}

type email struct {
	dialer *gomail.Dialer
	cfg    *Config
}

func NewEmail(cfg *Config) (Email, error) {
	if cfg == nil {
		return nil, NoConfigurationError
	}
	e := &email{cfg: cfg, dialer: gomail.NewDialer(cfg.URI, cfg.Port, cfg.Username, cfg.Password)}
	return e, nil
}

func (e *email) Send(cfg *EmailConfig) error {
	m := gomail.NewMessage()
	if cfg.From == "" {
		cfg.From = e.cfg.Address
	}
	m.SetHeader("From", cfg.From)
	m.SetHeader("To", cfg.To...)
	m.SetHeader("Subject", cfg.Subject)
	m.SetBody("text/html", cfg.Body)
	return e.dialer.DialAndSend(m)
}
