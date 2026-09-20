package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"

	chatAPI "github.com/ivankornilov/chat-server/internal/api/chat"
	"github.com/ivankornilov/chat-server/internal/config"
	"github.com/ivankornilov/chat-server/internal/dbtx"
	"github.com/ivankornilov/chat-server/internal/repository"
	chatRepository "github.com/ivankornilov/chat-server/internal/repository/chat"
	"github.com/ivankornilov/chat-server/internal/service"
	chatService "github.com/ivankornilov/chat-server/internal/service/chat"
)

type serviceProvider struct {
	cfg *config.Config

	pool           *pgxpool.Pool
	txManager      dbtx.Manager
	chatRepository repository.ChatRepository
	chatService    service.ChatService
	chatImpl       *chatAPI.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) Config() *config.Config {
	return s.cfg
}

func (s *serviceProvider) Pool(ctx context.Context) *pgxpool.Pool {
	if s.pool == nil {
		pool, err := pgxpool.Connect(ctx, s.Config().PG.DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		s.pool = pool
	}

	return s.pool
}

func (s *serviceProvider) TxManager(ctx context.Context) dbtx.Manager {
	if s.txManager == nil {
		s.txManager = dbtx.NewManager(s.Pool(ctx))
	}

	return s.txManager
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.Pool(ctx))
	}

	return s.chatRepository
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(s.ChatRepository(ctx), s.TxManager(ctx))
	}

	return s.chatService
}

func (s *serviceProvider) ChatImpl(ctx context.Context) *chatAPI.Implementation {
	if s.chatImpl == nil {
		s.chatImpl = chatAPI.NewImplementation(s.ChatService(ctx))
	}

	return s.chatImpl
}

func (s *serviceProvider) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}
