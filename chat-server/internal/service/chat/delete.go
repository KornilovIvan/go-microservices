package chat

import (
	"context"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	return s.chatRepository.Delete(ctx, id)
}
