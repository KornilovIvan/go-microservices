package config

import "github.com/joho/godotenv"

type Config struct {
	GRPC GRPCConfig
	PG   PGConfig
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

	pgCfg, err := NewPGConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		GRPC: grpcCfg,
		PG:   pgCfg,
	}, nil
}
