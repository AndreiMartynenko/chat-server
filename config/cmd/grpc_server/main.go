package main

import (
	"flag"
	desc "github.com/AndreiMartynenko/chat-server/grpc/pkg/chat_server_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
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
