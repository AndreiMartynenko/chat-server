package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AndreiMartynenko/chat-server/grpc/pkg/chat_v1"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// type Chat struct {
// 	id int
// }

const (
	dbDSN     = "host=localhost port=54321 dbname=chats user=chat-user password=chat-password sslmode=disable"
	grpcPort  = 50051
	dbTimeout = time.Second * 3
)

//var counts int64

type server struct {
	db *sql.DB
	chat_v1.UnimplementedChatAPIServicesServer
}

func main() {
	log.Println("Starting authentication service")

	ctx := context.Background()

	// Create a connection to the database
	con, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer con.Close(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()

	reflection.Register(srv)
	chat_v1.RegisterChatAPIServicesServer(srv, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func (srv *server) Create(ctx context.Context, req *chat_v1.CreateNewChatRequest) (*chat_v1.CreateNewChatResponse, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	//chats table
	row := srv.db.QueryRowContext(ctx, "INSERT INTO chats (title) VALUES ($1)")
	var id int64
	err := row.Scan(&id)
	if err != nil {
		log.Printf("Failed to get chat id: %v", err)
		return nil, err
	}

	log.Printf("Created New Chat with id: %d", id)

	return &chat_v1.CreateNewChatResponse{Id: id}, nil

}

func (srv *server) Delete(ctx context.Context, req *chat_v1.DeleteChatRequest) (*chat_v1.DeleteChatResponse, error) {
	log.Printf("Delete Chat request received: %v", req)

	return &chat_v1.DeleteChatResponse{DeleteResponse: &empty.Empty{}}, nil

}

func (srv *server) SendMessage(ctx context.Context, req *chat_v1.SendMessageRequest) (*chat_v1.SendMessageResponse, error) {
	log.Printf("SendMessageRequest received: %v", req)

	return &chat_v1.SendMessageResponse{SendMessageResponse: &empty.Empty{}}, nil
}

// Creating a connection to DB

// func openDB(dsn string) (*sql.DB, error) {
// 	db, err := sql.Open("pgx", dsn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return db, nil
// }
