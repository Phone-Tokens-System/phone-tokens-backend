package tokens

import "errors"

var (
	ErrNotFound          = errors.New("token not found")
	ErrForbidden         = errors.New("token does not belong to the user")
	ErrInvalidPermission = errors.New("invalid permission")
	ErrInvalidStatus     = errors.New("invalid status")
)
