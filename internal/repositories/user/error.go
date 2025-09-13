package user

import "errors"

var (
	ErrNullUser       = errors.New("user is null")
	ErrNoRowsAffected = errors.New("no rows affected")
)
