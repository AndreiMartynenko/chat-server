package service

import (
	"context"

	"github.com/AndreiMartynenko/chat-server/internal/model"
)

// ChatService is the interface for service communication.
type ChatService interface {
	Create(ctx context.Context, chat *model.Chat) (string, error)
	Delete(ctx context.Context, id string) error
	SendMessage(ctx context.Context, chatID string, message *model.Message) error
	Connect(chatID string, username string, stream model.Stream) error
	InitChannels(ctx context.Context) error
}
