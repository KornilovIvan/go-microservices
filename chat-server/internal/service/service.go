package service

//go:generate minimock -i github.com/ivankornilov/chat-server/internal/service.ChatService -o ./mocks -s _minimock.go

import (
	"context"

	"github.com/ivankornilov/chat-server/internal/model"
)

type ChatService interface {
	Create(ctx context.Context, usernames []string) (int64, error)
	Delete(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, message *model.Message) error
}
