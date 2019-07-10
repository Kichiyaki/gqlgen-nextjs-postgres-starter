package seed

import (
	"time"

	"github.com/kichiyaki/graphql-starter/backend/models"
)

func Users() []models.User {
	return []models.User{
		models.User{
			ID:        340,
			Login:     "Logineszko",
			Password:  "test123T",
			Email:     "test@test.com",
			Role:      models.AdministrativeRole,
			Activated: true,
			Slug:      "1-logineszko",
			CreatedAt: time.Now(),
		},
		models.User{
			ID:        341,
			Login:     "Logineszko2",
			Password:  "test123T",
			Email:     "tesasdt@test.com",
			Role:      models.DefaultRole,
			Activated: false,
			Slug:      "2-logineszko2",
			CreatedAt: time.Now(),
		},
	}
}
