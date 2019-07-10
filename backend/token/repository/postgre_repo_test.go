package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-pg/pg/orm"
	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/kichiyaki/graphql-starter/backend/postgre"
	"github.com/kichiyaki/graphql-starter/backend/seed"
	"github.com/kichiyaki/graphql-starter/backend/token"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestPostgreRepo(t *testing.T) {
	conn := newPostgreConn()
	defer conn.Close()
	conn.DropTable((*models.Token)(nil), &orm.DropTableOptions{
		IfExists: true,
	})
	repo, err := NewPostgreTokenRepository(conn)
	require.Equal(t, nil, err)
	tokens := seed.Tokens()

	t.Run("STORE", func(t *testing.T) {
		err := repo.Store(context.Background(), &tokens[0])
		require.Equal(t, nil, err)
		token, err := repo.Get(context.Background(), tokens[0].Type, tokens[0].Value)
		require.Equal(t, nil, err)
		require.Equal(t, tokens[0].ID, token.ID)
		clearTokensTable(conn)
	})

	t.Run("Get", func(t *testing.T) {
		seedPostgreDB(repo)
		t.Run("invalid token value", func(t *testing.T) {
			value := tokens[0].Value + "asdasda"
			_, err := repo.Get(context.Background(), tokens[0].Type, value)
			require.Equal(t, fmt.Errorf(notFoundTokenByValueErrorFormat, value), err)
		})
		t.Run("invalid type", func(t *testing.T) {
			_, err := repo.Get(context.Background(), tokens[0].Type+"asdasda", tokens[0].Value)
			require.Equal(t, fmt.Errorf(notFoundTokenByValueErrorFormat, tokens[0].Value), err)
		})
		t.Run("success", func(t *testing.T) {
			token, err := repo.Get(context.Background(), tokens[0].Type, tokens[0].Value)
			require.Equal(t, nil, err)
			require.Equal(t, tokens[0].ID, token.ID)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		ids := []int{tokens[0].ID}
		deletedTokens, err := repo.Delete(context.Background(), ids)
		require.Equal(t, nil, err)
		require.Equal(t, len(ids), len(deletedTokens))
		clearTokensTable(conn)
	})

	t.Run("DeleteByUserID", func(t *testing.T) {
		seedPostgreDB(repo)
		matchTokens := []models.Token{}
		for _, token := range tokens {
			if token.UserID == tokens[0].UserID && token.Type == tokens[0].Type {
				matchTokens = append(matchTokens, token)
			}
		}

		deletedTokens, err := repo.DeleteByUserID(context.Background(), tokens[0].Type, tokens[0].UserID)
		require.Equal(t, nil, err)
		require.Equal(t, len(matchTokens), len(deletedTokens))
		clearTokensTable(conn)
	})
}

func clearTokensTable(conn *postgre.Database) error {
	var tokens []models.Token
	_, err := conn.Query(&tokens, `SELECT * FROM tokens`)
	if err != nil {
		return err
	}

	if len(tokens) > 0 {
		_, err = conn.Model(&tokens).Delete()
		if err != nil {
			return err
		}
	}
	return nil
}

func seedPostgreDB(repo token.Repository) error {
	errgrp, ctx := errgroup.WithContext(context.Background())
	for _, token := range seed.Tokens() {
		t := token
		errgrp.Go(func() error {
			return repo.Store(ctx, &t)
		})
	}
	return errgrp.Wait()
}

func newPostgreConn() *postgre.Database {
	cfg := postgre.
		NewConfig().
		SetApplicationName("postgre-tests").
		SetDBName("tribalwars-api-tests").
		SetURI(os.Getenv("POSTGRE_URI")).
		SetUser(os.Getenv("POSTGRE_USER")).
		SetPassword(os.Getenv("POSTGRE_PASSWORD"))
	conn, _ := postgre.NewDatabase(cfg)
	return conn
}
