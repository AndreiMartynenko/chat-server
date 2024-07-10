package repository

import (
	"context"

	"github.com/AndreiMartynenko/chat-server/internal/model"
)

// ChatRepository is the interface for chat info repository communication.
type ChatRepository interface {
	Create(ctx context.Context, chat *model.Chat) (string, error)
	Delete(ctx context.Context, id string) error
	GetChats(ctx context.Context) ([]string, error)
}

// MessagesRepository is the interface for messages info repository communication.
type MessagesRepository interface {
	Create(ctx context.Context, chatID string, message *model.Message) error
	GetMessages(ctx context.Context, chatID string) ([]*model.Message, error)
	DeleteChat(ctx context.Context, chatID string) error
}

// LogRepository is the interface for transaction log repository communication.
type LogRepository interface {
	Log(ctx context.Context, log *model.Log) error
}
