package chat

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ivankornilov/chat-server/internal/model"
	desc "github.com/ivankornilov/chat-server/pkg/chat_v1"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.chatService.Delete(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, model.ErrChatNotFound) {
			return nil, status.Errorf(codes.NotFound, "chat with id %d not found", req.GetId())
		}
		return nil, mapError(err)
	}

	return &emptypb.Empty{}, nil
}
