package model

import "errors"

var (
	ErrChatIDRequired      = errors.New("chat_id is required")
	ErrFromAndTextRequired = errors.New("from and text are required")
	ErrChatNotFound        = errors.New("chat not found")
)
