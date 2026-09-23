package access

import (
	"github.com/ivankornilov/auth/internal/config"
	"github.com/ivankornilov/auth/internal/model"
	"github.com/ivankornilov/auth/internal/service"
)

type serv struct {
	accessSecretKey []byte
	accessibleRoles map[string]string
}

func NewService(tokenConfig config.TokenConfig) service.AccessService {
	return &serv{
		accessSecretKey: tokenConfig.AccessSecretKey(),
		accessibleRoles: map[string]string{
			model.ExamplePath: "ADMIN",
		},
	}
}
