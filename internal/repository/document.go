package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	models "docs_storage/internal/models"
)

type DocumentRepo struct {
	db *pgxpool.Pool
}

func NewDocsRepo(db *pgxpool.Pool) *DocumentRepo {
	return &DocumentRepo{db: db}
}

func (r *DocumentRepo) Save(ctx context.Context, doc *models.Document) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO documents (id, name, mime, file, public, token, grant_list, created_at, json_data, file_path)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		doc.ID, doc.Name, doc.Mime, doc.File, doc.Public, doc.Token, doc.Grant, doc.CreatedAt, doc.JSONData, doc.FilePath,
	)
	return err
}
