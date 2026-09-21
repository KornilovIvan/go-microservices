package log

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/KornilovIvan/platform_common/pkg/db"
	"github.com/ivankornilov/auth/internal/model"
	"github.com/ivankornilov/auth/internal/repository"
)

const (
	tableName = "user_logs"

	userIDColumn = "user_id"
	actionColumn = "action"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.LogRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, userLog *model.UserLog) error {
	query, args, err := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(userIDColumn, actionColumn).
		Values(userLog.UserID, userLog.Action).
		ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "log_repository.Create",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}
