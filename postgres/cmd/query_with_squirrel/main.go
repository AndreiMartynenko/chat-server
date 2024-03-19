package main

import (
	"database/sql"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const (
	dbDSN = "host=localhost port=54321 dbname=chats user=chat-user password=chat-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	// Create a new pool of connections to the database
	pool, err ;= pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()
}

// Make a query to add a record to the database chats
builderInsert := sq.Insert("chats").
	PlaceholderFormat(sq.Dollar).
	Columns("id").
	Values(randomID).
	Suffix("RETURNING id")

query, args, err := builderInsert.ToSql()
if err != nil {
	log.Fatalf("failed to build insert query: %v", err)
}

var chatID int64
err = pool.QueryRow(ctx, query, args...).Scan(&chatID)
if err != nil {
	log.Fatalf("failed to insert chat id : %v", err)
}

log.Printf("inserted chat with id: %d", chatID)

// Make a query to the database to get data from the chats
builderSelect := sq.Select("id").
	From("chats").
	PlaceholderFormat(sq.Dollar).
	OrderBy("id ASC").
	Limit(10)

query, args, err = builderSelect.ToSql()
if err != nil {
	log.Fatalf("failed to build select query: %v", err)
}

rows, err := pool.Query(ctx, query, args...)
if err != nil {
	log.Fatalf("failed to get chat : %v", err)
}

var id int64
var createdAt time.Time
var updatedAt sql.NullTime

for rows.Next() {
	err = rows.Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to scan chat : %v", err)
	}
	log.Printf("chat id: %d, created at: %v, updated at: %v", id, createdAt, updatedAt)
}

// Make a query to the database to update the chat
builderUpdate := sq.Update("chats").
	PlaceholderFormat(sq.Dollar).
	Set("updated_at", time.Now()).
	Where(sq.Eq{"id": chatID})