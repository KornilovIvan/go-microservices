package auth

import (
	"context"

	"github.com/ivankornilov/auth/internal/model"
	"github.com/ivankornilov/auth/internal/utils"
)

func (s *serv) GetAccessToken(_ context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, s.refreshSecretKey)
	if err != nil {
		return "", model.ErrInvalidRefreshToken
	}

	token, err := utils.GenerateToken(claims.Username, claims.Role, s.accessSecretKey, s.accessTokenExpiration)
	if err != nil {
		return "", err
	}

	return token, nil
}
