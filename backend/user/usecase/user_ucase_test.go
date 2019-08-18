package usecase

import (
	"context"
	"testing"

	"github.com/kichiyaki/graphql-starter/backend/middleware"
	"github.com/kichiyaki/graphql-starter/backend/utils"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/kichiyaki/graphql-starter/backend/models"

	_authMocks "github.com/kichiyaki/graphql-starter/backend/auth/mocks"
	"github.com/kichiyaki/graphql-starter/backend/seed"
	"github.com/kichiyaki/graphql-starter/backend/user/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var localizer *i18n.Localizer

func init() {
	localizer = utils.GetLocalizer(language.Polish, "../../i18n/locales/active.pl.json")
}

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
		_, err := usecase.Store(getContext(), mockInput)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrNotLoggedIn"), err)
	})

	t.Run("user needs administrative privileges", func(t *testing.T) {
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(false).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Store(getContext(), mockInput)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrUnauthorized"), err)
	})

	t.Run("store returns error", func(t *testing.T) {
		returnedErr := utils.GetErrorMsg(localizer, "ErrEmailIsOccupied")
		mockUserRepo.
			On("Store",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(returnedErr).
			Once()
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Store(getContext(), mockInput)
		require.Equal(t, returnedErr, err)
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
		user, err := usecase.Store(getContext(), mockInput)
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
		returnedErr := utils.GetErrorMsg(localizer, "ErrListOfUsersCannotBeGenerated")
		mockUserRepo.
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("pgpagination.Pagination"),
				mock.AnythingOfType("*pgfilter.Filter")).
			Return(nil, returnedErr).
			Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Fetch(getContext(), p, f)
		require.Equal(t, returnedErr, err)
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
		l, err := usecase.Fetch(getContext(), p, f)
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
		_, err := usecase.Update(getContext(), id, input)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrNotLoggedIn"), err)
	})

	t.Run("user needs administrative privileges", func(t *testing.T) {
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(false).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Update(getContext(), id, input)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrUnauthorized"), err)
	})

	t.Run("error - user not found", func(t *testing.T) {
		returnedErr := utils.GetErrorMsgWithData(localizer, "ErrNotFoundUserByID", map[string]interface{}{
			"ID": id,
		})
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(nil, returnedErr).
			Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Update(getContext(), id, input)
		require.Equal(t, returnedErr, err)
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
		_, err := usecase.Update(getContext(), id, models.UserInput{})
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrNothingChanged"), err)
	})

	t.Run("error - something went wrong with update", func(t *testing.T) {
		returnedErr := utils.GetErrorMsg(localizer, "ErrUserCannotBeUpdated")
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
			Return(returnedErr).
			Once()
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()

		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Update(getContext(), id, input)
		require.Equal(t, returnedErr, err)
	})

	t.Run("login is occupied ", func(t *testing.T) {
		returnedErr := utils.GetErrorMsg(localizer, "ErrLoginIsOccupied")
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
			Return(returnedErr).
			Once()
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Update(getContext(), id, input)
		require.Equal(t, returnedErr, err)
	})

	t.Run("email is occupied ", func(t *testing.T) {
		returnedErr := utils.GetErrorMsg(localizer, "ErrEmailIsOccupied")
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
			Return(returnedErr).
			Once()
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Update(getContext(), id, input)
		require.Equal(t, returnedErr, err)
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
		u, err := usecase.Update(getContext(), id, input)
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
		_, err := usecase.Delete(getContext(), ids)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrNotLoggedIn"), err)
	})

	t.Run("user needs administrative privileges", func(t *testing.T) {
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(false).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Delete(getContext(), ids)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrUnauthorized"), err)
	})

	t.Run("ids cannot contains current user ID", func(t *testing.T) {
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()
		mockAuthUcase.On("CurrentUser", mock.Anything).Return(deletedUsers[0]).Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Delete(getContext(), ids)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrUserCannotDeleteHisAccountByHimself"), err)
	})

	t.Run("cannot delete users", func(t *testing.T) {
		returnedErr := utils.GetErrorMsg(localizer, "ErrUsersCannotBeDeleted")
		mockAuthUcase.On("IsLogged", mock.Anything).Return(true).Once()
		mockAuthUcase.On("HasAdministrativePrivileges", mock.Anything).Return(true).Once()
		mockAuthUcase.On("CurrentUser", mock.Anything).Return(&users[0]).Once()
		mockUserRepo.
			On("Delete",
				mock.Anything,
				ids).
			Return(nil, returnedErr).
			Once()
		usecase := NewUserUsecase(mockUserRepo, mockAuthUcase)
		_, err := usecase.Delete(getContext(), ids)
		require.Equal(t, returnedErr, err)
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
		us, err := usecase.Delete(getContext(), ids)
		require.Equal(t, nil, err)
		require.Equal(t, len(deletedUsers), len(us))
	})
}

func getContext() context.Context {
	return middleware.StoreLocalizerInContext(context.Background(), localizer)
}
