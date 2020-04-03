package seed

import (
	"fmt"
	"time"

	"backend/models"

	"github.com/google/uuid"
)

func Users(limit int) []models.User {
	users := []models.User{}
	for i := 1; i <= limit; i++ {
		login := fmt.Sprintf("%dLogin%d", i, i)
		role := models.UserAdminRole
		if i%3 == 0 {
			role = models.UserAdminRole
		}
		users = append(users, models.User{
			ID:                 i + 150,
			Slug:               fmt.Sprintf("%d-%s", i, login),
			Login:              login,
			Email:              fmt.Sprintf("%s@gmail.com", login),
			Password:           fmt.Sprintf("%dpasswordinoElorino", i),
			Role:               role,
			ActivationToken:    uuid.New().String(),
			ResetPasswordToken: uuid.New().String(),
			Activated:          role == models.UserAdminRole || i%2 == 0,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		})
	}
	return users
}
