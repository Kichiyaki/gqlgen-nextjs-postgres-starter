package errors

import (
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func Wrap(msg string, errors ...error) error {
	return &gqlerror.Error{
		Message: msg,
		Extensions: map[string]interface{}{
			"detailed": errors,
		},
	}

}

func ToGqlError(err error) *gqlerror.Error {
	if gqlError, ok := err.(*gqlerror.Error); ok {
		return gqlError
	}
	return &gqlerror.Error{
		Message: err.Error(),
	}
}
