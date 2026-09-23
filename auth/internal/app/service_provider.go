package app

import (
	"context"
	"log"

	"github.com/KornilovIvan/platform_common/pkg/db"
	"github.com/KornilovIvan/platform_common/pkg/db/pg"
	"github.com/KornilovIvan/platform_common/pkg/db/transaction"
	accessAPI "github.com/ivankornilov/auth/internal/api/access"
	userAPI "github.com/ivankornilov/auth/internal/api/user"
	"github.com/ivankornilov/auth/internal/config"
	"github.com/ivankornilov/auth/internal/repository"
	logRepository "github.com/ivankornilov/auth/internal/repository/log"
	userRepository "github.com/ivankornilov/auth/internal/repository/user"
	"github.com/ivankornilov/auth/internal/service"
	accessService "github.com/ivankornilov/auth/internal/service/access"
	authService "github.com/ivankornilov/auth/internal/service/auth"
	userService "github.com/ivankornilov/auth/internal/service/user"
)

type serviceProvider struct {
	cfg *config.Config

	dbClient       db.Client
	txManager      db.TxManager
	userRepository repository.UserRepository
	logRepository  repository.LogRepository
	userService    service.UserService
	authService    service.AuthService
	accessService  service.AccessService
	userImpl       *userAPI.Implementation
	accessImpl     *accessAPI.Implementation
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

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepository(s.DBClient(ctx))
	}

	return s.logRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.LogRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.UserRepository(ctx), s.Config().Token)
	}

	return s.authService
}

func (s *serviceProvider) AccessService(_ context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessService.NewService(s.Config().Token)
	}

	return s.accessService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *userAPI.Implementation {
	if s.userImpl == nil {
		s.userImpl = userAPI.NewImplementation(s.UserService(ctx), s.AuthService(ctx))
	}

	return s.userImpl
}

func (s *serviceProvider) AccessImpl(ctx context.Context) *accessAPI.Implementation {
	if s.accessImpl == nil {
		s.accessImpl = accessAPI.NewImplementation(s.AccessService(ctx))
	}

	return s.accessImpl
}

func (s *serviceProvider) Close() {
	if s.dbClient != nil {
		_ = s.dbClient.Close()
	}
}
