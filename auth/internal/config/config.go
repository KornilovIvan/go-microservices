package config

import "github.com/joho/godotenv"

type Config struct {
	GRPC    GRPCConfig
	HTTP    HTTPConfig
	PG      PGConfig
	Swagger SwaggerConfig
	Token   TokenConfig
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

func Parse(path string) (*Config, error) {
	err := Load(path)
	if err != nil {
		return nil, err
	}

	grpcCfg, err := NewGRPCConfig()
	if err != nil {
		return nil, err
	}

	httpCfg, err := NewHTTPConfig()
	if err != nil {
		return nil, err
	}

	pgCfg, err := NewPGConfig()
	if err != nil {
		return nil, err
	}

	swaggerCfg, err := NewSwaggerConfig()
	if err != nil {
		return nil, err
	}

	tokenCfg, err := NewTokenConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		GRPC:    grpcCfg,
		HTTP:    httpCfg,
		PG:      pgCfg,
		Swagger: swaggerCfg,
		Token:   tokenCfg,
	}, nil
}
