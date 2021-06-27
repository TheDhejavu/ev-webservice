package errors

import "errors"

var (
	ErrEntityDoesNotExist = errors.New("Entity does not  exist")
	ErrEntityAlreadyExist = errors.New("Entity alrrady  exist")
)
