package access

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/ivankornilov/auth/internal/model"
	"github.com/ivankornilov/auth/internal/utils"
)

const authPrefix = "Bearer "

func (s *serv) Check(ctx context.Context, endpointAddress string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return model.ErrAuthHeaderMissing
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return model.ErrAuthHeaderMissing
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return model.ErrAuthHeaderInvalid
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	claims, err := utils.VerifyToken(accessToken, s.accessSecretKey)
	if err != nil {
		return model.ErrAccessTokenInvalid
	}

	role, ok := s.accessibleRoles[endpointAddress]
	if !ok {
		return nil
	}

	if role == claims.Role {
		return nil
	}

	return model.ErrAccessDenied
}
