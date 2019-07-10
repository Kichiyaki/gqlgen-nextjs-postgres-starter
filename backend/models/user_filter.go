package models

type UserFilter struct {
	ID            string `json:"id" structs:"id"`
	Login         string `json:"login" structs:"login"`
	Role          string `json:"role" structs:"role"`
	Email         string `json:"email" structs:"email"`
	OnlyActivated bool   `json:"onlyActivated" structs:"-"`
	Sort          string `json:"sort" structs:"sort"`
	SortBy        string `json:"sortBy" structs:"sortBy"`
}
