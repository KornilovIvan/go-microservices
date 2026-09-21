package chat

import (
	"github.com/KornilovIvan/platform_common/pkg/db"
	"github.com/ivankornilov/chat-server/internal/repository"
	"github.com/ivankornilov/chat-server/internal/service"
)

type serv struct {
	chatRepository repository.ChatRepository
	txManager      db.TxManager
}

func NewService(chatRepository repository.ChatRepository, txManager db.TxManager) service.ChatService {
	return &serv{
		chatRepository: chatRepository,
		txManager:      txManager,
	}
}
