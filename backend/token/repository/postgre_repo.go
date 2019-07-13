package repository

import (
	"context"
	"fmt"

	pgfilter "github.com/kichiyaki/pg-filter"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/kichiyaki/graphql-starter/backend/postgre"
	"github.com/kichiyaki/graphql-starter/backend/token"
)

var (
	errCannotGenerateSliceOfToken = fmt.Errorf("Nie znaleziono żadnego tokenu")
)

type postgreTokenRepo struct {
	conn *postgre.Database
}

func NewPostgreTokenRepository(conn *postgre.Database) (token.Repository, error) {
	return &postgreTokenRepo{conn},
		conn.CreateTable((*models.Token)(nil), &orm.CreateTableOptions{
			IfNotExists: true,
		})
}

func (repo *postgreTokenRepo) Store(ctx context.Context, t *models.Token) error {
	return repo.conn.Insert(t)
}

func (repo *postgreTokenRepo) Fetch(ctx context.Context, f *pgfilter.Filter) ([]*models.Token, error) {
	tokens := []*models.Token{}
	query := repo.conn.Model(&tokens)
	if f != nil {
		query = query.Apply(f.Filter)
	}
	err := query.Select()
	if err != nil && err != pg.ErrNoRows {
		fmt.Println(err)
		return nil, errCannotGenerateSliceOfToken
	}
	return tokens, nil
}

func (repo *postgreTokenRepo) Delete(ctx context.Context, ids []int) ([]*models.Token, error) {
	tokens := []*models.Token{}

	_, err := repo.conn.Model(&tokens).
		Where("id IN (?)", pg.In(ids)).
		Returning("*").
		Delete()

	return tokens, err
}

func (repo *postgreTokenRepo) DeleteByUserID(ctx context.Context, t string, id int) ([]*models.Token, error) {
	tokens := []*models.Token{}

	_, err := repo.conn.Model(&tokens).
		Where("type = ?", t).
		Where("user_id = ?", id).
		Returning("*").
		Delete()

	return tokens, err
}
