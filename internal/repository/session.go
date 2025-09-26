package repository

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"docs_storage/internal/models"
)

type SessionRepo struct {
	db *pgxpool.Pool
}

func NewSessionRepo(db *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) Create(ctx context.Context, s *models.Session) error {
    builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

    q := builder.
        Insert("sessions").
        Columns("token", "user_id", "login").
        Values(s.Token, s.UserID, s.Login)

    sqlStr, args, err := q.ToSql()
    if err != nil {
        return err
    }

    _, err = r.db.Exec(ctx, sqlStr, args...)
    return err
}

func (r *SessionRepo) GetByToken(ctx context.Context, token string) (*models.Session, error) {
    builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

    q := builder.
        Select("token", "user_id", "login", "created_at").
        From("sessions").
        Where(sq.Eq{"token": token}).
        Limit(1)

    sqlStr, args, err := q.ToSql()
    if err != nil {
        return nil, err
    }

    row := r.db.QueryRow(ctx, sqlStr, args...)
    var s models.Session
    if err := row.Scan(&s.Token, &s.UserID, &s.Login, &s.CreatedAt); err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return &s, nil
}

func (r *SessionRepo) Delete(ctx context.Context, token string) error {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := builder.
		Delete("sessions").
		Where(sq.Eq{"token": token})

	sqlStr, args, err := q.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sqlStr, args...)
	return err
}
