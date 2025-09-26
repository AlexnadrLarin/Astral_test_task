package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	models "docs_storage/internal/models"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrAccessDenied = errors.New("access denied")
)

type DocumentRepo struct {
	db *pgxpool.Pool
}

func NewDocsRepo(db *pgxpool.Pool) *DocumentRepo {
	return &DocumentRepo{db: db}
}

func (r *DocumentRepo) Save(ctx context.Context, doc *models.Document) error {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := builder.
		Insert("documents").
		Columns(
			"id", "name", "mime", "file", "public",
			"owner_login", "grant_list",
			"created_at", "json_data", "file_path",
		).
		Values(
			doc.ID, doc.Name, doc.Mime, doc.File, doc.Public,
			doc.OwnerLogin, doc.Grant,
			doc.CreatedAt, doc.JSONData, doc.FilePath,
		)

	sqlStr, args, err := q.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *DocumentRepo) List(ctx context.Context, requesterLogin, login, key, value string, limit int) ([]models.Document, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := builder.
		Select("id", "name", "mime", "file", "public", "owner_login", "grant_list", "created_at", "json_data", "file_path").
		From("documents")

	if login != "" {
		q = q.Where(sq.Eq{"owner_login": login})
	} else {
		q = q.Where(sq.Or{
			sq.Eq{"owner_login": requesterLogin},
			sq.Eq{"public": true},
			sq.Expr("? = ANY(grant_list)", requesterLogin),
		})
	}

	allowedKeys := map[string]bool{
		"id":         true,
		"name":       true,
		"mime":       true,
		"file":       true,
		"public":     true,
		"created_at": true,
	}

	if key != "" && value != "" {
		if !allowedKeys[key] {
			return nil, fmt.Errorf("invalid filter key: %s", key)
		}

		if key == "file" || key == "public" {
			b, err := strconv.ParseBool(value)
			if err != nil {
				return nil, fmt.Errorf("invalid boolean value for %s: %w", key, err)
			}
			q = q.Where(sq.Eq{key: b})
		} else if key == "created_at" {
			var parsed time.Time
			var err error
			formats := []string{
				time.RFC3339, "2006-01-02", "2006-01-02 15:04:05",
			}
			for _, f := range formats {
				parsed, err = time.Parse(f, value)
				if err == nil {
					break
				}
			}
			if err == nil {
				q = q.Where(sq.Eq{key: parsed})
			} else {
				q = q.Where(sq.Eq{key: value})
			}
		} else {
			q = q.Where(sq.Eq{key: value})
		}
	}

	q = q.OrderBy("name", "created_at")
	if limit > 0 {
		q = q.Limit(uint64(limit))
	}

	sqlStr, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []models.Document
	for rows.Next() {
		var d models.Document
		if err := rows.Scan(
			&d.ID, &d.Name, &d.Mime, &d.File, &d.Public,
			&d.OwnerLogin, &d.Grant, &d.CreatedAt, &d.JSONData, &d.FilePath,
		); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}

	return docs, nil
}

func (r *DocumentRepo) GetByID(ctx context.Context, id string) (*models.Document, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := builder.
		Select("id", "name", "mime", "file", "public", "owner_login", "grant_list", "created_at", "json_data", "file_path").
		From("documents").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sqlStr, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, sqlStr, args...)
	var d models.Document
	if err := row.Scan(
		&d.ID, &d.Name, &d.Mime, &d.File, &d.Public,
		&d.OwnerLogin, &d.Grant, &d.CreatedAt, &d.JSONData, &d.FilePath,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

func (r *DocumentRepo) Delete(ctx context.Context, id string) error {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := builder.Delete("documents").Where(sq.Eq{"id": id})

	sqlStr, args, err := q.ToSql()
	if err != nil {
		return err
	}

	cmd, err := r.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
