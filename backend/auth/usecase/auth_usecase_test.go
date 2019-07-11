package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/kichiyaki/graphql-starter/backend/middleware"

	"github.com/google/uuid"
	_authErrors "github.com/kichiyaki/graphql-starter/backend/auth/errors"
	_emailMock "github.com/kichiyaki/graphql-starter/backend/email/mocks"
	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/kichiyaki/graphql-starter/backend/seed"
	_tokenMocks "github.com/kichiyaki/graphql-starter/backend/token/mocks"
	_userErrors "github.com/kichiyaki/graphql-starter/backend/user/errors"
	"github.com/kichiyaki/graphql-starter/backend/user/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSignup(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockTokenRepo := new(_tokenMocks.Repository)
	mockEmail := new(_emailMock.Email)
	mockUser := &seed.Users()[0]
	role := models.AdministrativeRole
	mockInput := models.UserInput{
		Login:    &mockUser.Login,
		Password: &mockUser.Password,
		Role:     &role,
		Email:    &mockUser.Email,
	}

	t.Run("user cannot be logged in", func(t *testing.T) {
		usecase := NewAuthUsecase(mockUserRepo, mockTokenRepo, mockEmail)
		_, err := usecase.Signup(middleware.StoreUserInContext(context.Background(), mockUser), mockInput)
		require.Equal(t, _authErrors.ErrCannotCreateAccountWhileLoggedIn, err)
	})

	t.Run("login is occupied ", func(t *testing.T) {
		mockUserRepo.
			On("Store",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(_userErrors.ErrLoginIsOccupied).
			Once()
		usecase := NewAuthUsecase(mockUserRepo, mockTokenRepo, mockEmail)
		_, err := usecase.Signup(context.TODO(), mockInput)
		require.Equal(t, _userErrors.ErrLoginIsOccupied, err)
	})

	t.Run("email is occupied ", func(t *testing.T) {
		mockUserRepo.
			On("Store",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(_userErrors.ErrEmailIsOccupied).
			Once()
		usecase := NewAuthUsecase(mockUserRepo, mockTokenRepo, mockEmail)
		_, err := usecase.Signup(context.TODO(), mockInput)
		require.Equal(t, _userErrors.ErrEmailIsOccupied, err)
	})

	t.Run("token cannot be created", func(t *testing.T) {
		mockUserRepo.On("Store", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()
		mockTokenRepo.On("Store", mock.Anything, mock.Anything).Return(fmt.Errorf("Error")).Once()
		usecase := NewAuthUsecase(mockUserRepo, mockTokenRepo, mockEmail)
		_, err := usecase.Signup(context.TODO(), mockInput)
		require.Equal(t, _authErrors.ErrActivationTokenCannotBeCreated, err)
	})

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Store", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()
		mockTokenRepo.On("Store", mock.Anything, mock.Anything).Return(nil).Once()
		mockEmail.On("Send", mock.Anything, mock.Anything).Return(nil).Once()
		usecase := NewAuthUsecase(mockUserRepo, mockTokenRepo, mockEmail)
		user, err := usecase.Signup(context.TODO(), mockInput)
		require.Equal(t, nil, err)
		require.Equal(t, mockUser.Login, user.Login)
		require.Equal(t, mockUser.Email, user.Email)
	})
}

func TestLogin(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockTokenRepo := new(_tokenMocks.Repository)
	mockEmail := new(_emailMock.Email)
	mockUser := &seed.Users()[0]

	t.Run("user cannot be logged in", func(t *testing.T) {
		usecase := NewAuthUsecase(mockUserRepo, mockTokenRepo, mockEmail)
		_, err := usecase.Login(middleware.StoreUserInContext(context.Background(), mockUser), mockUser.Login, mockUser.Password)
		require.Equal(t, _authErrors.ErrCannotLoginWhileLoggedIn, err)
	})

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetByCredentials", mock.Anything, mockUser.Login, mockUser.Password).Return(mockUser, nil).Once()

		usecase := NewAuthUsecase(mockUserRepo, mockTokenRepo, mockEmail)
		user, err := usecase.Login(context.Background(), mockUser.Login, mockUser.Password)
		require.Equal(t, nil, err)
		require.Equal(t, user.Login, mockUser.Login)
	})
}

func TestLogout(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockTokenRepo := new(_tokenMocks.Repository)
	mockEmail := new(_emailMock.Email)
	mockUser := &seed.Users()[0]

	t.Run("user cannot be logged out", func(t *testing.T) {
		usecase := NewAuthUsecase(mockUserRepo, mockTokenRepo, mockEmail)
		err := usecase.Logout(context.Background())
		require.Equal(t, _authErrors.ErrNotLoggedIn, err)
	})

	t.Run("success", func(t *testing.T) {
		usecase := NewAuthUsecase(mockUserRepo, mockTokenRepo, mockEmail)
		err := usecase.Logout(middleware.StoreUserInContext(context.Background(), mockUser))
		require.Equal(t, nil, err)
	})
}

func TestActivate(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockTokenRepo := new(_tokenMocks.Repository)
	mockEmail := new(_emailMock.Email)
	users := seed.Users()
	tokens := seed.Tokens()

	t.Run("user account is activated", func(t *testing.T) {
		id := users[0].ID
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[0], nil).
			Once()
		usecase := NewAuthUsecase(mockUserRepo, mockTokenRepo, mockEmail)
		token := uuid.New().String()
		_, err := usecase.Activate(context.TODO(), id, token)
		require.Equal(t, _authErrors.ErrAccountHasBeenActivated, err)
	})

	t.Run("wrong user id", func(t *testing.T) {
		id := users[1].ID
		tok := tokens[0]
		tok.Type = "asdf"
		tok.UserID += 153
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[1], nil).
			Once()
		mockTokenRepo.
			On("Get",
				mock.Anything,
				models.AccountActivationTokenType,
				tok.Value).
			Return(&tok, nil).
			Once()
		usecase := NewAuthUsecase(mockUserRepo, mockTokenRepo, mockEmail)
		_, err := usecase.Activate(context.TODO(), id, tok.Value)
		require.Equal(t, _authErrors.ErrInvalidActivationToken, err)
	})

	t.Run("cannot update user", func(t *testing.T) {
		id := users[1].ID
		tok := tokens[0]
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[1], nil).
			Once()
		mockTokenRepo.
			On("Get",
				mock.Anything,
				models.AccountActivationTokenType,
				tok.Value).
			Return(&tok, nil).
			Once()
		mockUserRepo.
			On("Update",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(fmt.Errorf("error")).
			Once()
		usecase := NewAuthUsecase(mockUserRepo, mockTokenRepo, mockEmail)
		_, err := usecase.Activate(context.TODO(), id, tok.Value)
		require.Equal(t, _authErrors.ErrAccountCannotBeActivated, err)
	})

	t.Run("success", func(t *testing.T) {
		id := users[1].ID
		tok := tokens[0]
		users[1].Activated = false
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[1], nil).
			Once()
		mockTokenRepo.
			On("Get",
				mock.Anything,
				models.AccountActivationTokenType,
				tok.Value).
			Return(&tok, nil).
			Once()
		mockTokenRepo.
			On("DeleteByUserID",
				mock.Anything,
				models.AccountActivationTokenType,
				id).
			Return([]*models.Token{&tok}, nil).
			Once()
		mockUserRepo.
			On("Update",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(nil).
			Once()
		mockTokenRepo.On("Delete", mock.Anything, []int{tok.ID}).Return([]*models.Token{&tok}, nil).Once()
		usecase := NewAuthUsecase(mockUserRepo, mockTokenRepo, mockEmail)
		user, err := usecase.Activate(context.TODO(), id, tok.Value)
		require.Equal(t, nil, err)
		require.Equal(t, true, user.Activated)
	})
}

func TestGenerateNewActivationToken(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockTokenRepo := new(_tokenMocks.Repository)
	mockEmail := new(_emailMock.Email)
	users := seed.Users()

	t.Run("user account is activated", func(t *testing.T) {
		id := users[0].ID
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[0], nil).
			Once()
		usecase := NewAuthUsecase(mockUserRepo, mockTokenRepo, mockEmail)
		err := usecase.GenerateNewActivationToken(context.TODO(), id)
		require.Equal(t, _authErrors.ErrAccountHasBeenActivated, err)
	})

	t.Run("cannot create token", func(t *testing.T) {
		id := users[1].ID
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[1], nil).
			Once()
		mockTokenRepo.
			On("Store",
				mock.Anything,
				mock.AnythingOfType("*models.Token")).
			Return(fmt.Errorf("Error")).
			Once()

		usecase := NewAuthUsecase(mockUserRepo, mockTokenRepo, mockEmail)
		err := usecase.GenerateNewActivationToken(context.TODO(), id)
		require.Equal(t, _authErrors.ErrActivationTokenCannotBeCreated, err)
	})

	t.Run("success", func(t *testing.T) {
		id := users[1].ID
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[1], nil).
			Once()
		mockTokenRepo.
			On("Store",
				mock.Anything,
				mock.AnythingOfType("*models.Token")).
			Return(nil).
			Once()
		mockEmail.
			On("Send", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		usecase := NewAuthUsecase(mockUserRepo, mockTokenRepo, mockEmail)
		err := usecase.GenerateNewActivationToken(context.TODO(), id)
		require.Equal(t, nil, err)
	})
}
