package common

import "errors"

var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
