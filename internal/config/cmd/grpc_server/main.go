package main

import (
	"flag"
	desc "github.com/AndreiMartynenko/chat-server/grpc/pkg/chat_server_v1"
	"github.com/AndreiMartynenko/chat-server/internal/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedChatAPIServicesV1Server
	pool *pgxpool.Pool
}

func (srv *server) Create(ctx context.Context, req *desc.CreateNewChatRequest) (*desc.CreateNewChatResponse, error) {

}

func (srv *server) Delete(ctx context.Context, req *desc.DeleteChatRequest) (*emptypb.Empty, error) {

}

func main() {
	flag.Parse()

	ctx := context.Background()
	log.Println("Starting server...")

	// Load config
	//cfg := config.MustLoad()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: #{err}")

	}
	// Connect to database
	pool, err := pgxpool.Connect(context.Background(), cfg.Database.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Create server
	srv := &server{
		pool: pool,
	}

	// Create gRPC server
	s := grpc.NewServer()
	desc.RegisterChatAPIServicesV1Server(s, srv)

	// Start server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
