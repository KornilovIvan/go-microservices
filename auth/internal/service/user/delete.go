package user

import (
	"context"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	return s.userRepository.Delete(ctx, id)
}
