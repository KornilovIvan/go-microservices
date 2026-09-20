package model

import "errors"

var (
	ErrNameEmailPasswordRequired = errors.New("name, email and password are required")
	ErrPasswordMismatch          = errors.New("password and password_confirm do not match")
	ErrNameOrEmailRequired       = errors.New("name or email is required")
	ErrUserNotFound              = errors.New("user not found")
	ErrUserAlreadyExists         = errors.New("user with this email already exists")
)
