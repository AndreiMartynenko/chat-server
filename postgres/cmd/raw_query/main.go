package main

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"log"

	"github.com/jackc/pgx/v4"
)

const (
	dbDSN = "host=localhost port=54321 dbname=chats user=chat-user password=chat-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	// Create a new connection to the database
	con, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer con.Close(ctx)

	// Generate a random chat id
	var randomID uint64
	err = binary.Read(rand.Reader, binary.BigEndian, &randomID)
	if err != nil {
		log.Fatalf("failed to generate random number: %v", err)
	}

	// Make a query to the database
	res, err := con.Exec(ctx, "INSERT INTO chats (id) VALUES ($1)", randomID)
	if err != nil {
		log.Fatalf("failed to insert chat id : %v", err)
	}

	log.Printf("inserted chat with id: %d", randomID)
	log.Printf("inserted %d rows", res.RowsAffected())

	// Make a query to the database to get data from the chat
	rows, err := con.Query(ctx, "SELECT id FROM chats")
	if err != nil {
		log.Fatalf("failed to get chat : %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			log.Fatalf("failed to scan chat : %v", err)
		}
		log.Printf("chat id: %d", id)
	}
}
