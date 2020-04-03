package models

type UserList struct {
	Total int     `json:"total"`
	Items []*User `json:"items"`
}
