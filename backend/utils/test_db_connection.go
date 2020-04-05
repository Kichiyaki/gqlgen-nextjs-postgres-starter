package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/go-pg/pg/v9"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Println(q.FormattedQuery())
	return nil
}

func ConnectToPostgreTestDB(logger bool) *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     os.Getenv("POSTGRE_USER"),
		Password: os.Getenv("POSTGRE_PASSWORD"),
		Database: os.Getenv("POSTGRE_TEST_DATABASE"),
		Addr:     os.Getenv("POSTGRE_ADDR"),
	})
	if logger {
		db.AddQueryHook(dbLogger{})
	}
	return db
}
