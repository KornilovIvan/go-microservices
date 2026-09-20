package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	chatAPI "github.com/ivankornilov/chat-server/internal/api/chat"
	"github.com/ivankornilov/chat-server/internal/config"
	"github.com/ivankornilov/chat-server/internal/dbtx"
	chatRepository "github.com/ivankornilov/chat-server/internal/repository/chat"
	chatService "github.com/ivankornilov/chat-server/internal/service/chat"
	desc "github.com/ivankornilov/chat-server/pkg/chat_v1"
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
	desc.RegisterChatV1Server(
		s,
		chatAPI.NewImplementation(
			chatService.NewService(
				chatRepository.NewRepository(pool),
				dbtx.NewManager(pool),
			),
		),
	)

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
