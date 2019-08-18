package usecase

import (
	"context"
	"fmt"
	"testing"

	"golang.org/x/text/language"

	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/kichiyaki/graphql-starter/backend/middleware"

	"github.com/google/uuid"
	_emailMock "github.com/kichiyaki/graphql-starter/backend/email/mocks"
	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/kichiyaki/graphql-starter/backend/seed"
	_sessionsMocks "github.com/kichiyaki/graphql-starter/backend/sessions/mocks"
	_tokenMocks "github.com/kichiyaki/graphql-starter/backend/token/mocks"
	"github.com/kichiyaki/graphql-starter/backend/user/mocks"
	"github.com/kichiyaki/graphql-starter/backend/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var localizer *i18n.Localizer

func init() {
	localizer = utils.GetLocalizer(language.Polish, "../../i18n/locales/active.pl.json")
}

func TestSignup(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockTokenRepo := new(_tokenMocks.Repository)
	mockEmail := new(_emailMock.Email)
	mockSessStore := new(_sessionsMocks.Store)
	mockUser := &seed.Users()[0]
	role := models.AdministrativeRole
	mockInput := models.UserInput{
		Login:    &mockUser.Login,
		Password: &mockUser.Password,
		Role:     &role,
		Email:    &mockUser.Email,
	}
	cfg := &Config{
		SessStore:       mockSessStore,
		UserRepo:        mockUserRepo,
		TokenRepo:       mockTokenRepo,
		Email:           mockEmail,
		ApplicationName: "AppName",
		FrontendURL:     "http://localhost:3000",
	}

	t.Run("user cannot be logged in", func(t *testing.T) {
		usecase := NewAuthUsecase(cfg)
		_, err := usecase.Signup(getContext(localizer, mockUser), mockInput)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrCannotCreateAccountWhileLoggedIn"), err)
	})

	t.Run("login is occupied ", func(t *testing.T) {
		errorReturned := utils.GetErrorMsg(localizer, "ErrLoginIsOccupied")
		mockUserRepo.
			On("Store",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(errorReturned).
			Once()
		usecase := NewAuthUsecase(cfg)
		_, err := usecase.Signup(getContext(localizer, nil), mockInput)
		require.Equal(t, errorReturned, err)
	})

	t.Run("email is occupied ", func(t *testing.T) {
		errorReturned := utils.GetErrorMsg(localizer, "ErrEmailIsOccupied")
		mockUserRepo.
			On("Store",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(errorReturned).
			Once()
		usecase := NewAuthUsecase(cfg)
		_, err := usecase.Signup(getContext(localizer, nil), mockInput)
		require.Equal(t, errorReturned, err)
	})

	t.Run("token cannot be created", func(t *testing.T) {
		mockUserRepo.On("Store", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()
		mockTokenRepo.On("Store", mock.Anything, mock.Anything).Return(fmt.Errorf("Error")).Once()
		usecase := NewAuthUsecase(cfg)
		_, err := usecase.Signup(getContext(localizer, nil), mockInput)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrActivationTokenCannotBeCreated"), err)
	})

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Store", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()
		mockTokenRepo.On("Store", mock.Anything, mock.Anything).Return(nil).Once()
		mockEmail.On("Send", mock.Anything, mock.Anything).Return(nil).Once()
		usecase := NewAuthUsecase(cfg)
		user, err := usecase.Signup(getContext(localizer, nil), mockInput)
		require.Equal(t, nil, err)
		require.Equal(t, mockUser.Login, user.Login)
		require.Equal(t, mockUser.Email, user.Email)
	})
}

func TestLogin(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockTokenRepo := new(_tokenMocks.Repository)
	mockEmail := new(_emailMock.Email)
	mockSessStore := new(_sessionsMocks.Store)
	mockUser := &seed.Users()[0]
	cfg := &Config{
		SessStore:       mockSessStore,
		UserRepo:        mockUserRepo,
		TokenRepo:       mockTokenRepo,
		Email:           mockEmail,
		ApplicationName: "AppName",
		FrontendURL:     "http://localhost:3000",
	}

	t.Run("user cannot be logged in", func(t *testing.T) {
		usecase := NewAuthUsecase(cfg)
		_, err := usecase.Login(getContext(localizer, mockUser), mockUser.Login, mockUser.Password)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrCannotLoginWhileLoggedIn"), err)
	})

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetByCredentials", mock.Anything, mockUser.Login, mockUser.Password).Return(mockUser, nil).Once()

		usecase := NewAuthUsecase(cfg)
		user, err := usecase.Login(getContext(localizer, nil), mockUser.Login, mockUser.Password)
		require.Equal(t, nil, err)
		require.Equal(t, user.Login, mockUser.Login)
	})
}

func TestLogout(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockTokenRepo := new(_tokenMocks.Repository)
	mockEmail := new(_emailMock.Email)
	mockSessStore := new(_sessionsMocks.Store)
	mockUser := &seed.Users()[0]
	cfg := &Config{
		SessStore:       mockSessStore,
		UserRepo:        mockUserRepo,
		TokenRepo:       mockTokenRepo,
		Email:           mockEmail,
		ApplicationName: "AppName",
		FrontendURL:     "http://localhost:3000",
	}

	t.Run("user cannot be logged out", func(t *testing.T) {
		usecase := NewAuthUsecase(cfg)
		err := usecase.Logout(getContext(localizer, nil))
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrNotLoggedIn"), err)
	})

	t.Run("success", func(t *testing.T) {
		usecase := NewAuthUsecase(cfg)
		err := usecase.Logout(getContext(localizer, mockUser))
		require.Equal(t, nil, err)
	})
}

func TestGenerateNewActivationToken(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockTokenRepo := new(_tokenMocks.Repository)
	mockEmail := new(_emailMock.Email)
	mockSessStore := new(_sessionsMocks.Store)
	users := seed.Users()
	cfg := &Config{
		SessStore:       mockSessStore,
		UserRepo:        mockUserRepo,
		TokenRepo:       mockTokenRepo,
		Email:           mockEmail,
		ApplicationName: "AppName",
		FrontendURL:     "http://localhost:3000",
	}

	t.Run("user account is activated", func(t *testing.T) {
		id := users[0].ID
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[0], nil).
			Once()
		usecase := NewAuthUsecase(cfg)
		err := usecase.GenerateNewActivationToken(getContext(localizer, nil), id)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrAccountHasBeenActivated"), err)
	})

	t.Run("the token limit has been reached", func(t *testing.T) {
		id := users[1].ID
		tokens := seed.Tokens()
		fetchedTokens := []*models.Token{}
		for i := 0; i < limitOfActivationTokens+1; i++ {
			token := tokens[0]
			fetchedTokens = append(fetchedTokens, &token)
		}
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[1], nil).
			Once()
		mockTokenRepo.
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("*pgfilter.Filter")).
			Return(fetchedTokens, nil).
			Once()

		usecase := NewAuthUsecase(cfg)
		err := usecase.GenerateNewActivationToken(getContext(localizer, nil), id)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrReachedLimitOfActivationTokens"), err)
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
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("*pgfilter.Filter")).
			Return([]*models.Token{}, nil)
		mockTokenRepo.
			On("Store",
				mock.Anything,
				mock.AnythingOfType("*models.Token")).
			Return(fmt.Errorf("Error")).
			Once()

		usecase := NewAuthUsecase(cfg)
		err := usecase.GenerateNewActivationToken(getContext(localizer, nil), id)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrActivationTokenCannotBeCreated"), err)
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
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("*pgfilter.Filter")).
			Return([]*models.Token{}, nil)
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

		usecase := NewAuthUsecase(cfg)
		err := usecase.GenerateNewActivationToken(getContext(localizer, nil), id)
		require.Equal(t, nil, err)
	})
}

func TestGenerateNewActivationTokenForCurrentUser(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockTokenRepo := new(_tokenMocks.Repository)
	mockEmail := new(_emailMock.Email)
	mockSessStore := new(_sessionsMocks.Store)
	users := seed.Users()
	cfg := &Config{
		SessStore:       mockSessStore,
		UserRepo:        mockUserRepo,
		TokenRepo:       mockTokenRepo,
		Email:           mockEmail,
		ApplicationName: "AppName",
		FrontendURL:     "http://localhost:3000",
	}

	t.Run("user cannot be logged out", func(t *testing.T) {
		usecase := NewAuthUsecase(cfg)
		err := usecase.GenerateNewActivationTokenForCurrentUser(getContext(localizer, nil))
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrNotLoggedIn"), err)
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
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("*pgfilter.Filter")).
			Return([]*models.Token{}, nil)
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

		usecase := NewAuthUsecase(cfg)
		err := usecase.GenerateNewActivationTokenForCurrentUser(getContext(localizer, &users[1]))
		require.Equal(t, nil, err)
	})
}

func TestActivate(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockTokenRepo := new(_tokenMocks.Repository)
	mockEmail := new(_emailMock.Email)
	mockSessStore := new(_sessionsMocks.Store)
	users := seed.Users()
	tokens := seed.Tokens()
	cfg := &Config{
		SessStore:       mockSessStore,
		UserRepo:        mockUserRepo,
		TokenRepo:       mockTokenRepo,
		Email:           mockEmail,
		ApplicationName: "AppName",
		FrontendURL:     "http://localhost:3000",
	}

	t.Run("user account is activated", func(t *testing.T) {
		id := users[0].ID
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[0], nil).
			Once()
		usecase := NewAuthUsecase(cfg)
		token := uuid.New().String()
		_, err := usecase.Activate(getContext(localizer, nil), id, token)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrAccountHasBeenActivated"), err)
	})

	t.Run("no token found", func(t *testing.T) {
		id := users[1].ID
		tok := tokens[0]
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&users[1], nil).
			Once()
		mockTokenRepo.
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("*pgfilter.Filter")).
			Return([]*models.Token{}, nil).
			Once()
		usecase := NewAuthUsecase(cfg)
		_, err := usecase.Activate(getContext(localizer, nil), id, tok.Value)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrInvalidActivationToken"), err)
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
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("*pgfilter.Filter")).
			Return([]*models.Token{&tok}, nil).
			Once()
		mockUserRepo.
			On("Update",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(fmt.Errorf("error")).
			Once()
		usecase := NewAuthUsecase(cfg)
		_, err := usecase.Activate(getContext(localizer, nil), id, tok.Value)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrAccountCannotBeActivated"), err)
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
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("*pgfilter.Filter")).
			Return([]*models.Token{&tok}, nil).
			Once()
		mockTokenRepo.
			On("Delete",
				mock.Anything,
				[]int{id}).
			Return([]*models.Token{&tok}, nil).
			Once()
		mockTokenRepo.
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("*pgfilter.Filter")).
			Return([]*models.Token{&tok}, nil).
			Once()
		mockUserRepo.
			On("Update",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(nil).
			Once()
		mockTokenRepo.On("Delete", mock.Anything, []int{tok.ID}).Return([]*models.Token{&tok}, nil).Once()
		usecase := NewAuthUsecase(cfg)
		user, err := usecase.Activate(getContext(localizer, nil), id, tok.Value)
		require.Equal(t, nil, err)
		require.Equal(t, true, user.Activated)
	})
}

func TestGenerateNewResetPasswordToken(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockTokenRepo := new(_tokenMocks.Repository)
	mockEmail := new(_emailMock.Email)
	mockSessStore := new(_sessionsMocks.Store)
	users := seed.Users()
	email := users[0].Email
	cfg := &Config{
		SessStore:       mockSessStore,
		UserRepo:        mockUserRepo,
		TokenRepo:       mockTokenRepo,
		Email:           mockEmail,
		ApplicationName: "AppName",
		FrontendURL:     "http://localhost:3000",
	}

	t.Run("the token limit has been reached", func(t *testing.T) {
		tokens := seed.Tokens()
		fetchedTokens := []*models.Token{}
		for i := 0; i < limitOfActivationTokens+1; i++ {
			token := tokens[0]
			fetchedTokens = append(fetchedTokens, &token)
		}
		mockUserRepo.
			On("GetByEmail",
				mock.Anything,
				email).
			Return(&users[1], nil).
			Once()
		mockTokenRepo.
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("*pgfilter.Filter")).
			Return(fetchedTokens, nil).
			Once()

		usecase := NewAuthUsecase(cfg)
		err := usecase.GenerateNewResetPasswordToken(getContext(localizer, nil), email)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrReachedLimitOfResetPasswordTokens"), err)
	})

	t.Run("cannot create token", func(t *testing.T) {
		mockUserRepo.
			On("GetByEmail",
				mock.Anything,
				email).
			Return(&users[1], nil).
			Once()
		mockTokenRepo.
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("*pgfilter.Filter")).
			Return([]*models.Token{}, nil)
		mockTokenRepo.
			On("Store",
				mock.Anything,
				mock.AnythingOfType("*models.Token")).
			Return(fmt.Errorf("Error")).
			Once()

		usecase := NewAuthUsecase(cfg)
		err := usecase.GenerateNewResetPasswordToken(getContext(localizer, nil), email)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrResetPasswordTokenCannotBeCreated"), err)
	})

	t.Run("success", func(t *testing.T) {
		mockUserRepo.
			On("GetByEmail",
				mock.Anything,
				email).
			Return(&users[1], nil).
			Once()
		mockTokenRepo.
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("*pgfilter.Filter")).
			Return([]*models.Token{}, nil)
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

		usecase := NewAuthUsecase(cfg)
		err := usecase.GenerateNewResetPasswordToken(getContext(localizer, nil), email)
		require.Equal(t, nil, err)
	})
}

func TestResetPassword(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockTokenRepo := new(_tokenMocks.Repository)
	mockEmail := new(_emailMock.Email)
	mockSessStore := new(_sessionsMocks.Store)
	users := seed.Users()
	tokens := seed.Tokens()
	id := users[0].ID
	cfg := &Config{
		SessStore:       mockSessStore,
		UserRepo:        mockUserRepo,
		TokenRepo:       mockTokenRepo,
		Email:           mockEmail,
		ApplicationName: "AppName",
		FrontendURL:     "http://localhost:3000",
	}

	t.Run("no token found", func(t *testing.T) {
		tok := tokens[1]
		user := users[0]
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&user, nil).
			Once()
		mockTokenRepo.
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("*pgfilter.Filter")).
			Return([]*models.Token{}, nil).
			Once()
		usecase := NewAuthUsecase(cfg)
		err := usecase.ResetPassword(getContext(localizer, nil), id, tok.Value)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrInvalidResetPasswordToken"), err)
	})

	t.Run("cannot update user", func(t *testing.T) {
		tok := tokens[1]
		user := users[0]
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&user, nil).
			Once()
		mockTokenRepo.
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("*pgfilter.Filter")).
			Return([]*models.Token{&tok}, nil).
			Once()
		mockUserRepo.
			On("Update",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(fmt.Errorf("error")).
			Once()
		usecase := NewAuthUsecase(cfg)
		err := usecase.ResetPassword(getContext(localizer, nil), id, tok.Value)
		require.Equal(t, utils.GetErrorMsg(localizer, "ErrUserCannotBeUpdated"), err)
	})

	t.Run("success", func(t *testing.T) {
		tok := tokens[1]
		user := users[0]
		mockUserRepo.
			On("GetByID",
				mock.Anything,
				id).
			Return(&user, nil).
			Once()
		mockTokenRepo.
			On("Fetch",
				mock.Anything,
				mock.AnythingOfType("*pgfilter.Filter")).
			Return([]*models.Token{&tok}, nil).
			Once()
		mockUserRepo.
			On("Update",
				mock.Anything,
				mock.AnythingOfType("*models.User")).
			Return(nil).
			Once()
		mockTokenRepo.
			On("Delete", mock.Anything, []int{tok.ID}).
			Return([]*models.Token{&tok}, nil).
			Once()
		mockEmail.
			On("Send", mock.Anything, mock.Anything).
			Return(nil).
			Once()
		mockSessStore.On("GetAll").Return([]*sessions.Session{}, nil).Once()
		mockSessStore.On("DeleteByID", mock.AnythingOfType("[]string")).Return(nil).Once()

		usecase := NewAuthUsecase(cfg)
		err := usecase.ResetPassword(getContext(localizer, nil), id, tok.Value)
		require.Equal(t, nil, err)
	})
}

func getContext(localizer *i18n.Localizer, user *models.User) context.Context {
	return middleware.StoreLocalizerInContext(middleware.StoreUserInContext(context.Background(), user), localizer)
}
