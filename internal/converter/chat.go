package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/AndreiMartynenko/chat-server/internal/model"
	desc "github.com/AndreiMartynenko/chat-server/pkg/chat_v1"
)

// ToChatFromService converts service layer model to structure of API layer.
func ToChatFromService(chat *model.Chat) *desc.Chat {
	return &desc.Chat{
		Usernames: chat.Usernames,
	}
}

// ToChatFromDesc converts structure of API layer to service layer model.
func ToChatFromDesc(chat *desc.Chat) *model.Chat {
	return &model.Chat{
		Usernames: chat.Usernames,
	}
}

// ToMessageFromDesc converts structure of API layer to service layer model.
func ToMessageFromDesc(message *desc.Message) *model.Message {
	return &model.Message{
		From:      message.From,
		Text:      message.Text,
		Timestamp: message.Timestamp.AsTime(),
	}
}

// ToStreamFromDesc converts interface of API layer to service layer interface.
func ToStreamFromDesc(stream desc.ChatV1_ConnectServer) model.Stream {
	return stream.(model.Stream)
}

// ToMessageFromService converts service layer model to structure of API layer.
func ToMessageFromService(message *model.Message) *desc.Message {
	return &desc.Message{
		From:      message.From,
		Text:      message.Text,
		Timestamp: timestamppb.New(message.Timestamp),
	}
}
