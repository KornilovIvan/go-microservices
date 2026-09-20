package chat

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ivankornilov/chat-server/internal/model"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, model.ErrChatIDRequired),
		errors.Is(err, model.ErrFromAndTextRequired):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, model.ErrChatNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Errorf(codes.Internal, "%v", err)
	}
}
