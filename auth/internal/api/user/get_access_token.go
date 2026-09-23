package user

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ivankornilov/auth/internal/model"
	desc "github.com/ivankornilov/auth/pkg/auth_v1"
)

func (i *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	accessToken, err := i.authService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		if errors.Is(err, model.ErrInvalidRefreshToken) {
			return nil, status.Errorf(codes.Aborted, "invalid refresh token")
		}
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &desc.GetAccessTokenResponse{AccessToken: accessToken}, nil
}
