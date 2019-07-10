package models

type List struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}
