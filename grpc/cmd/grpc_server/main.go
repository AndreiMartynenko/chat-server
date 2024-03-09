package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"

	"github.com/AndreiMartynenko/chat-server/grpc/pkg/chat_v1"
	"github.com/golang/protobuf/ptypes/empty"
)

const (
	dbDSN     = "host=localhost port=54321 dbname=chats user=chat-user password=chat-password sslmode=disable"
	grpcPort  = 50051
	dbTimeout = time.Second * 3
)

// var counts int64
var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	fmt.Println(configPath)
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	//grpcConfig
	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %s", err.Error())
	}

	//pgConfig
	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %s ", err.Error())
	}

	list, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgx.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to db: %s", err.Error())
	}

	defer func() {
		err = pool.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	s := grpc.NewServer()
	reflection.Register(s)
	chat_v1.RegisterChatAPIServicesServer(
		s,
		chatServerV1.NewChatServer(pool, sq.StatementBuilder.PlaceholderFormat(sq.Dollar)),
	)

	log.Printf("server listerning at %v", list.Addr())

	if err = s.Serve(list); err != nil {
		log.Fatalf("failed to serv %v", err)
	}
}

//type server struct {
//	db *sql.DB
//	chat_v1.UnimplementedChatAPIServicesServer
//}

/*
func main() {
	log.Println("Starting the chat service")

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
*/

func (srv *server) Create(ctx context.Context, req *chat_v1.CreateNewChatRequest) (*chat_v1.CreateNewChatResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	row := srv.db.QueryRowContext(ctx, "INSERT INTO chats DEFAULT VALUES RETURNING id")
	var id int64
	err := row.Scan(&id)
	if err != nil {
		log.Printf("Failed to get chat id: %v", err)
		return nil, err
	}

	log.Printf("Created New Chat with id: %d", id)

	return &chat_v1.CreateNewChatResponse{Id: id}, nil
}

// Delete
func (srv *server) Delete(ctx context.Context, req *chat_v1.DeleteChatRequest) (*chat_v1.DeleteChatResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	_, err := srv.db.ExecContext(ctx, "DELETE FROM chats WHERE id = $1", req.GetId())
	if err != nil {
		log.Printf("Failed to delete chat: %v", err)
		return nil, err
	}

	log.Printf("Deleted chat with id: %d", req.GetId())

	return &chat_v1.DeleteChatResponse{DeleteResponse: &empty.Empty{}}, nil
}

// func (srv *server) SendMessage(ctx context.Context, req *chat_v1.SendMessageRequest) (*chat_v1.SendMessageResponse, error) {
// 	log.Printf("SendMessageRequest received: %v", req)

// 	return &chat_v1.SendMessageResponse{SendMessageResponse: &empty.Empty{}}, nil
// }

// GetUserID
func (srv *server) getUserID(username string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var userID int64
	err := srv.db.QueryRowContext(ctx, "SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no user found with username %s", username)
		}
		return 0, err
	}

	return userID, nil
}

// SendMessage
func (srv *server) SendMessage(ctx context.Context, req *chat_v1.SendMessageRequest) (*chat_v1.SendMessageResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	// Get the chat_id and username from the request
	chatID := req.GetChatId()
	username := req.GetFrom()
	text := req.GetText()
	timestamp := req.GetTimestamp().AsTime()

	// Look up the user_id based on the username
	// This assumes you have a function getUserID that takes a username and returns a user_id
	userID, err := srv.getUserID(username)
	if err != nil {
		log.Printf("Failed to get user id: %v", err)
		return nil, err
	}

	_, err = srv.db.ExecContext(ctx, "INSERT INTO messages (chat_id, user_id, text, timestamp) VALUES ($1, $2, $3, $4)", chatID, userID, text, timestamp)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return nil, err
	}

	log.Printf("Sent message from user %d in chat %d", userID, chatID)

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
