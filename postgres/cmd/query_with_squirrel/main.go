package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const (
	dbDSN = "host=localhost port=54321 dbname=chats user=chat-user password=chat-password sslmode=disable"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Set a timeout of 5 seconds
	defer cancel()

	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	randomID, err := generateRandomID()
	if err != nil {
		log.Fatalf("failed to generate random ID: %v", err)
	}

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

	builderSelect := sq.Select("id", "created_at", "updated_at").
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
	defer rows.Close()

	for rows.Next() {
		var id int64
		var createdAt time.Time
		var updatedAt sql.NullTime
		if err := rows.Scan(&id, &createdAt, &updatedAt); err != nil {
			log.Fatalf("failed to scan chat : %v", err)
		}
		log.Printf("chat id: %d, created at: %v, updated at: %v", id, createdAt, updatedAt)
	}

	builderUpdate := sq.Update("chats").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": chatID})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build update query: %v", err)
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update chat : %v", err)
	}

	log.Printf("updated %d rows", res.RowsAffected())

	builderSelect = sq.Select("id, created_at, updated_at").
		From("chats").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": chatID}).
		Limit(1)

	query, args, err = builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build select query: %v", err)
	}

	var createdAt time.Time
	var updatedAt sql.NullTime
	err = pool.QueryRow(ctx, query, args...).Scan(&chatID, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to select chats : %v", err)
	}
	log.Printf("chat id: %d, created at: %v, updated at: %v", chatID, createdAt, updatedAt)
}

func generateRandomID() (uint64, error) {
	var randomID uint64
	err := binary.Read(rand.Reader, binary.BigEndian, &randomID)
	if err != nil {
		return 0, err
	}
	return randomID, nil
}
