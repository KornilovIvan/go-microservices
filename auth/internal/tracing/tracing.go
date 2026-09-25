package tracing

import (
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"

	"github.com/ivankornilov/auth/internal/logger"
)

const serviceName = "auth"

func Init() {
	cfg, err := config.FromEnv()
	if err != nil {
		logger.Fatal("failed to init tracing", zap.Error(err))
	}

	cfg.Sampler = &config.SamplerConfig{
		Type:  "const",
		Param: 1,
	}

	_, err = cfg.InitGlobalTracer(serviceName)
	if err != nil {
		logger.Fatal("failed to init tracing", zap.Error(err))
	}
}
