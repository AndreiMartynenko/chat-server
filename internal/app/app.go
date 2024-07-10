package app

import (
	"context"
	"flag"
	"log"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/AndreiMartynenko/chat-server/internal/config"
	desc "github.com/AndreiMartynenko/chat-server/pkg/chat_v1"
	"github.com/AndreiMartynenko/common/pkg/closer"
)

// App structure contains main application structures.
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

var configPath string

func init() {
	flag.StringVar(&configPath, "config", ".env", "Path to config file")
}

// NewApp creates new App object.
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Run executes the application.
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGrpcServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGrpcServer,
		a.initTracing,
		a.initChatChannels,
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
	// Parse the command-line flags from os.Args[1:].
	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGrpcServer(ctx context.Context) error {
	cfgGrpc := a.serviceProvider.GrpcConfig()
	creds, err := credentials.NewServerTLSFromFile(cfgGrpc.CertPath(), cfgGrpc.KeyPath())
	if err != nil {
		return err
	}

	c := a.serviceProvider.InterceptorClient()
	a.grpcServer = grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(c.PolicyInterceptor),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	// Upon the client's request, the server will automatically provide information on the supported methods.
	reflection.Register(a.grpcServer)

	// Register service with corresponded interface.
	desc.RegisterChatV1Server(a.grpcServer, a.serviceProvider.ChatImpl(ctx))

	return nil
}

func (a *App) initTracing(ctx context.Context) error {
	cfg := a.serviceProvider.TracingConfig()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName(cfg.ServiceName()),
		),
	)
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(cfg.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return err
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	closer.Add(func() error {
		return tracerProvider.Shutdown(ctx)
	})

	return nil
}

func (a *App) initChatChannels(ctx context.Context) error {
	return a.serviceProvider.ChatService(ctx).InitChannels(ctx)
}

func (a *App) runGrpcServer() error {
	// Open IP and port for server.
	lis, err := net.Listen(a.serviceProvider.GrpcConfig().Transport(), a.serviceProvider.GrpcConfig().Address())
	if err != nil {
		return err
	}

	log.Printf("gRPC server running on %v", a.serviceProvider.GrpcConfig().Address())

	err = a.grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}
