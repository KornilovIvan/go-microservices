package user

import (
	"github.com/ivankornilov/auth/internal/repository"
	"github.com/ivankornilov/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) service.UserService {
	return &serv{
		userRepository: userRepository,
	}
}
