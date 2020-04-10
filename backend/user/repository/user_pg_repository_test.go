package repository

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"backend/user"

	_errors "backend/errors"
	"backend/models"
	"backend/utils"
	"backend/utils/seed"

	"github.com/go-pg/urlstruct"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestPgRepository(t *testing.T) {
	conn := utils.ConnectToPostgreTestDB(false)
	tx, err := conn.Begin()
	defer tx.Rollback()
	defer conn.Close()
	require.Equal(t, nil, err)
	repo, err := NewPostgreUserRepository(tx)
	require.Equal(t, nil, err)
	seedUsers := seed.Users(5)
	err = seedDatabase(repo)
	require.Equal(t, nil, err)

	t.Run("Fetch", func(t *testing.T) {
		t.Run("Without filter", func(t *testing.T) {
			users, err := repo.Fetch(context.Background(), &models.UserFilter{})
			require.Equal(t, nil, err)
			require.Equal(t, len(seedUsers), users.Total)
		})

		t.Run("With filter", func(t *testing.T) {
			f := new(models.UserFilter)
			err := urlstruct.Unmarshal(context.Background(),
				url.Values{
					"login__neq": {seedUsers[0].Login},
					"email__neq": {seedUsers[1].Email},
					"order":      {"id DESC"},
				},
				f)
			require.Equal(t, nil, err)
			data, err := repo.Fetch(context.Background(), f)
			require.Equal(t, nil, err)
			require.Equal(t, len(seedUsers)-2, data.Total)
			users := data.Items
			maxID := 0
			for _, user := range seedUsers {
				if user.ID > maxID && user.Login != seedUsers[0].Login && user.Email != seedUsers[1].Email {
					maxID = user.ID
				}
			}
			require.Equal(t, maxID, users[0].ID)
		})
	})

	t.Run("GetByID", func(t *testing.T) {
		t.Run("User not found in database", func(t *testing.T) {
			_, err := repo.GetByID(context.Background(), seedUsers[len(seedUsers)-1].ID+1)
			require.Equal(t, true, strings.Contains(err.Error(), _errors.ErrUserNotFound))
		})

		t.Run("User found in database", func(t *testing.T) {
			user, err := repo.GetByID(context.Background(), seedUsers[0].ID)
			require.Equal(t, nil, err)
			require.Equal(t, seedUsers[0].Login, user.Login)
			require.Equal(t, seedUsers[0].Email, user.Email)
			require.Equal(t, seedUsers[0].Role, user.Role)
			require.Equal(t, seedUsers[0].Activated, user.Activated)
		})
	})

	t.Run("GetBySlug", func(t *testing.T) {
		t.Run("User not found in database", func(t *testing.T) {
			_, err := repo.GetBySlug(context.Background(), seedUsers[0].Slug+"asdf123")
			require.Equal(t, true, strings.Contains(err.Error(), _errors.ErrUserNotFound))
		})

		t.Run("User found in database", func(t *testing.T) {
			user, err := repo.GetBySlug(context.Background(), seedUsers[0].Slug)
			require.Equal(t, nil, err)
			require.Equal(t, seedUsers[0].Login, user.Login)
			require.Equal(t, seedUsers[0].Email, user.Email)
			require.Equal(t, seedUsers[0].Role, user.Role)
			require.Equal(t, seedUsers[0].Activated, user.Activated)
		})
	})

	t.Run("GetByEmail", func(t *testing.T) {
		t.Run("User not found in database", func(t *testing.T) {
			_, err := repo.GetByEmail(context.Background(), seedUsers[0].Email+"asdf123")
			require.Equal(t, true, strings.Contains(err.Error(), _errors.ErrUserNotFound))
		})

		t.Run("User found in database", func(t *testing.T) {
			user, err := repo.GetByEmail(context.Background(), seedUsers[0].Email)
			require.Equal(t, nil, err)
			require.Equal(t, seedUsers[0].Login, user.Login)
			require.Equal(t, seedUsers[0].Email, user.Email)
			require.Equal(t, seedUsers[0].Role, user.Role)
			require.Equal(t, seedUsers[0].Activated, user.Activated)
		})
	})

	t.Run("GetByCredentials", func(t *testing.T) {
		t.Run("Credentials are invalid", func(t *testing.T) {
			_, err := repo.GetByCredentials(context.Background(), seedUsers[0].Login, seedUsers[0].Password+"asdf")
			require.Equal(t, true, strings.Contains(err.Error(), _errors.ErrInvalidCredentials))
		})

		t.Run("User found in database", func(t *testing.T) {
			user, err := repo.GetByCredentials(context.Background(), seedUsers[0].Login, seedUsers[0].Password)
			require.Equal(t, nil, err)
			require.Equal(t, seedUsers[0].Login, user.Login)
			require.Equal(t, seedUsers[0].Email, user.Email)
			require.Equal(t, seedUsers[0].Role, user.Role)
			require.Equal(t, seedUsers[0].Activated, user.Activated)
		})
	})

	t.Run("Store", func(t *testing.T) {
		newUser := models.User{
			Login:     "NewUser",
			Password:  "Password123",
			Email:     "newUserEmail@gmail.com",
			Role:      models.UserDefaultRole,
			Activated: true,
			Slug:      "NewUser2-slug",
		}

		t.Run("Successfully stored in database", func(t *testing.T) {
			err := repo.Store(context.Background(), &newUser)
			require.Equal(t, nil, err)
			seedUsers = append(seedUsers, newUser)
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("User not found in database", func(t *testing.T) {
			u := seedUsers[0]
			u.ID = seedUsers[len(seedUsers)-1].ID + 1
			u.Login = "123"
			err := repo.Update(context.Background(), &u)
			require.Equal(t, true, strings.Contains(err.Error(), _errors.ErrUserNotFound))
		})

		t.Run("User found in database and successfully updated", func(t *testing.T) {
			seedUsers[0].Login = "123"
			err := repo.Update(context.Background(), &seedUsers[0])
			require.Equal(t, nil, err)
			users, err := repo.Fetch(context.Background(), &models.UserFilter{
				Login: []string{seedUsers[0].Login},
			})
			require.Equal(t, nil, err)
			require.Equal(t, 1, users.Total)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("User found in database", func(t *testing.T) {
			users, err := repo.Delete(context.Background(), &models.UserFilter{
				ID: []int{seedUsers[0].ID},
			})
			require.Equal(t, nil, err)
			require.Equal(t, 1, len(users))
			require.Equal(t, seedUsers[0].Login, users[0].Login)
			require.Equal(t, seedUsers[0].Email, users[0].Email)
			require.Equal(t, seedUsers[0].Role, users[0].Role)
			require.Equal(t, seedUsers[0].Activated, users[0].Activated)
		})
	})
}

func seedDatabase(repo user.Repository) error {
	errgrp, ctx := errgroup.WithContext(context.Background())
	for _, user := range seed.Users(5) {
		uCopy := user
		errgrp.Go(func() error {
			return repo.Store(ctx, &uCopy)
		})
	}
	return errgrp.Wait()
}
