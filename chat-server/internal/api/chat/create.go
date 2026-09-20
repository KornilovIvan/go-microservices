package chat

import (
	"context"
	"log"

	desc "github.com/ivankornilov/chat-server/pkg/chat_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.chatService.Create(ctx, req.GetUsernames())
	if err != nil {
		return nil, mapError(err)
	}

	log.Printf("created chat id=%d usernames=%v", id, req.GetUsernames())

	return &desc.CreateResponse{Id: id}, nil
}
