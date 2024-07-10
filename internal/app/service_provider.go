package app

import (
	"context"
	"log"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	descAccess "github.com/AndreiMartynenko/auth/pkg/access_v1"
	"github.com/AndreiMartynenko/chat-server/internal/api/chat"
	rpc "github.com/AndreiMartynenko/chat-server/internal/client/rpc"
	rpcAuth "github.com/AndreiMartynenko/chat-server/internal/client/rpc/auth"
	"github.com/AndreiMartynenko/chat-server/internal/config"
	"github.com/AndreiMartynenko/chat-server/internal/config/env"
	"github.com/AndreiMartynenko/chat-server/internal/interceptor"
	"github.com/AndreiMartynenko/chat-server/internal/repository"
	chatRepository "github.com/AndreiMartynenko/chat-server/internal/repository/chat"
	logRepository "github.com/AndreiMartynenko/chat-server/internal/repository/log"
	messagesRepository "github.com/AndreiMartynenko/chat-server/internal/repository/messages"
	"github.com/AndreiMartynenko/chat-server/internal/service"
	chatService "github.com/AndreiMartynenko/chat-server/internal/service/chat"
	"github.com/AndreiMartynenko/common/pkg/closer"
	"github.com/AndreiMartynenko/common/pkg/db"
	"github.com/AndreiMartynenko/common/pkg/db/pg"
	"github.com/AndreiMartynenko/common/pkg/db/transaction"
)

type serviceProvider struct {
	pgConfig      config.PgConfig
	grpcConfig    config.GrpcConfig
	authConfig    config.AuthConfig
	tracingConfig config.TracingConfig

	authClient        rpc.AuthClient
	dbClient          db.Client
	txManager         db.TxManager
	interceptorClient *interceptor.Client

	chatRepository     repository.ChatRepository
	messagesRepository repository.MessagesRepository
	logRepository      repository.LogRepository
	chatService        service.ChatService
	chatImpl           *chat.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PgConfig() config.PgConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPgConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GrpcConfig() config.GrpcConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGrpcConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) AuthConfig() config.AuthConfig {
	if s.authConfig == nil {
		cfg, err := env.NewAuthConfig()
		if err != nil {
			log.Fatalf("failed to get authentication service config: %v", err)
		}

		s.authConfig = cfg
	}

	return s.authConfig
}

func (s *serviceProvider) TracingConfig() config.TracingConfig {
	if s.tracingConfig == nil {
		cfg, err := env.NewTracingConfig()
		if err != nil {
			log.Fatalf("failed to get tracing config: %v", err)
		}

		s.tracingConfig = cfg
	}

	return s.tracingConfig
}

func (s *serviceProvider) AuthClient() rpc.AuthClient {
	if s.authClient == nil {
		cfg := s.AuthConfig()
		creds, err := credentials.NewClientTLSFromFile(cfg.CertPath(), "")
		if err != nil {
			log.Fatalf("failed to get credentials of authentication service: %v", err)
		}

		conn, err := grpc.Dial(
			cfg.Address(),
			grpc.WithTransportCredentials(creds),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		)
		if err != nil {
			log.Fatalf("failed to connect to authentication service: %v", err)
		}

		s.authClient = rpcAuth.NewAuthClient(descAccess.NewAccessV1Client(conn))
	}

	return s.authClient
}

func (s *serviceProvider) InterceptorClient() *interceptor.Client {
	if s.interceptorClient == nil {
		s.interceptorClient = &interceptor.Client{
			Client: s.AuthClient(),
		}
	}
	return s.interceptorClient
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		c, err := pg.New(ctx, s.PgConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = c.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping database: %v", err)
		}

		closer.Add(c.Close)

		s.dbClient = c
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

func (s *serviceProvider) MessagesRepository(ctx context.Context) repository.MessagesRepository {
	if s.messagesRepository == nil {
		s.messagesRepository = messagesRepository.NewRepository(s.DBClient(ctx))
	}
	return s.messagesRepository
}

func (s *serviceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepository(s.DBClient(ctx))
	}
	return s.logRepository
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(s.ChatRepository(ctx), s.MessagesRepository(ctx), s.LogRepository(ctx), s.TxManager(ctx))
	}
	return s.chatService
}

func (s *serviceProvider) ChatImpl(ctx context.Context) *chat.Implementation {
	if s.chatImpl == nil {
		s.chatImpl = chat.NewImplementation(s.ChatService(ctx))
	}
	return s.chatImpl
}
