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

	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		id, err = s.userRepository.Create(ctx, info)
		if err != nil {
			return err
		}

		return s.logRepository.Create(ctx, &model.UserLog{
			UserID: id,
			Action: model.LogActionCreate,
		})
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}
