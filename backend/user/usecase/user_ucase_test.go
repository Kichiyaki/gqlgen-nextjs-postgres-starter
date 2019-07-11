package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/kichiyaki/graphql-starter/backend/models"

	_authErrors "github.com/kichiyaki/graphql-starter/backend/auth/errors"
	_authMocks "github.com/kichiyaki/graphql-starter/backend/auth/mocks"
	"github.com/kichiyaki/graphql-starter/backend/seed"
	"github.com/kichiyaki/graphql-starter/backend/user/errors"
	"github.com/kichiyaki/graphql-starter/backend/user/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockAuthUcase := new(_authMocks.Usecase)
	mockUser := &seed.Users()[0]
	role := models.AdministrativeRole
	mockInput := models.UserInput{
		Login:    &mockUser.Login,
		Password: &mockUser.Password,
		Role:     &role,
		Email:    &mockUser.Email,
	}

	t.Run("user must be logged in", func(t *testing.T) {
		mockAuthUcase.On("IsLogged", mock.Anything).Return(false).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Store(context.TODO(), mockInput)
		require.Equal(t, _authErrors.ErrNotLoggedIn, err)
	})

	t.Run("user needs administrative privileges", func(t *testing.T) {
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(false).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Store(context.TODO(), mockInput)
		require.Equal(t, _authErrors.ErrUnauthorized, err)
	})

	t.Run("store returns error", func(t *testing.T) {
		mockUserRepo.
			On("Store",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(errors.ErrEmailIsOccupied).
			Once()
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Store(context.TODO(), mockInput)
		require.Equal(t, errors.ErrEmailIsOccupied, err)
	})

	t.Run("success", func(t *testing.T) {
		mockUserRepo.
			On("Store",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(nil).
			Once()
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		user, err := usecase.Store(context.TODO(), mockInput)
		require.Equal(t, nil, err)
		require.Equal(t, mockUser.Login, user.Login)
		require.Equal(t, mockUser.Email, user.Email)
	})
}

func TestFetch(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockAuthUcase := new(_authMocks.Usecase)
	users := seed.Users()
	p := models.Pagination{Page: 1, Limit: 100}
	f := &models.UserFilter{OnlyActivated: true, Login: "asdf"}
	list := &models.List{
		Total: len(users),
		Items: users,
	}

	t.Run("error", func(t *testing.T) {
		mockUserRepo.
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("pgpagination.Pagination"),
				mock.AnythingOfType("*pgfilter.Filter")).
			Return(nil, errors.ErrListOfUsersCannotBeGenerated).
			Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Fetch(context.Background(), p, f)
		require.Equal(t, errors.ErrListOfUsersCannotBeGenerated, err)
	})

	t.Run("success", func(t *testing.T) {
		mockUserRepo.
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("pgpagination.Pagination"),
				mock.AnythingOfType("*pgfilter.Filter")).
			Return(list, nil).
			Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		l, err := usecase.Fetch(context.Background(), p, f)
		require.Equal(t, nil, err)
		require.Equal(t, list.Total, l.Total)
	})
}

func TestUpdate(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockAuthUcase := new(_authMocks.Usecase)
	users := seed.Users()
	newLogin := "newLoginasd"
	newPassword := "newPassword123"
	newEmail := "newEmail@gmail.com"
	input := models.UserInput{
		Login:    &newLogin,
		Password: &newPassword,
		Email:    &newEmail,
	}
	id := users[0].ID

	t.Run("user must be logged in", func(t *testing.T) {
		mockAuthUcase.On("IsLogged", mock.Anything).Return(false).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Update(context.TODO(), id, input)
		require.Equal(t, _authErrors.ErrNotLoggedIn, err)
	})

	t.Run("user needs administrative privileges", func(t *testing.T) {
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(false).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Update(context.TODO(), id, input)
		require.Equal(t, _authErrors.ErrUnauthorized, err)
	})

	t.Run("error - user not found", func(t *testing.T) {
		e := fmt.Errorf(errors.NotFoundUserByIDErrFormat, id)
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(nil, e).
			Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Update(context.TODO(), id, input)
		require.Equal(t, e, err)
	})

	t.Run("error - nothing to validate", func(t *testing.T) {
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[0], nil).
			Once()
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Update(context.TODO(), id, models.UserInput{})
		require.Equal(t, errors.ErrNothingChanged, err)
	})

	t.Run("error - something went wrong with update", func(t *testing.T) {
		e := fmt.Errorf(errors.UserCannotBeUpdatedErrFormatWithID, users[0].Login)
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[0], nil).
			Once()
		mockUserRepo.
			On("Update",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(e).
			Once()
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()

		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Update(context.TODO(), id, input)
		require.Equal(t, e, err)
	})

	t.Run("login is occupied ", func(t *testing.T) {
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[0], nil).
			Once()
		mockUserRepo.
			On("Update",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(errors.ErrLoginIsOccupied).
			Once()
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Update(context.TODO(), id, input)
		require.Equal(t, errors.ErrLoginIsOccupied, err)
	})

	t.Run("email is occupied ", func(t *testing.T) {
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[0], nil).
			Once()
		mockUserRepo.
			On("Update",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(errors.ErrEmailIsOccupied).
			Once()
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Update(context.TODO(), id, input)
		require.Equal(t, errors.ErrEmailIsOccupied, err)
	})

	t.Run("success", func(t *testing.T) {
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[0], nil).
			Once()
		mockUserRepo.
			On("Update",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(nil).
			Once()
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()

		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		u, err := usecase.Update(context.TODO(), id, input)
		require.Equal(t, nil, err)
		require.Equal(t, id, u.ID)
	})
}

func TestDelete(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockAuthUcase := new(_authMocks.Usecase)
	users := seed.Users()
	deletedUsers := []*models.User{&users[1]}
	ids := []int{users[1].ID}

	t.Run("user must be logged in", func(t *testing.T) {
		mockAuthUcase.On("IsLogged", mock.Anything).Return(false).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Delete(context.TODO(), ids)
		require.Equal(t, _authErrors.ErrNotLoggedIn, err)
	})

	t.Run("user needs administrative privileges", func(t *testing.T) {
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(false).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Delete(context.TODO(), ids)
		require.Equal(t, _authErrors.ErrUnauthorized, err)
	})

	t.Run("ids cannot contains current user ID", func(t *testing.T) {
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()
		mockAuthUcase.On("CurrentUser", mock.Anything).Return(deletedUsers[0]).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Delete(context.TODO(), ids)
		require.Equal(t, errors.ErrUserCannotDeleteHisAccountByHimself, err)
	})

	t.Run("cannot delete users", func(t *testing.T) {
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()
		mockAuthUcase.On("CurrentUser", mock.Anything).Return(&users[0]).Once()
		mockUserRepo.
			On("Delete",
				mock.Anything,
				ids).
			Return(nil, errors.ErrUsersCannotBeDeleted).
			Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Delete(context.TODO(), ids)
		require.Equal(t, errors.ErrUsersCannotBeDeleted, err)
	})

	t.Run("success", func(t *testing.T) {
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()
		mockAuthUcase.On("CurrentUser", mock.Anything).Return(&users[0]).Once()
		mockUserRepo.
			On("Delete",
				mock.Anything,
				ids).
			Return(deletedUsers, nil).
			Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		us, err := usecase.Delete(context.TODO(), ids)
		require.Equal(t, nil, err)
		require.Equal(t, len(deletedUsers), len(us))
	})
}
