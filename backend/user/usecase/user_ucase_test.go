package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/kichiyaki/graphql-starter/backend/middleware"

	"github.com/kichiyaki/graphql-starter/backend/models"

	"github.com/kichiyaki/graphql-starter/backend/seed"
	"github.com/kichiyaki/graphql-starter/backend/user/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockUser := &seed.Users()[0]
	role := models.AdministrativeRole
	mockInput := models.UserInput{
		Login:    &mockUser.Login,
		Password: &mockUser.Password,
		Role:     &role,
		Email:    &mockUser.Email,
	}

	t.Run("login is occupied ", func(t *testing.T) {
		mockUserRepo.
			On("Store",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(fmt.Errorf("duplicate key value violates unique login")).
			Once()
		usecase := NewUserUsecase(mockUserRepo)
		_, err := usecase.Store(context.TODO(), mockInput)
		require.Equal(t, "Podany login jest zajęty", err.Error())
	})

	t.Run("email is occupied ", func(t *testing.T) {
		mockUserRepo.
			On("Store",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(fmt.Errorf("duplicate key value violates unique email")).
			Once()
		usecase := NewUserUsecase(mockUserRepo)
		_, err := usecase.Store(context.TODO(), mockInput)
		require.Equal(t, "Podany email jest zajęty", err.Error())
	})

	t.Run("success", func(t *testing.T) {
		mockUserRepo.
			On("Store",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(nil).
			Once()
		usecase := NewUserUsecase(mockUserRepo)
		user, err := usecase.Store(context.TODO(), mockInput)
		require.Equal(t, nil, err)
		require.Equal(t, mockUser.Login, user.Login)
		require.Equal(t, mockUser.Email, user.Email)
	})
}

func TestFetch(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
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
			Return(nil, fmt.Errorf("Some error")).
			Once()
		usecase := NewUserUsecase(mockUserRepo)
		_, err := usecase.Fetch(context.TODO(), p, f)
		require.Equal(t, "Nie udało się wygenerować listy użytkowników", err.Error())
	})

	t.Run("success", func(t *testing.T) {
		mockUserRepo.
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("pgpagination.Pagination"),
				mock.AnythingOfType("*pgfilter.Filter")).
			Return(list, nil).
			Once()
		usecase := NewUserUsecase(mockUserRepo)
		l, err := usecase.Fetch(context.TODO(), p, f)
		require.Equal(t, nil, err)
		require.Equal(t, list.Total, l.Total)
	})
}

func TestUpdate(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
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

	t.Run("error - user not found", func(t *testing.T) {
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(nil, fmt.Errorf("User not found")).
			Once()
		usecase := NewUserUsecase(mockUserRepo)
		_, err := usecase.Update(context.TODO(), id, input)
		require.Equal(t, fmt.Sprintf("Nie znaleziono użytkownika o ID: %d", id), err.Error())
	})

	t.Run("error - nothing to validate", func(t *testing.T) {
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[0], nil).
			Once()
		usecase := NewUserUsecase(mockUserRepo)
		_, err := usecase.Update(context.TODO(), id, models.UserInput{})
		require.Equal(t, "Nie wprowadziłeś żadnych zmian w konfiguracji użytkownika", err.Error())
	})

	t.Run("error - something went wrong with update", func(t *testing.T) {
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
			Return(fmt.Errorf("Something went wrong :(")).
			Once()

		usecase := NewUserUsecase(mockUserRepo)
		_, err := usecase.Update(context.TODO(), id, input)
		require.Equal(t, "Nie udało się zaaktualizować użytkownika", err.Error())
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
			Return(fmt.Errorf("duplicate key value violates unique login")).
			Once()
		usecase := NewUserUsecase(mockUserRepo)
		_, err := usecase.Update(context.TODO(), id, input)
		require.Equal(t, "Podany login jest zajęty", err.Error())
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
			Return(fmt.Errorf("duplicate key value violates unique email")).
			Once()
		usecase := NewUserUsecase(mockUserRepo)
		_, err := usecase.Update(context.TODO(), id, input)
		require.Equal(t, "Podany email jest zajęty", err.Error())
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

		usecase := NewUserUsecase(mockUserRepo)
		u, err := usecase.Update(context.TODO(), id, input)
		require.Equal(t, nil, err)
		require.Equal(t, id, u.ID)
	})
}

func TestDelete(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	users := seed.Users()
	deletedUsers := []*models.User{&users[1]}
	ids := []int{users[1].ID}

	t.Run("ids cannot contains current user ID", func(t *testing.T) {
		ctx := middleware.StoreUserInContext(context.Background(), &users[1])
		usecase := NewUserUsecase(mockUserRepo)
		_, err := usecase.Delete(ctx, ids)
		require.Equal(t, "Nie możesz sam usunąć swojego konta", err.Error())
	})

	t.Run("cannot delete users", func(t *testing.T) {
		ctx := middleware.StoreUserInContext(context.Background(), &users[0])
		mockUserRepo.
			On("Delete",
				mock.Anything,
				ids).
			Return(nil, fmt.Errorf("Something went wrong :(")).
			Once()
		usecase := NewUserUsecase(mockUserRepo)
		_, err := usecase.Delete(ctx, ids)
		require.Equal(t, "Nie udało się usunąć kont użytkowników", err.Error())
	})

	t.Run("success", func(t *testing.T) {
		ctx := middleware.StoreUserInContext(context.Background(), &users[0])
		mockUserRepo.
			On("Delete",
				mock.Anything,
				ids).
			Return(deletedUsers, nil).
			Once()
		usecase := NewUserUsecase(mockUserRepo)
		us, err := usecase.Delete(ctx, ids)
		require.Equal(t, nil, err)
		require.Equal(t, len(deletedUsers), len(us))
	})
}
