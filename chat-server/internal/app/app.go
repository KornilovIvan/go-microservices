package app

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ivankornilov/chat-server/internal/config"
	desc "github.com/ivankornilov/chat-server/pkg/chat_v1"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	configPath      string
}

func NewApp(ctx context.Context, configPath string) (*App, error) {
	a := &App{configPath: configPath}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer a.serviceProvider.Close()

	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	cfg, err := config.Parse(a.configPath)
	if err != nil {
		return err
	}

	if a.serviceProvider == nil {
		a.serviceProvider = newServiceProvider()
	}
	a.serviceProvider.cfg = cfg

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	if a.serviceProvider == nil {
		a.serviceProvider = newServiceProvider()
	}

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(a.serviceProvider.AccessInterceptor(ctx).Unary),
	)
	reflection.Register(a.grpcServer)
	desc.RegisterChatV1Server(a.grpcServer, a.serviceProvider.ChatImpl(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	address := a.serviceProvider.Config().GRPC.Address()
	log.Printf("server listening at %s", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	return a.grpcServer.Serve(lis)
}
