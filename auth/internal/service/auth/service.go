package auth

import (
	"time"

	"github.com/ivankornilov/auth/internal/config"
	"github.com/ivankornilov/auth/internal/repository"
	"github.com/ivankornilov/auth/internal/service"
)

type serv struct {
	userRepository         repository.UserRepository
	refreshSecretKey       []byte
	accessSecretKey        []byte
	refreshTokenExpiration time.Duration
	accessTokenExpiration  time.Duration
}

func NewService(userRepository repository.UserRepository, tokenConfig config.TokenConfig) service.AuthService {
	return &serv{
		userRepository:         userRepository,
		refreshSecretKey:       tokenConfig.RefreshSecretKey(),
		accessSecretKey:        tokenConfig.AccessSecretKey(),
		refreshTokenExpiration: tokenConfig.RefreshTokenExpiration(),
		accessTokenExpiration:  tokenConfig.AccessTokenExpiration(),
	}
}
