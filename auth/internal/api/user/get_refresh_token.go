package user

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ivankornilov/auth/internal/model"
	desc "github.com/ivankornilov/auth/pkg/auth_v1"
)

func (i *Implementation) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	refreshToken, err := i.authService.GetRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		if errors.Is(err, model.ErrInvalidRefreshToken) {
			return nil, status.Errorf(codes.Aborted, "invalid refresh token")
		}
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &desc.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
}
