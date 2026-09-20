package user

import (
	"context"

	"github.com/ivankornilov/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, user *model.UpdateUser) error {
	if user.Name == nil && user.Email == nil {
		return model.ErrNameOrEmailRequired
	}

	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		err := s.userRepository.Update(ctx, user)
		if err != nil {
			return err
		}

		return s.logRepository.Create(ctx, &model.UserLog{
			UserID: user.ID,
			Action: model.LogActionUpdate,
		})
	})
}
