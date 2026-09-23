package user

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ivankornilov/auth/internal/model"
	desc "github.com/ivankornilov/auth/pkg/auth_v1"
)

func (i *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	refreshToken, err := i.authService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &desc.LoginResponse{RefreshToken: refreshToken}, nil
}
