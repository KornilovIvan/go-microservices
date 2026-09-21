package service

//go:generate minimock -i github.com/ivankornilov/auth/internal/service.UserService -o ./mocks -s _minimock.go

import (
	"context"

	"github.com/ivankornilov/auth/internal/model"
)

type UserService interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.UpdateUser) error
	Delete(ctx context.Context, id int64) error
}
