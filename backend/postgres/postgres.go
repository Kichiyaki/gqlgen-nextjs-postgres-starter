package postgres

import (
	"context"
	"io"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

const (
	functions = `
		CREATE OR REPLACE FUNCTION slugify(id bigint, "value" TEXT)
		RETURNS TEXT AS $$
		-- removes accents (diacritic signs) from a given string --
		WITH "unaccented" AS (
			SELECT unaccent("value") AS "value"
		),
		-- lowercases the string
		"lowercase" AS (
			SELECT lower("value") AS "value"
			FROM "unaccented"
		),
		-- remove single and double quotes
		"removed_quotes" AS (
			SELECT regexp_replace("value", '[''"]+', '', 'gi') AS "value"
			FROM "lowercase"
		),
		-- replaces anything that's not a letter, number, hyphen('-'), or underscore('_') with a hyphen('-')
		"hyphenated" AS (
			SELECT regexp_replace("value", '[^a-z0-9\\-_]+', '-', 'gi') AS "value"
			FROM "removed_quotes"
		),
		-- trims hyphens('-') if they exist on the head or tail of the string
		"trimmed" AS (
			SELECT regexp_replace(regexp_replace("value", '\-+$', ''), '^\-', '') AS "value"
			FROM "hyphenated"
		)
		SELECT id || '-' || "value" as value FROM "trimmed";
		$$ LANGUAGE SQL STRICT IMMUTABLE;
		CREATE OR REPLACE FUNCTION set_slug_from_login() RETURNS trigger
			LANGUAGE plpgsql
			AS $$
		BEGIN
		NEW.slug := slugify(NEW.id, NEW.login);
		RETURN NEW;
		END
		$$;
	`
	extensions = `
		CREATE EXTENSION IF NOT EXISTS "unaccent";
	`
	triggers = `	  
		DROP TRIGGER IF EXISTS set_slug_user ON users;
		CREATE TRIGGER set_slug_user
		BEFORE INSERT ON users FOR EACH ROW
		WHEN (NEW.login IS NOT NULL AND NEW.slug IS NULL) EXECUTE PROCEDURE set_slug_from_login();
	`
)

type DB interface {
	Begin() (*pg.Tx, error)
	Close() error
	Context() context.Context
	CopyFrom(r io.Reader, query interface{}, params ...interface{}) (res pg.Result, err error)
	CopyTo(w io.Writer, query interface{}, params ...interface{}) (res pg.Result, err error)
	CreateTable(model interface{}, opt *orm.CreateTableOptions) error
	Delete(model interface{}) error
	DropTable(model interface{}, opt *orm.DropTableOptions) error
	Exec(query interface{}, params ...interface{}) (pg.Result, error)
	ExecContext(c context.Context, query interface{}, params ...interface{}) (pg.Result, error)
	ExecOne(query interface{}, params ...interface{}) (pg.Result, error)
	ExecOneContext(c context.Context, query interface{}, params ...interface{}) (pg.Result, error)
	ForceDelete(model interface{}) error
	Formatter() orm.QueryFormatter
	Insert(model ...interface{}) error
	Model(model ...interface{}) *orm.Query
	ModelContext(c context.Context, model ...interface{}) *orm.Query
	Prepare(q string) (*pg.Stmt, error)
	Query(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)
	QueryContext(c context.Context, model interface{}, query interface{}, params ...interface{}) (pg.Result, error)
	QueryOne(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)
	QueryOneContext(c context.Context, model interface{}, query interface{}, params ...interface{}) (pg.Result, error)
	RunInTransaction(fn func(*pg.Tx) error) error
	Select(model interface{}) error
	Update(model interface{}) error
}

func LoadFunctionsAndTriggers(db DB) error {
	if _, err := db.Exec(extensions); err != nil {
		return err
	}
	if _, err := db.Exec(functions); err != nil {
		return err
	}
	if _, err := db.Exec(triggers); err != nil {
		return err
	}
	return nil
}
