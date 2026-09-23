package auth

import (
	"context"
	"errors"

	"github.com/ivankornilov/auth/internal/model"
	"github.com/ivankornilov/auth/internal/utils"
)

func (s *serv) Login(ctx context.Context, username string, password string) (string, error) {
	user, err := s.userRepository.GetByName(ctx, username)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return "", model.ErrInvalidCredentials
		}
		return "", err
	}

	if user.Info.Password != password {
		return "", model.ErrInvalidCredentials
	}

	refreshToken, err := utils.GenerateToken(user.Info.Name, user.Info.Role, s.refreshSecretKey, s.refreshTokenExpiration)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
