package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/ivankornilov/auth/internal/config"
	"github.com/ivankornilov/auth/internal/interceptor"
	desc "github.com/ivankornilov/auth/pkg/auth_v1"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
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

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to run GRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
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
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)
	reflection.Register(a.grpcServer)
	desc.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.UserImpl(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := desc.RegisterAuthV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.Config().GRPC.Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.Config().HTTP.Address(),
		Handler: corsMiddleware.Handler(mux),
	}

	return nil
}

func (a *App) runGRPCServer() error {
	address := a.serviceProvider.Config().GRPC.Address()
	log.Printf("gRPC server listening at %s", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	return a.grpcServer.Serve(lis)
}

func (a *App) runHTTPServer() error {
	address := a.serviceProvider.Config().HTTP.Address()
	log.Printf("HTTP server listening at %s", address)

	return a.httpServer.ListenAndServe()
}
