package repository

import (
	"context"

	"github.com/ivankornilov/chat-server/internal/model"
)

type ChatRepository interface {
	Create(ctx context.Context) (int64, error)
	AddUser(ctx context.Context, chatID int64, username string) error
	Delete(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, message *model.Message) error
}
