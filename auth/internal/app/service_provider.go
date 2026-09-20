package app

import (
	"context"
	"log"

	userAPI "github.com/ivankornilov/auth/internal/api/user"
	"github.com/ivankornilov/auth/internal/client/db"
	"github.com/ivankornilov/auth/internal/client/db/pg"
	"github.com/ivankornilov/auth/internal/config"
	"github.com/ivankornilov/auth/internal/repository"
	userRepository "github.com/ivankornilov/auth/internal/repository/user"
	"github.com/ivankornilov/auth/internal/service"
	userService "github.com/ivankornilov/auth/internal/service/user"
)

type serviceProvider struct {
	cfg *config.Config

	dbClient       db.Client
	userRepository repository.UserRepository
	userService    service.UserService
	userImpl       *userAPI.Implementation
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

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx))
	}

	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *userAPI.Implementation {
	if s.userImpl == nil {
		s.userImpl = userAPI.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}

func (s *serviceProvider) Close() {
	if s.dbClient != nil {
		_ = s.dbClient.Close()
	}
}
