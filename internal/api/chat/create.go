package chat

import (
	"context"

	"github.com/AndreiMartynenko/chat-server/internal/converter"
	desc "github.com/AndreiMartynenko/chat-server/pkg/chat_v1"
)

// Create is used for creating new chat.
func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.chatService.Create(ctx, converter.ToChatFromDesc(req.GetChat()))
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
