package chat

import (
	"errors"
	"log"

	"github.com/AndreiMartynenko/chat-server/internal/converter"
	"github.com/AndreiMartynenko/chat-server/internal/model"
)

func (s *serv) Connect(chatID string, username string, stream model.Stream) error {
	s.mxChannels.RLock()
	chatChan, ok := s.channels[chatID]
	s.mxChannels.RUnlock()

	if !ok {
		return errors.New("chat not found")
	}

	// Init streams for chat, if they don't exist
	s.mxChat.Lock()
	if _, okChat := s.chats[chatID]; !okChat {
		s.chats[chatID] = &chat{
			streams: make(map[string]model.Stream),
		}
	}
	s.mxChat.Unlock()

	// Set stream for user
	s.chats[chatID].m.Lock()
	s.chats[chatID].streams[username] = stream
	s.chats[chatID].m.Unlock()

	if err := s.loadHistory(chatID, stream); err != nil {
		// If history not loaded, there's no problem, you can still send messages
		log.Printf("failed to load history: %v", err)
	}

	for {
		select {
		case msg, okCh := <-chatChan:
			// Check if channel is closed
			if !okCh {
				return nil
			}

			// Send message for everyone in chat
			for _, st := range s.chats[chatID].streams {
				if err := st.Send(converter.ToMessageFromService(msg)); err != nil {
					return err
				}
			}
		case <-stream.Context().Done():
			// Delete stream for user when context is dead
			s.chats[chatID].m.Lock()
			delete(s.chats[chatID].streams, username)
			s.chats[chatID].m.Unlock()
			return nil
		}
	}
}

func (s *serv) loadHistory(chatID string, stream model.Stream) error {
	messages, err := s.messagesRepository.GetMessages(stream.Context(), chatID)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		if err := stream.Send(converter.ToMessageFromService(msg)); err != nil {
			return err
		}
	}
	return nil
}
