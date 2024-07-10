package chat

import (
	"context"
	"errors"

	"github.com/AndreiMartynenko/chat-server/internal/model"
)

func (s *serv) SendMessage(ctx context.Context, chatID string, message *model.Message) error {
	s.mxChannels.RLock()
	chatChan, ok := s.channels[chatID]
	s.mxChannels.RUnlock()

	if !ok {
		return errors.New("chat not found")
	}

	// Save message in repository
	err := s.messagesRepository.Create(ctx, chatID, message)
	if err != nil {
		return errors.New("failed to save message")
	}

	chatChan <- message
	return nil
}
