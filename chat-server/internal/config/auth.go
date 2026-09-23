package config

import (
	"errors"
	"net"
	"os"
)

const (
	authGRPCHostEnvName = "AUTH_GRPC_HOST"
	authGRPCPortEnvName = "AUTH_GRPC_PORT"
)

type AuthConfig interface {
	Address() string
}

type authConfig struct {
	host string
	port string
}

func NewAuthConfig() (AuthConfig, error) {
	host := os.Getenv(authGRPCHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("auth grpc host not found")
	}

	port := os.Getenv(authGRPCPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("auth grpc port not found")
	}

	return &authConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *authConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
