package converter

import (
	model "github.com/AndreiMartynenko/chat-server/internal/model"
	modelRepo "github.com/AndreiMartynenko/chat-server/internal/repository/messages/model"
)

// ToMessagesFromRepo converts repository layer model to structure of service layer.
func ToMessagesFromRepo(messages []*modelRepo.Message) []*model.Message {
	var res []*model.Message
	for _, m := range messages {
		res = append(res, &model.Message{
			From:      m.From,
			Text:      m.Text,
			Timestamp: m.Timestamp,
		})
	}
	return res
}
