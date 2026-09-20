package user

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ivankornilov/auth/internal/model"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, model.ErrNameEmailPasswordRequired),
		errors.Is(err, model.ErrPasswordMismatch),
		errors.Is(err, model.ErrNameOrEmailRequired):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, model.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, model.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Errorf(codes.Internal, "%v", err)
	}
}
