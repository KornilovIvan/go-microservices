package access

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ivankornilov/auth/internal/model"
	desc "github.com/ivankornilov/auth/pkg/access_v1"
)

func (i *Implementation) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {
	err := i.accessService.Check(ctx, req.GetEndpointAddress())
	if err != nil {
		switch {
		case errors.Is(err, model.ErrAuthHeaderMissing),
			errors.Is(err, model.ErrAuthHeaderInvalid),
			errors.Is(err, model.ErrAccessTokenInvalid):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		case errors.Is(err, model.ErrAccessDenied):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	}

	return &emptypb.Empty{}, nil
}
