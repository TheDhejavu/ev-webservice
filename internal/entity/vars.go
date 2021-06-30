package entity

import "errors"

var (
	ErrNotFound        = errors.New("Entity does not  exist")
	ErrAlreadyExist    = errors.New("Entity alrrady  exist")
	ErrInvalidId       = errors.New("Invalid ID")
	ErrInvalidIdentity = errors.New("digits or password is wrong")
	ErrInvalidUser     = errors.New("username or password is wrong")
)
