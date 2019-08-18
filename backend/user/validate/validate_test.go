package validate

import (
	"fmt"
	"testing"

	"github.com/kichiyaki/graphql-starter/backend/seed"
	"github.com/kichiyaki/graphql-starter/backend/utils"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/kichiyaki/graphql-starter/backend/models"
	"github.com/stretchr/testify/require"
)

var localizer *i18n.Localizer

func init() {
	localizer = utils.GetLocalizer(language.Polish, "../../i18n/locales/active.pl.json")
}

func TestValidate(t *testing.T) {
	cfg := NewConfig(localizer)
	cfg.Email = true
	cfg.Login = true
	cfg.Password = true
	cfg.Role = true

	fmt.Println(utils.GetErrorMsgWithDataAndPluralCount(localizer, "ErrMinimalLengthOfLogin", map[string]interface{}{
		"Characters": minimalLengthOfLogin,
	}, minimalLengthOfLogin).Error())

	t.Run("Login", func(t *testing.T) {
		t.Run("must provide login", func(t *testing.T) {
			err := cfg.Validate(models.User{})
			require.Equal(t, utils.GetErrorMsg(localizer, "ErrMustProvideLogin"), err)
		})

		t.Run(fmt.Sprintf("login must be at least %d characters", minimalLengthOfLogin), func(t *testing.T) {
			err := cfg.Validate(models.User{Login: "a"})
			require.Equal(t,
				utils.GetErrorMsgWithDataAndPluralCount(localizer, "ErrMinimalLengthOfLogin", map[string]interface{}{
					"Characters": minimalLengthOfLogin,
				}, minimalLengthOfLogin),
				err)
		})

		t.Run(fmt.Sprintf("login must be shorter than %d characters", maximalLengthOfLogin), func(t *testing.T) {
			l := "a"
			for i := 1; i <= maximalLengthOfLogin+5; i++ {
				l += "b"
			}
			err := cfg.Validate(models.User{Login: l})
			require.Equal(t,
				utils.GetErrorMsgWithDataAndPluralCount(localizer, "ErrMaximalLengthOfLogin", map[string]interface{}{
					"Characters": maximalLengthOfLogin,
				}, maximalLengthOfLogin),
				err)
		})
	})

	t.Run("Password", func(t *testing.T) {
		t.Run("must provide password", func(t *testing.T) {
			err := cfg.Validate(models.User{Login: "test12"})
			require.Equal(t, utils.GetErrorMsg(localizer, "ErrMustProvidePassword"), err)
		})

		t.Run(fmt.Sprintf("password must be at least %d characters", minimalLengthOfPassword), func(t *testing.T) {
			err := cfg.Validate(models.User{Login: "abcd", Password: "12"})
			require.Equal(t,
				utils.GetErrorMsgWithDataAndPluralCount(localizer, "ErrMinimalLengthOfPassword", map[string]interface{}{
					"Characters": minimalLengthOfPassword,
				}, minimalLengthOfPassword),
				err)
		})

		t.Run(fmt.Sprintf("password must be shorter than %d characters", maximalLengthOfPassword), func(t *testing.T) {
			l := "a"
			for i := 1; i <= maximalLengthOfPassword+5; i++ {
				l += "b"
			}
			err := cfg.Validate(models.User{Login: "asds", Password: l})
			require.Equal(t,
				utils.GetErrorMsgWithDataAndPluralCount(localizer, "ErrMaximalLengthOfPassword", map[string]interface{}{
					"Characters": maximalLengthOfPassword,
				}, maximalLengthOfPassword),
				err)
		})

		t.Run("password must contain at least one uppercase letter", func(t *testing.T) {
			err := cfg.Validate(models.User{Login: "asds", Password: "asdasdadasda"})
			require.Equal(t, utils.GetErrorMsg(localizer, "ErrPasswordMustContainAtLeastOneUppercase"), err)
		})

		t.Run("password must contain at least one lowercase letter", func(t *testing.T) {
			err := cfg.Validate(models.User{Login: "asds", Password: "ASDASDADSADA"})
			require.Equal(t, utils.GetErrorMsg(localizer, "ErrPasswordMustContainAtLeastOneLowercase"), err)
		})

		t.Run("password must contain at least one digit", func(t *testing.T) {
			err := cfg.Validate(models.User{Login: "asds", Password: "ASDASDADSADAasd"})
			require.Equal(t, utils.GetErrorMsg(localizer, "ErrPasswordMustContainAtLeastOneDigit"), err)
		})
	})

	t.Run("Email", func(t *testing.T) {
		t.Run("must provide email address", func(t *testing.T) {
			u := seed.Users()[0]
			u.Email = ""
			err := cfg.Validate(u)
			require.Equal(t, utils.GetErrorMsg(localizer, "ErrMustProvideEmailAddress"), err)
		})

		t.Run("invalid email address", func(t *testing.T) {
			u := seed.Users()[0]
			u.Email = "tesasd2sd..as"
			err := cfg.Validate(u)
			require.Equal(t,
				utils.GetErrorMsgWithData(localizer, "ErrInvalidEmailAddress", map[string]interface{}{
					"Email": u.Email,
				}),
				err)
		})
	})

	t.Run("Role", func(t *testing.T) {
		t.Run("invalid role", func(t *testing.T) {
			u := seed.Users()[0]
			u.Role = "eloszka"
			err := cfg.Validate(u)
			require.Equal(t,
				utils.GetErrorMsgWithData(localizer, "ErrInvalidRole", map[string]interface{}{
					"Role": u.Role,
				}),
				err)
		})
	})

	t.Run("success", func(t *testing.T) {
		u := seed.Users()[0]
		err := cfg.Validate(u)
		require.Equal(t, nil, err)
	})
}
