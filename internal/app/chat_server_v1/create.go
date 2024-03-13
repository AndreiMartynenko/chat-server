package chat_server_v1

import (
	desc "github.com/AndreiMartynenko/chat-server/pkg/chat_v1"
	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type server struct {
	chat_v1.UnimplementedChatAPIServicesServer
	pool *pgxpool.Pool
}

// Create
func (srv *server) Create(ctx context.Context, req *desc.CreateNewChatRequest) (*desc.CreateNewChatResponse, error) {
	log.Printf("Create New Chat request received: %v", req)

	//For testing purposes
	builderInsert := sq.Insert("chats").
		PlaceholderFormat(sq.Dollar).
		Columns("id").
		Values(gofakeit.City(), gofakeit.Address().Street).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var chatID int64
	err = srv.pool.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		log.Fatalf("failed to insert chat: %v", err)
	}

	log.Printf("inserted chat with id: %d", chatID)

	return &chat_v1.CreateNewChatResponse{
		Id: chatID,
	}, nil

}
