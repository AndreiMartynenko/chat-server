package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4"
)

const (
	dbDSN = "host=localhost port=54321 dbname=chat user=chat-user password=chat-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	// Create a connection to the database
	con, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer con.Close(ctx)

	// Making a request to insert a record into the chat table
	res, err := con.Exec(ctx, "INSERT INTO chat (title, body) VALUES ($1, $2)", gofakeit.City(), gofakeit.Address().Street)
	if err != nil {
		log.Fatalf("failed to insert chat: %v", err)
	}

	log.Printf("inserted %d rows", res.RowsAffected())

	// We make a request to select records from the chat table
	rows, err := con.Query(ctx, "SELECT id, title, body, created_at, updated_at FROM chat")
	if err != nil {
		log.Fatalf("failed to select chat_users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title, body string
		var createdAt time.Time
		var updatedAt sql.NullTime

		err = rows.Scan(&id, &title, &body, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("failed to scan chat: %v", err)
		}

		log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", id, title, body, createdAt, updatedAt)
	}
}
