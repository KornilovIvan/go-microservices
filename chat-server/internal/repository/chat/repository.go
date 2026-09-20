package chat

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ivankornilov/chat-server/internal/dbtx"
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

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type repo struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) repository.ChatRepository {
	return &repo{pool: pool}
}

func (r *repo) querier(ctx context.Context) querier {
	if tx, ok := dbtx.Extract(ctx); ok {
		return tx
	}
	return r.pool
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

	var chatID int64
	err = r.querier(ctx).QueryRow(ctx, query, args...).Scan(&chatID)
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

	_, err = r.querier(ctx).Exec(ctx, query, args...)
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

	tag, err := r.querier(ctx).Exec(ctx, query, args...)
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

	_, err = r.querier(ctx).Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return model.ErrChatNotFound
		}
		return err
	}

	return nil
}
