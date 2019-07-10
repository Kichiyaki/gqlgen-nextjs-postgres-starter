package seed

import (
	"github.com/kichiyaki/graphql-starter/backend/models"
)

func Tokens() []models.Token {
	users := Users()
	return []models.Token{
		models.Token{
			ID:     340,
			Value:  "token1",
			UserID: users[1].ID,
			Type:   models.AccountActivationTokenType,
		},
	}
}
