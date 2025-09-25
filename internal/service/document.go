package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	models "docs_storage/internal/models"
)

type docsRepository interface {
	Save(ctx context.Context, doc *models.Document) error
}

type fileStorage interface {
	Save(fileName string, data []byte) (string, error)
}

type DocsService struct {
	docsRepo    docsRepository
	fileStorage fileStorage
}

func NewDocsService(docRepo docsRepository, fileStorage fileStorage) *DocsService {
	return &DocsService{docsRepo: docRepo, fileStorage: fileStorage}
}

func (s *DocsService) Create(ctx context.Context, meta *models.Document, fileName string, fileData []byte, jsonData []byte) (*models.Document, error) {
	doc := &models.Document{
		ID:        uuid.New().String(),
		Name:      meta.Name,
		Mime:      meta.Mime,
		File:      meta.File,
		Public:    meta.Public,
		Grant:     meta.Grant,
		CreatedAt: time.Now(),
		JSONData:  jsonData,
	}

	if meta.File && len(fileData) > 0 {
		path, err := s.fileStorage.Save(fileName, fileData)
		if err != nil {
			return nil, err
		}
		doc.FilePath = path
	}

	if err := s.docsRepo.Save(ctx, doc); err != nil {
		return nil, err
	}

	return doc, nil
}
