package chat

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/AndreiMartynenko/chat-server/internal/model"
)

func (s *serv) Create(ctx context.Context, chat *model.Chat) (string, error) {
	var id string

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.chatRepository.Create(ctx, chat)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, &model.Log{
			Text: fmt.Sprintf("Created chat with id: %v", id),
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		log.Print(err)
		return "", errors.New("failed to create chat")
	}

	// Create buffered channel for new chat
	s.channels[id] = make(chan *model.Message, messagesBuffer)

	return id, nil
}
