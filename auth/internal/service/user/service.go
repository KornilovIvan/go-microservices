package user

import (
	"github.com/ivankornilov/auth/internal/client/db"
	"github.com/ivankornilov/auth/internal/repository"
	"github.com/ivankornilov/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	logRepository  repository.LogRepository
	txManager      db.TxManager
}

func NewService(
	userRepository repository.UserRepository,
	logRepository repository.LogRepository,
	txManager db.TxManager,
) service.UserService {
	return &serv{
		userRepository: userRepository,
		logRepository:  logRepository,
		txManager:      txManager,
	}
}
