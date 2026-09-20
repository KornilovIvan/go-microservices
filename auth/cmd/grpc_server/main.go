package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userAPI "github.com/ivankornilov/auth/internal/api/user"
	"github.com/ivankornilov/auth/internal/config"
	userRepository "github.com/ivankornilov/auth/internal/repository/user"
	userService "github.com/ivankornilov/auth/internal/service/user"
	desc "github.com/ivankornilov/auth/pkg/auth_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "local.env", "path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	cfg, err := config.Parse(configPath)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, cfg.PG.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", cfg.GRPC.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(
		s,
		userAPI.NewImplementation(userService.NewService(userRepository.NewRepository(pool))),
	)

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
