package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	pgfilter "github.com/kichiyaki/pg-filter"
	pgpagination "github.com/kichiyaki/pg-pagination"

	"golang.org/x/sync/errgroup"

	"github.com/go-pg/pg/orm"
	"github.com/gosimple/slug"
	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/kichiyaki/graphql-starter/backend/seed"
	"github.com/kichiyaki/graphql-starter/backend/user"

	"github.com/kichiyaki/graphql-starter/backend/postgre"
	"github.com/stretchr/testify/require"
)

func TestPostgreUserRepository(t *testing.T) {
	conn := newPostgreConn()
	defer conn.Close()
	conn.DropTable((*models.User)(nil), &orm.DropTableOptions{
		IfExists: true,
	})
	repo, err := NewPostgreUserRepository(conn)
	require.Equal(t, nil, err)
	users := seed.Users()

	t.Run("STORE / must create user", func(t *testing.T) {
		user := users[0]
		err = repo.Store(context.TODO(), &user)
		require.Equal(t, nil, err)
		_, err = repo.GetByID(context.TODO(), user.ID)
		require.Equal(t, nil, err)
		err = clearUsersTable(conn)
		require.Equal(t, nil, err)
	})

	t.Run("GetByID", func(t *testing.T) {
		err := seedPostgreDB(repo)
		require.Equal(t, nil, err)
		t.Run("should return error if user does not exists", func(t *testing.T) {
			_, err := repo.GetByID(context.TODO(), 1)
			require.Equal(t, fmt.Errorf(notFoundUserByIDErrorFormat, 1), err)
		})

		t.Run("should return user", func(t *testing.T) {
			user, err := repo.GetByID(context.TODO(), users[0].ID)
			require.Equal(t, nil, err)
			require.Equal(t, users[0].Login, user.Login)
		})
	})

	t.Run("GetBySlug", func(t *testing.T) {
		t.Run("should return error if user does not exists", func(t *testing.T) {
			_, err := repo.GetBySlug(context.TODO(), "asdf")
			require.Equal(t, fmt.Errorf(notFoundUserBySlugErrorFormat, "asdf"), err)
		})

		t.Run("should return user", func(t *testing.T) {
			user, err := repo.GetBySlug(context.TODO(), slug.MakeLang(fmt.Sprintf("%d-%s", users[0].ID, users[0].Login), "pl"))
			require.Equal(t, nil, err)
			require.Equal(t, users[0].Login, user.Login)
		})
	})

	t.Run("GetByEmail", func(t *testing.T) {
		t.Run("should return error if user does not exists", func(t *testing.T) {
			email := "elówa@gmailo.com"
			_, err := repo.GetByEmail(context.TODO(), email)
			require.Equal(t, fmt.Errorf(notFoundUserByEmailErrorFormat, email), err)
		})

		t.Run("should return user", func(t *testing.T) {
			user, err := repo.GetByEmail(context.TODO(), users[0].Email)
			require.Equal(t, nil, err)
			require.Equal(t, users[0].Login, user.Login)
		})
	})

	t.Run("GetByCredentials", func(t *testing.T) {
		t.Run("should return error if password is invalid", func(t *testing.T) {
			_, err := repo.GetByCredentials(context.TODO(), users[0].Login+"asdasda", users[0].Password)
			require.Equal(t, errInvalidLoginOrPassword, err)
		})

		t.Run("should return error if password is invalid", func(t *testing.T) {
			_, err := repo.GetByCredentials(context.TODO(), users[0].Login, "asdasdadsa")
			require.Equal(t, errInvalidLoginOrPassword, err)
		})

		t.Run("should return user", func(t *testing.T) {
			user, err := repo.GetByCredentials(context.TODO(), users[0].Login, users[0].Password)
			require.Equal(t, nil, err)
			require.Equal(t, users[0].Login, user.Login)
		})
	})

	t.Run("Fetch", func(t *testing.T) {
		t.Run("Without filter", func(t *testing.T) {
			p := pgpagination.Pagination{Limit: 100, Page: 1}
			list, err := repo.Fetch(context.TODO(), p, nil)
			require.Equal(t, nil, err)
			us, ok := list.Items.([]*models.User)
			require.Equal(t, true, ok)
			require.Equal(t, len(users), len(us))
		})

		t.Run("With filter", func(t *testing.T) {
			m := make(map[string]string)
			m["activated"] = "eq__true"
			m["login"] = "ieq__" + users[0].Login + "%"
			f := pgfilter.New(m)
			p := pgpagination.Pagination{Limit: 100, Page: 1}
			list, err := repo.Fetch(context.TODO(), p, f)
			require.Equal(t, nil, err)
			us, ok := list.Items.([]*models.User)
			require.Equal(t, true, ok)
			require.Equal(t, 1, len(us))
		})
	})

	t.Run("Update", func(t *testing.T) {
		u := users[0]
		u.Login = "NewLoginKiszk"
		err := repo.Update(context.TODO(), &users[0])
		require.Equal(t, nil, err)
		err = clearUsersTable(conn)
		require.Equal(t, nil, err)
	})

	t.Run("Update", func(t *testing.T) {
		err := seedPostgreDB(repo)
		require.Equal(t, nil, err)
		ids := []int{users[0].ID}
		us, err := repo.Delete(context.TODO(), ids)
		require.Equal(t, nil, err)
		require.Equal(t, len(ids), len(us))
	})
}

func clearUsersTable(conn *postgre.Database) error {
	var users []models.User
	_, err := conn.Query(&users, `SELECT * FROM users`)
	if err != nil {
		return err
	}

	if len(users) > 0 {
		_, err = conn.Model(&users).Delete()
		if err != nil {
			return err
		}
	}
	return nil
}

func seedPostgreDB(repo user.Repository) error {
	errgrp, ctx := errgroup.WithContext(context.Background())
	for _, user := range seed.Users() {
		u := user
		errgrp.Go(func() error {
			return repo.Store(ctx, &u)
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
