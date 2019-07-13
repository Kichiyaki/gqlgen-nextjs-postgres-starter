package models

import (
	"fmt"

	"github.com/kichiyaki/structs"
)

type TokenFilter struct {
	Type      string `structs:"type"`
	Value     string `structs:"value"`
	CreatedAt string `structs:"created_at"`
	UserID    string `structs:"user_id"`
}

func (f *TokenFilter) ToMap() map[string]string {
	m := make(map[string]string)
	if f != nil {
		fi := structs.Map(f)
		for key, value := range fi {
			strKey := fmt.Sprintf("%v", key)
			strValue := fmt.Sprintf("%v", value)

			m[strKey] = strValue
		}
	}
	return m
}
