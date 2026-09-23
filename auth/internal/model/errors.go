package model

import "errors"

var (
	ErrNameEmailPasswordRequired = errors.New("name, email and password are required")
	ErrPasswordMismatch          = errors.New("password and password_confirm do not match")
	ErrNameOrEmailRequired       = errors.New("name or email is required")
	ErrUserNotFound              = errors.New("user not found")
	ErrUserAlreadyExists         = errors.New("user with this email already exists")
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrInvalidRefreshToken       = errors.New("invalid refresh token")
	ErrAccessTokenInvalid        = errors.New("access token is invalid")
	ErrAuthHeaderMissing         = errors.New("authorization header is not provided")
	ErrAuthHeaderInvalid         = errors.New("invalid authorization header format")
	ErrAccessDenied              = errors.New("access denied")
)
