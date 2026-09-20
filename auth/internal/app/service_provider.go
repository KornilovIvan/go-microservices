package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"

	userAPI "github.com/ivankornilov/auth/internal/api/user"
	"github.com/ivankornilov/auth/internal/config"
	"github.com/ivankornilov/auth/internal/repository"
	userRepository "github.com/ivankornilov/auth/internal/repository/user"
	"github.com/ivankornilov/auth/internal/service"
	userService "github.com/ivankornilov/auth/internal/service/user"
)

type serviceProvider struct {
	cfg *config.Config

	pool           *pgxpool.Pool
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

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.Pool(ctx))
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
	if s.pool != nil {
		s.pool.Close()
	}
}
