package chat

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ivankornilov/chat-server/internal/converter"
	"github.com/ivankornilov/chat-server/internal/model"
	desc "github.com/ivankornilov/chat-server/pkg/chat_v1"
)

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	err := i.chatService.SendMessage(ctx, converter.ToMessageFromDesc(req))
	if err != nil {
		if errors.Is(err, model.ErrChatNotFound) {
			return nil, status.Errorf(codes.NotFound, "chat with id %d not found", req.GetChatId())
		}
		return nil, mapError(err)
	}

	return &emptypb.Empty{}, nil
}
