package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"

	desc "github.com/AndreiMartynenko/chat-server/grpc/pkg/chat_server_v1"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedChatAPIServicesV1Server
}

func (srv *server) Create(ctx context.Context, req *desc.CreateNewChatRequest) (*desc.CreateNewChatResponse, error) {
	log.Printf("Create New Chat request received: %v", req)

	//For testing purposes
	// response := &chat_server_v1.CreateNewChatResponse{
	// 	Id: 1345,
	// }
	// return response, nil
	return &desc.CreateNewChatResponse{
		Id: gofakeit.Int64(),
	}, nil

}

func (srv *server) Delete(ctx context.Context, req *desc.DeleteChatRequest) (emptypb.Empty, error) {
	log.Printf("Delete Chat request received: %v", req)
	return nil

}

//func (srv *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*desc.SendMessageResponse, error) {
//	log.Printf("SendMessageRequest received: %v", req)
//
//	return &desc.SendMessageResponse{SendMessageResponse: &empty.Empty{}}, nil
//}

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
	desc.RegisterChatAPIServicesV1Server(srv, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
