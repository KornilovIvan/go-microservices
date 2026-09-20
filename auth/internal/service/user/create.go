package user

import (
	"context"

	"github.com/ivankornilov/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.UserInfo) (int64, error) {
	if info.Name == "" || info.Email == "" || info.Password == "" {
		return 0, model.ErrNameEmailPasswordRequired
	}
	if info.Password != info.PasswordConfirm {
		return 0, model.ErrPasswordMismatch
	}

	return s.userRepository.Create(ctx, info)
}
