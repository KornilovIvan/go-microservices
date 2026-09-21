package chat

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"

	"github.com/KornilovIvan/platform_common/pkg/db"
	"github.com/ivankornilov/chat-server/internal/model"
	"github.com/ivankornilov/chat-server/internal/repository"
)

const (
	chatsTable     = "chats"
	chatUsersTable = "chat_users"
	messagesTable  = "messages"

	idColumn        = "id"
	chatIDColumn    = "chat_id"
	usernameColumn  = "username"
	fromUserColumn  = "from_user"
	textColumn      = "text"
	createdAtColumn = "created_at"
	sentAtColumn    = "sent_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context) (int64, error) {
	query, args, err := sq.Insert(chatsTable).
		PlaceholderFormat(sq.Dollar).
		Columns(createdAtColumn).
		Values(time.Now()).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "chat_repository.Create",
		QueryRaw: query,
	}

	var chatID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

func (r *repo) AddUser(ctx context.Context, chatID int64, username string) error {
	query, args, err := sq.Insert(chatUsersTable).
		PlaceholderFormat(sq.Dollar).
		Columns(chatIDColumn, usernameColumn).
		Values(chatID, username).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_repository.AddUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	query, args, err := sq.Delete(chatsTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_repository.Delete",
		QueryRaw: query,
	}

	tag, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrChatNotFound
	}

	return nil
}

func (r *repo) SendMessage(ctx context.Context, message *model.Message) error {
	query, args, err := sq.Insert(messagesTable).
		PlaceholderFormat(sq.Dollar).
		Columns(chatIDColumn, fromUserColumn, textColumn, sentAtColumn).
		Values(message.ChatID, message.From, message.Text, message.SentAt).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_repository.SendMessage",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return model.ErrChatNotFound
		}
		return err
	}

	return nil
}
