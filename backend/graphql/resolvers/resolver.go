package resolvers

import (
	"github.com/kichiyaki/graphql-starter/backend/auth"
	"github.com/kichiyaki/graphql-starter/backend/graphql/generated"
	"github.com/kichiyaki/graphql-starter/backend/user"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	UserUcase user.Usecase
	AuthUcase auth.Usecase
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) UserList() generated.UserListResolver {
	return &userListResolver{r}
}

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }
