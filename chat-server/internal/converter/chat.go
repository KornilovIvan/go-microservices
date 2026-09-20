package converter

import (
	"time"

	"github.com/ivankornilov/chat-server/internal/model"
	desc "github.com/ivankornilov/chat-server/pkg/chat_v1"
)

func ToMessageFromDesc(req *desc.SendMessageRequest) *model.Message {
	sentAt := time.Now()
	if req.GetTimestamp() != nil && req.GetTimestamp().IsValid() {
		sentAt = req.GetTimestamp().AsTime()
	}

	return &model.Message{
		ChatID: req.GetChatId(),
		From:   req.GetFrom(),
		Text:   req.GetText(),
		SentAt: sentAt,
	}
}
