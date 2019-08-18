package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/kichiyaki/graphql-starter/backend/middleware"
	"github.com/kichiyaki/graphql-starter/backend/utils"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/gosimple/slug"
	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/kichiyaki/graphql-starter/backend/postgre"
	"github.com/kichiyaki/graphql-starter/backend/user"
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
	localizer, _ := middleware.LocalizerFromContext(ctx)
	err := repo.conn.Insert(u)
	if err != nil {
		if strings.Contains(err.Error(), postgre.DuplicateKeyValueMsg) {
			if strings.Contains(err.Error(), "login") {
				return utils.GetErrorMsg(localizer, "ErrLoginIsOccupied")
			} else if strings.Contains(err.Error(), "email") {
				return utils.GetErrorMsg(localizer, "ErrEmailIsOccupied")
			}
		}
		return utils.GetErrorMsg(localizer, "ErrUserCannotBeCreated")
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
	localizer, _ := middleware.LocalizerFromContext(ctx)
	u := &models.User{}
	repo.conn.Model(u).Where("id = ?", id).First()
	if u.ID > 0 {
		return u, nil
	}
	return nil, utils.GetErrorMsgWithData(localizer, "ErrNotFoundUserByID", map[string]interface{}{
		"ID": id,
	})
}

func (repo *postgreUserRepo) GetBySlug(ctx context.Context, slug string) (*models.User, error) {
	localizer, _ := middleware.LocalizerFromContext(ctx)
	u := &models.User{}
	repo.conn.Model(u).Where("slug = ?", slug).First()
	if u.ID > 0 {
		return u, nil
	}
	return nil, utils.GetErrorMsgWithData(localizer, "ErrNotFoundUserBySlug", map[string]interface{}{
		"Slug": slug,
	})
}

func (repo *postgreUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	localizer, _ := middleware.LocalizerFromContext(ctx)
	u := &models.User{}
	repo.conn.Model(u).Where("email = ?", email).First()
	if u.ID > 0 {
		return u, nil
	}
	return nil, utils.GetErrorMsgWithData(localizer, "ErrNotFoundUserByEmail", map[string]interface{}{
		"Email": email,
	})
}

func (repo *postgreUserRepo) GetByCredentials(ctx context.Context, login, password string) (*models.User, error) {
	localizer, _ := middleware.LocalizerFromContext(ctx)
	u := &models.User{}
	repo.
		conn.
		Model(u).
		Where("login = ?", login).
		First()
	if u.ID > 0 {
		err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
		if err != nil {
			return u, utils.GetErrorMsg(localizer, "ErrInvalidLoginOrPassword")
		}
		return u, nil
	}
	return nil, utils.GetErrorMsg(localizer, "ErrInvalidLoginOrPassword")
}

func (repo *postgreUserRepo) Fetch(ctx context.Context,
	p pgpagination.Pagination,
	f *pgfilter.Filter) (*models.List, error) {
	localizer, _ := middleware.LocalizerFromContext(ctx)
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
		return nil, utils.GetErrorMsg(localizer, "ErrListOfUsersCannotBeGenerated")
	}
	list.Items = users
	list.Total = count

	return &list, nil
}

func (repo *postgreUserRepo) Update(ctx context.Context, u *models.User) error {
	localizer, _ := middleware.LocalizerFromContext(ctx)
	_, err := repo.
		conn.
		Model(u).
		WherePK().
		Returning("*").
		Update()
	if err != nil {
		if strings.Contains(err.Error(), postgre.DuplicateKeyValueMsg) {
			if strings.Contains(err.Error(), "login") {
				return utils.GetErrorMsg(localizer, "ErrLoginIsOccupied")
			} else if strings.Contains(err.Error(), "email") {
				return utils.GetErrorMsg(localizer, "ErrEmailIsOccupied")
			}
		}
		return utils.GetErrorMsg(localizer, "ErrUserCannotBeUpdated")
	}
	return nil
}

func (repo *postgreUserRepo) Delete(ctx context.Context, ids []int) ([]*models.User, error) {
	localizer, _ := middleware.LocalizerFromContext(ctx)
	users := []*models.User{}

	_, err := repo.conn.Model(&users).
		Where("id IN (?)", pg.In(ids)).
		Returning("*").
		Delete()
	if err != nil {
		return nil, utils.GetErrorMsg(localizer, "ErrUsersCannotBeDeleted")
	}

	return users, err
}
