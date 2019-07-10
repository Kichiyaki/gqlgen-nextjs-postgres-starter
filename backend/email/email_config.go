package email

type EmailConfig struct {
	From    string
	To      []string
	Subject string
	Body    string
}

func NewEmailConfig() *EmailConfig {
	return &EmailConfig{}
}

func (cfg *EmailConfig) SetFrom(from string) *EmailConfig {
	cfg.From = from
	return cfg
}

func (cfg *EmailConfig) SetTo(to []string) *EmailConfig {
	cfg.To = to
	return cfg
}

func (cfg *EmailConfig) SetSubject(subject string) *EmailConfig {
	cfg.Subject = subject
	return cfg
}

func (cfg *EmailConfig) SetBody(body string) *EmailConfig {
	cfg.Body = body
	return cfg
}
