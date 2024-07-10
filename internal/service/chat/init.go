package chat

import (
	"context"
	"errors"

	"github.com/AndreiMartynenko/chat-server/internal/model"
)

func (s *serv) InitChannels(ctx context.Context) error {
	// Get chats from repository
	ids, err := s.chatRepository.GetChats(ctx)
	if err != nil {
		return errors.New("failed to init existing chats")
	}

	// Fill chats and channels for already existing chats
	for _, id := range ids {
		s.channels[id] = make(chan *model.Message, messagesBuffer)
	}

	return nil
}
