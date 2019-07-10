package repository

import (
	"context"
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/kichiyaki/graphql-starter/backend/postgre"
	"github.com/kichiyaki/graphql-starter/backend/token"
)

var (
	notFoundTokenByValueErrorFormat = "Nie znaleziono tokenu: %s"
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

func (repo *postgreTokenRepo) Get(ctx context.Context, t, value string) (*models.Token, error) {
	token := &models.Token{}
	repo.conn.Model(token).Where("type = ?", t).Where("value = ?", value).First()
	if token.ID > 0 {
		return token, nil
	}
	return nil, fmt.Errorf(notFoundTokenByValueErrorFormat, value)
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
