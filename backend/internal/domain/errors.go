package domain

import "errors"

var (
	ErrInvalidInput       = errors.New("invalid_input")
	ErrInvalidCredentials = errors.New("invalid_credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrUserExists         = errors.New("user_exists")
	ErrUserNotFound       = errors.New("user_not_found")
)
