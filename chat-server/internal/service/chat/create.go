package chat

import (
	"context"
)

func (s *serv) Create(ctx context.Context, usernames []string) (int64, error) {
	var chatID int64

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		id, err := s.chatRepository.Create(ctx)
		if err != nil {
			return err
		}
		chatID = id

		for _, username := range usernames {
			if username == "" {
				continue
			}
			if err = s.chatRepository.AddUser(ctx, chatID, username); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return chatID, nil
}
