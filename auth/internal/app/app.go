package app

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/ivankornilov/auth/internal/config"
	"github.com/ivankornilov/auth/internal/interceptor"
	accessDesc "github.com/ivankornilov/auth/pkg/access_v1"
	desc "github.com/ivankornilov/auth/pkg/auth_v1"
	_ "github.com/ivankornilov/auth/statik"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
	metricsServer   *http.Server
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
	wg.Add(4)

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

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			log.Fatalf("failed to run Swagger server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runMetricsServer()
		if err != nil {
			log.Fatalf("failed to run metrics server: %v", err)
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
		a.initSwaggerServer,
		a.initMetricsServer,
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
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.LogInterceptor,
				interceptor.MetricsInterceptor,
				interceptor.ValidateInterceptor,
			),
		),
	)
	reflection.Register(a.grpcServer)
	desc.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.UserImpl(ctx))
	accessDesc.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessImpl(ctx))

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

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:    a.serviceProvider.Config().Swagger.Address(),
		Handler: mux,
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

func (a *App) runSwaggerServer() error {
	address := a.serviceProvider.Config().Swagger.Address()
	log.Printf("Swagger server listening at %s", address)

	return a.swaggerServer.ListenAndServe()
}

func (a *App) initMetricsServer(_ context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	a.metricsServer = &http.Server{
		Addr:    a.serviceProvider.Config().Metrics.Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) runMetricsServer() error {
	address := a.serviceProvider.Config().Metrics.Address()
	log.Printf("metrics server listening at %s", address)

	return a.metricsServer.ListenAndServe()
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
