package user

import (
	"context"

	"github.com/ivankornilov/auth/internal/model"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		err := s.userRepository.Delete(ctx, id)
		if err != nil {
			return err
		}

		return s.logRepository.Create(ctx, &model.UserLog{
			UserID: id,
			Action: model.LogActionDelete,
		})
	})
}
