package email

type Config struct {
	URI      string
	Port     int
	Username string
	Password string
	Address  string
}

func NewConfig() *Config {
	return &Config{
		Port: 587,
	}
}

func (cfg *Config) SetURI(uri string) *Config {
	cfg.URI = uri
	return cfg
}

func (cfg *Config) SetPort(port int) *Config {
	cfg.Port = port
	return cfg
}

func (cfg *Config) SetUsername(uname string) *Config {
	cfg.Username = uname
	return cfg
}

func (cfg *Config) SetPassword(pswd string) *Config {
	cfg.Password = pswd
	return cfg
}

func (cfg *Config) SetAddress(address string) *Config {
	cfg.Address = address
	return cfg
}
