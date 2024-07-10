package chat

import (
	"log"

	"github.com/AndreiMartynenko/chat-server/internal/converter"
	desc "github.com/AndreiMartynenko/chat-server/pkg/chat_v1"
)

// Connect is used for connecting to a chat.
func (i *Implementation) Connect(req *desc.ConnectRequest, stream desc.ChatV1_ConnectServer) error {
	err := i.chatService.Connect(req.GetChatId(), req.GetUsername(), converter.ToStreamFromDesc(stream))
	log.Println(err)
	return err
}
