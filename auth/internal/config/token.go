package config

import (
	"errors"
	"os"
	"time"
)

const (
	refreshTokenSecretEnvName = "REFRESH_TOKEN_SECRET_KEY"
	accessTokenSecretEnvName  = "ACCESS_TOKEN_SECRET_KEY"

	refreshTokenExpiration = 60 * time.Minute
	accessTokenExpiration  = 5 * time.Minute
)

type TokenConfig interface {
	RefreshSecretKey() []byte
	AccessSecretKey() []byte
	RefreshTokenExpiration() time.Duration
	AccessTokenExpiration() time.Duration
}

type tokenConfig struct {
	refreshSecretKey []byte
	accessSecretKey  []byte
}

func NewTokenConfig() (TokenConfig, error) {
	refreshSecretKey := os.Getenv(refreshTokenSecretEnvName)
	if len(refreshSecretKey) == 0 {
		return nil, errors.New("refresh token secret key not found")
	}

	accessSecretKey := os.Getenv(accessTokenSecretEnvName)
	if len(accessSecretKey) == 0 {
		return nil, errors.New("access token secret key not found")
	}

	return &tokenConfig{
		refreshSecretKey: []byte(refreshSecretKey),
		accessSecretKey:  []byte(accessSecretKey),
	}, nil
}

func (cfg *tokenConfig) RefreshSecretKey() []byte {
	return cfg.refreshSecretKey
}

func (cfg *tokenConfig) AccessSecretKey() []byte {
	return cfg.accessSecretKey
}

func (cfg *tokenConfig) RefreshTokenExpiration() time.Duration {
	return refreshTokenExpiration
}

func (cfg *tokenConfig) AccessTokenExpiration() time.Duration {
	return accessTokenExpiration
}
