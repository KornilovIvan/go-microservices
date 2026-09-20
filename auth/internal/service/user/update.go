package user

import (
	"context"

	"github.com/ivankornilov/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, user *model.UpdateUser) error {
	if user.Name == nil && user.Email == nil {
		return model.ErrNameOrEmailRequired
	}

	return s.userRepository.Update(ctx, user)
}
