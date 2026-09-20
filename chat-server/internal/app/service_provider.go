package app

import (
	"context"
	"log"

	chatAPI "github.com/ivankornilov/chat-server/internal/api/chat"
	"github.com/ivankornilov/chat-server/internal/client/db"
	"github.com/ivankornilov/chat-server/internal/client/db/pg"
	"github.com/ivankornilov/chat-server/internal/client/db/transaction"
	"github.com/ivankornilov/chat-server/internal/config"
	"github.com/ivankornilov/chat-server/internal/repository"
	chatRepository "github.com/ivankornilov/chat-server/internal/repository/chat"
	"github.com/ivankornilov/chat-server/internal/service"
	chatService "github.com/ivankornilov/chat-server/internal/service/chat"
)

type serviceProvider struct {
	cfg *config.Config

	dbClient       db.Client
	txManager      db.TxManager
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

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.Config().PG.DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.DBClient(ctx))
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
	if s.dbClient != nil {
		_ = s.dbClient.Close()
	}
}
