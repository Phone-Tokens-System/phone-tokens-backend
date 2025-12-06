package tokens

import "errors"

var (
	ErrNotFound  = errors.New("token not found")
	ErrForbidden = errors.New("token does not belong to the user")
)
