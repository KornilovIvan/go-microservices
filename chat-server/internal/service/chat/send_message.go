package chat

import (
	"context"

	"github.com/ivankornilov/chat-server/internal/model"
)

func (s *serv) SendMessage(ctx context.Context, message *model.Message) error {
	if message.ChatID == 0 {
		return model.ErrChatIDRequired
	}
	if message.From == "" || message.Text == "" {
		return model.ErrFromAndTextRequired
	}

	return s.chatRepository.SendMessage(ctx, message)
}
