package postgre

type Config struct {
	URI             string
	Port            string
	User            string
	Password        string
	DBName          string
	ApplicationName string
}

func NewConfig() *Config {
	return &Config{Port: "5432", ApplicationName: "go-pg", User: "postgre", URI: "127.0.0.1"}
}

func (cfg *Config) SetURI(uri string) *Config {
	cfg.URI = uri
	return cfg
}

func (cfg *Config) SetPort(port string) *Config {
	cfg.Port = port
	return cfg
}

func (cfg *Config) SetUser(user string) *Config {
	cfg.User = user
	return cfg
}

func (cfg *Config) SetPassword(password string) *Config {
	cfg.Password = password
	return cfg
}

func (cfg *Config) SetDBName(dbName string) *Config {
	cfg.DBName = dbName
	return cfg
}

func (cfg *Config) SetApplicationName(appName string) *Config {
	cfg.ApplicationName = appName
	return cfg
}
