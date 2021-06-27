package entity

import "errors"

var (
	ErrNotFound     = errors.New("Entity does not  exist")
	ErrAlreadyExist = errors.New("Entity alrrady  exist")
	ErrInvalidId    = errors.New("Invalid ID")
)
