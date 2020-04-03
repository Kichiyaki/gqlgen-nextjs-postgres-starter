package repository

import (
	"backend/user"
	"context"
	"strings"

	"github.com/sirupsen/logrus"

	_errors "backend/errors"
	"backend/middleware"
	"backend/models"
	"backend/postgres"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

type postgreRepository struct {
	postgres.DB
	logrus *logrus.Entry
}

func NewPostgreUserRepository(conn postgres.DB) (user.Repository, error) {
	log := logrus.WithField("package", "user/repository")
	if err := conn.CreateTable((*models.User)(nil), &orm.CreateTableOptions{
		IfNotExists: true,
	}); err != nil {
		log.Debugf("Cannot create user table: %s", err.Error())
		return nil, err
	}
	return &postgreRepository{conn,
		log,
	}, nil
}

func (repo *postgreRepository) Fetch(ctx context.Context, f *models.UserFilter) (models.UserList, error) {
	var err error
	users := []*models.User{}
	pagination := models.UserList{}
	query := repo.Model(&users)
	log := repo.logrus.WithField("filter", f)
	log.Debug("Fetch")

	if f != nil {
		query = query.
			WhereStruct(f).
			Limit(f.Limit).
			Offset(f.Offset)

		if len(f.Order) > 0 {
			query = query.Order(f.Order...)
		}

		if f.Activated == "true" {
			query = query.Where("activated = true")
		} else if f.Activated == "false" {
			query = query.Where("activated = false")
		}
	}

	if pagination.Total, err = query.
		SelectAndCount(); err != nil && err != pg.ErrNoRows {
		log.Debugf("Fetch err: %s", err.Error())
		return pagination, _errors.Wrap(_errors.ErrInternalServerError, err)
	}
	pagination.Items = users

	return pagination, nil
}

func (repo *postgreRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{
		ID: id,
	}
	log := repo.logrus.WithField("id", id)
	log.Debug("GetByID")
	if err := repo.Select(user); err != nil {
		log.Debugf("GetByID err: %s", err.Error())
		return user, _errors.Wrap(_errors.ErrUserNotFound, err)
	}
	return user, nil
}

func (repo *postgreRepository) GetBySlug(ctx context.Context, slug string) (*models.User, error) {
	user := &models.User{}
	log := repo.logrus.WithField("slug", slug)
	log.Debug("GetBySlug")
	if err := repo.
		Model(user).
		Where("slug = ?", slug).
		Limit(1).
		Select(); err != nil {
		log.Debugf("GetBySlug err: %s", err.Error())
		return user, _errors.Wrap(_errors.ErrUserNotFound, err)
	}
	return user, nil
}

func (repo *postgreRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	log := repo.logrus.WithField("email", email)
	log.Debug("GetByEmail")
	if err := repo.
		Model(user).
		Where("email = ?", email).
		Limit(1).
		Select(); err != nil {
		log.Debugf("GetByEmail err: %s", err.Error())
		return user, _errors.Wrap(_errors.ErrUserNotFound, err)
	}
	return user, nil
}

//we assume that the password is hashed
func (repo *postgreRepository) GetByCredentials(ctx context.Context, login, password string) (*models.User, error) {
	u := &models.User{}
	log := repo.logrus.WithField("login", login).WithField("password", password)
	log.Debug("GetByCredentials")
	if err := repo.
		Model(u).
		Where("login = ?", login).
		Limit(1).
		Select(); err != nil {
		log.Debugf("GetByCredentials err: %s", err.Error())
		return u, _errors.Wrap(_errors.ErrInvalidCredentials, err)
	}
	if err := u.CompareHashAndPassword(password); err != nil {
		log.Debugf("GetByCredentials err: %s", err.Error())
		return nil, err
	}
	return u, nil
}

func (repo *postgreRepository) Update(ctx context.Context, u *models.User) error {
	log := repo.logrus.WithField("user", u)
	log.Debug("Update")
	if _, err := repo.
		Model(u).
		WherePK().
		Apply(middleware.UpdateNotZero).
		Returning("*").
		Update(); err != nil {
		log.Debugf("Update err: %s", err.Error())
		if err == pg.ErrNoRows {
			return _errors.Wrap(_errors.ErrUserNotFound, err)
		}

		if strings.Contains(err.Error(), "login") {
			return _errors.Wrap(_errors.ErrLoginMustBeUnique, err)
		} else if strings.Contains(err.Error(), "email") {
			return _errors.Wrap(_errors.ErrEmailMustBeUnique, err)
		}

		return _errors.Wrap(_errors.ErrInternalServerError, err)
	}
	return nil
}

func (repo *postgreRepository) Store(ctx context.Context, u *models.User) error {
	log := repo.logrus.WithField("user", u)
	log.Debug("Store")
	if _, err := repo.Model(u).Insert(); err != nil {
		log.Debugf("Store err: %s", err.Error())
		if strings.Contains(err.Error(), "login") {
			return _errors.Wrap(_errors.ErrLoginMustBeUnique, err)
		} else if strings.Contains(err.Error(), "email") {
			return _errors.Wrap(_errors.ErrEmailMustBeUnique, err)
		}

		return _errors.Wrap(_errors.ErrInternalServerError, err)
	}
	return nil
}

func (repo *postgreRepository) Delete(ctx context.Context, f *models.UserFilter) ([]*models.User, error) {
	users := []*models.User{}
	query := repo.Model(&users)
	log := repo.logrus.WithField("filter", f)
	log.Debug("Delete")
	if f != nil {
		query = query.
			WhereStruct(f)

		if f.Activated == "true" {
			query = query.Where("activated = true")
		} else if f.Activated == "false" {
			query = query.Where("activated = false")
		}
	}
	_, err := query.
		Returning("*").
		Delete()
	if err != nil && err != pg.ErrNoRows {
		log.Debugf("Delete err: %s", err.Error())
		return nil, _errors.Wrap(_errors.ErrInternalServerError, err)
	}
	return users, err
}
