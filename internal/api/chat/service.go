package chat

import (
	"github.com/AndreiMartynenko/chat-server/internal/service"
	desc "github.com/AndreiMartynenko/chat-server/pkg/chat_v1"
)

// Implementation structure describes API layer.
type Implementation struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService
}

// NewImplementation creates new object of API layer.
func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{
		chatService: chatService,
	}
}
