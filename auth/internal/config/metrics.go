package config

import (
	"errors"
	"net"
	"os"
)

const (
	metricsHostEnvName = "METRICS_HOST"
	metricsPortEnvName = "METRICS_PORT"
)

type MetricsConfig interface {
	Address() string
}

type metricsConfig struct {
	host string
	port string
}

func NewMetricsConfig() (MetricsConfig, error) {
	host := os.Getenv(metricsHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("metrics host not found")
	}

	port := os.Getenv(metricsPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("metrics port not found")
	}

	return &metricsConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *metricsConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
