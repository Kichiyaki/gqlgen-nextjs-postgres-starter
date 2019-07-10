package models

type UserInput struct {
	Login     *string
	Email     *string
	Password  *string
	Role      *string
	Activated *bool
}

func (input UserInput) ToUser() *User {
	user := &User{}
	if input.Login != nil {
		user.Login = *input.Login
	}
	if input.Password != nil {
		user.Password = *input.Password
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Role != nil {
		user.Role = *input.Role
	}
	if input.Activated != nil {
		user.Activated = *input.Activated
	}

	return user
}
