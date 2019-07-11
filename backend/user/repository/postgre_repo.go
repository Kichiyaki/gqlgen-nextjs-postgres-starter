package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/gosimple/slug"
	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/kichiyaki/graphql-starter/backend/postgre"
	"github.com/kichiyaki/graphql-starter/backend/user"
	"github.com/kichiyaki/graphql-starter/backend/user/errors"
	pgfilter "github.com/kichiyaki/pg-filter"
	pgpagination "github.com/kichiyaki/pg-pagination"
	"golang.org/x/crypto/bcrypt"
)

type postgreUserRepo struct {
	conn *postgre.Database
}

func NewPostgreUserRepository(conn *postgre.Database) (user.Repository, error) {
	return &postgreUserRepo{conn},
		conn.CreateTable((*models.User)(nil), &orm.CreateTableOptions{
			IfNotExists: true,
		})
}

func (repo *postgreUserRepo) Store(ctx context.Context, u *models.User) error {
	err := repo.conn.Insert(u)
	if err != nil {
		if strings.Contains(err.Error(), postgre.DuplicateKeyValueMsg) {
			if strings.Contains(err.Error(), "login") {
				return errors.ErrLoginIsOccupied
			} else if strings.Contains(err.Error(), "email") {
				return errors.ErrEmailIsOccupied
			}
		}
		return errors.ErrUserCannotBeCreated
	}
	s := slug.MakeLang(fmt.Sprintf("%d-%s", u.ID, u.Login), "pl")
	_, err = repo.
		conn.
		Model(u).
		Set("slug = ?", s).
		WherePK().
		Returning("*").
		Update()
	return err
}

func (repo *postgreUserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	u := &models.User{}
	repo.conn.Model(u).Where("id = ?", id).First()
	if u.ID > 0 {
		return u, nil
	}
	return nil, fmt.Errorf(errors.NotFoundUserByIDErrFormat, id)
}

func (repo *postgreUserRepo) GetBySlug(ctx context.Context, slug string) (*models.User, error) {
	u := &models.User{}
	repo.conn.Model(u).Where("slug = ?", slug).First()
	if u.ID > 0 {
		return u, nil
	}
	return nil, fmt.Errorf(errors.NotFoundUserBySlugErrFormat, slug)
}

func (repo *postgreUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	repo.conn.Model(u).Where("email = ?", email).First()
	if u.ID > 0 {
		return u, nil
	}
	return nil, fmt.Errorf(errors.NotFoundUserByEmailErrFormat, email)
}

func (repo *postgreUserRepo) GetByCredentials(ctx context.Context, login, password string) (*models.User, error) {
	u := &models.User{}
	repo.
		conn.
		Model(u).
		Where("login = ?", login).
		First()
	if u.ID > 0 {
		err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
		if err != nil {
			return u, errors.ErrInvalidLoginOrPassword
		}
		return u, nil
	}
	return nil, errors.ErrInvalidLoginOrPassword
}

func (repo *postgreUserRepo) Fetch(ctx context.Context,
	p pgpagination.Pagination,
	f *pgfilter.Filter) (*models.List, error) {
	list := models.List{}
	users := []*models.User{}

	q := repo.conn.Model(&users)

	if p.GetPage() == 0 {
		p.SetPage(1)
	}

	if f != nil {
		q.Apply(f.Filter)
	}

	count, err := q.Apply(p.Paginate).SelectAndCount()
	if err != nil && err != pg.ErrNoRows {
		return nil, errors.ErrListOfUsersCannotBeGenerated
	}
	list.Items = users
	list.Total = count

	return &list, nil
}

func (repo *postgreUserRepo) Update(ctx context.Context, u *models.User) error {
	_, err := repo.
		conn.
		Model(u).
		WherePK().
		Returning("*").
		Update()
	if err != nil {
		if strings.Contains(err.Error(), postgre.DuplicateKeyValueMsg) {
			if strings.Contains(err.Error(), "login") {
				return errors.ErrLoginIsOccupied
			} else if strings.Contains(err.Error(), "email") {
				return errors.ErrEmailIsOccupied
			}
		}
		if u.Login != "" {
			return fmt.Errorf(errors.UserCannotBeUpdatedErrFormatWithLogin, u.Login)
		}
		return fmt.Errorf(errors.UserCannotBeUpdatedErrFormatWithID, u.ID)
	}
	return nil
}

func (repo *postgreUserRepo) Delete(ctx context.Context, ids []int) ([]*models.User, error) {
	users := []*models.User{}

	_, err := repo.conn.Model(&users).
		Where("id IN (?)", pg.In(ids)).
		Returning("*").
		Delete()
	if err != nil {
		return nil, errors.ErrUsersCannotBeDeleted
	}

	return users, err
}
