package user

import (
	"context"
	"log"

	"github.com/ivankornilov/auth/internal/converter"
	desc "github.com/ivankornilov/auth/pkg/auth_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.userService.Create(ctx, converter.ToUserInfoFromDesc(req))
	if err != nil {
		return nil, mapError(err)
	}

	log.Printf("created user id=%d", id)

	return &desc.CreateResponse{Id: id}, nil
}
