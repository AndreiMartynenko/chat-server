package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/AndreiMartynenko/chat-server/config/pkg/chat_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

const grpcPort = 50051

type server struct {
	chat_v1.UnimplementedChatAPIServicesServer
}

func (srv *server) Create(ctx context.Context, req *chat_v1.CreateNewChatRequest) (*chat_v1.CreateNewChatResponse, error) {
	log.Printf("Create New Chat request received: %v", req)

	//For testing purposes
	// response := &chat_v1.CreateNewChatResponse{
	// 	Id: 1345,
	// }
	// return response, nil
	return &chat_v1.CreateNewChatResponse{
		Id: gofakeit.Int64(),
	}, nil

}

func (srv *server) Delete(ctx context.Context, req *chat_v1.DeleteChatRequest) (*chat_v1.DeleteChatResponse, error) {
	log.Printf("Delete Chat request received: %v", req)

	return &chat_v1.DeleteChatResponse{DeleteResponse: &empty.Empty{}}, nil

}

func (srv *server) SendMessage(ctx context.Context, req *chat_v1.SendMessageRequest) (*chat_v1.SendMessageResponse, error) {
	log.Printf("SendMessageRequest received: %v", req)

	return &chat_v1.SendMessageResponse{SendMessageResponse: &empty.Empty{}}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()

	// Reflection in this context allows gRPC clients to query information
	// about the gRPC server's services dynamically at runtime.
	// It enables tools like gRPC's command-line interface (grpc_cli)
	// and gRPC's web-based GUI (grpcui) to inspect the server's
	// services and make RPC calls without needing to know
	// the specifics of each service beforehand.

	reflection.Register(srv)
	chat_v1.RegisterChatAPIServicesServer(srv, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
