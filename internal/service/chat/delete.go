package chat

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/AndreiMartynenko/chat-server/internal/model"
)

func (s *serv) Delete(ctx context.Context, id string) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.messagesRepository.DeleteChat(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.chatRepository.Delete(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, &model.Log{
			Text: fmt.Sprintf("Deleted chat with id: %v", id),
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		log.Print(err)
		return errors.New("failed to delete chat")
	}

	// Delete channel associated with chat
	delete(s.channels, id)
	return nil
}
