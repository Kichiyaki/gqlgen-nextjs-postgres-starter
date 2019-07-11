package postgre

import (
	"github.com/go-pg/pg"
)

const (
	DuplicateKeyValueMsg = "duplicate key value violates unique"
)

type Database struct {
	*pg.DB
	cfg *Config
}

func NewDatabase(cfg *Config) (*Database, error) {
	if cfg == nil {
		cfg = NewConfig()
	}

	db := &Database{pg.Connect(&pg.Options{
		Addr:            cfg.URI + ":" + cfg.Port,
		User:            cfg.User,
		Password:        cfg.Password,
		Database:        cfg.DBName,
		ApplicationName: cfg.ApplicationName,
	}), cfg}

	return db, nil
}
