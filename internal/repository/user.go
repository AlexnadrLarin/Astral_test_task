package repository

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"docs_storage/internal/models"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := builder.
		Insert("users").
		Columns("login", "password_hash").
		Values(u.Login, u.Password)

	sqlStr, args, err := q.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *UserRepo) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := builder.
		Select("id", "login", "password_hash", "created_at").
		From("users").
		Where(sq.Eq{"login": login}).
		Limit(1)

	sqlStr, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, sqlStr, args...)
	var u models.User
	if err := row.Scan(&u.ID, &u.Login, &u.Password, &u.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
