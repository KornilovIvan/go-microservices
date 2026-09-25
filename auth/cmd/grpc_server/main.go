package main

import (
	"context"
	"flag"
	"log"

	"github.com/ivankornilov/auth/internal/app"
	"github.com/ivankornilov/auth/internal/logger"
	"github.com/ivankornilov/auth/internal/metric"
)

var (
	configPath string
	logLevel   string
)

func init() {
	flag.StringVar(&configPath, "config-path", "local.env", "path to config file")
	flag.StringVar(&logLevel, "l", "info", "log level")
}

func main() {
	flag.Parse()

	core, err := logger.Core(logLevel)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	logger.Init(core)

	ctx := context.Background()

	err = metric.Init(ctx)
	if err != nil {
		log.Fatalf("failed to init metrics: %v", err)
	}

	a, err := app.NewApp(ctx, configPath)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
