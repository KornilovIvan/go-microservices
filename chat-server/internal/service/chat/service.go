package chat

import (
	"github.com/ivankornilov/chat-server/internal/dbtx"
	"github.com/ivankornilov/chat-server/internal/repository"
	"github.com/ivankornilov/chat-server/internal/service"
)

type serv struct {
	chatRepository repository.ChatRepository
	txManager      dbtx.Manager
}

func NewService(chatRepository repository.ChatRepository, txManager dbtx.Manager) service.ChatService {
	return &serv{
		chatRepository: chatRepository,
		txManager:      txManager,
	}
}
