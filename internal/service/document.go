package service

import (
	"slices"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	models "docs_storage/internal/models"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrAccessDenied = errors.New("access denied")
)

type docsRepository interface {
	Save(ctx context.Context, doc *models.Document) error
	List(ctx context.Context, requesterLogin, login, key, value string, limit int) ([]models.Document, error)
	GetByID(ctx context.Context, id string) (*models.Document, error)
	Delete(ctx context.Context, id string) error
}

type fileStorage interface {
	Save(fileName string, data []byte) (string, error)
	Delete(fileName string) error
}

type sessionRepo interface {
	GetByToken(ctx context.Context, token string) (*models.Session, error)
}

type DocsService struct {
	docsRepo    docsRepository
	fileStorage fileStorage
	sessions    sessionRepo
}

func NewDocsService(docRepo docsRepository, fileStorage fileStorage, sessions sessionRepository) *DocsService {
	return &DocsService{
		docsRepo:    docRepo,
		fileStorage: fileStorage,
		sessions:    sessions,
	}
}

func (s *DocsService) Create(ctx context.Context, meta *models.Document, fileName string, fileData []byte, jsonData []byte, token string) (*models.Document, error) {
	session, err := s.sessions.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrAccessDenied
	}

	doc := &models.Document{
		ID:         uuid.New().String(),
		Name:       meta.Name,
		Mime:       meta.Mime,
		File:       meta.File,
		Public:     meta.Public,
		OwnerLogin: session.Login,
		Grant:      meta.Grant,
		CreatedAt:  time.Now(),
		JSONData:   jsonData,
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

func (s *DocsService) List(ctx context.Context, token, login, key, value string, limit int) ([]models.Document, error) {
	session, err := s.sessions.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrAccessDenied
	}

	return s.docsRepo.List(ctx, session.Login, login, key, value, limit)
}

func (s *DocsService) GetByID(ctx context.Context, id, token string) (*models.Document, error) {
	session, err := s.sessions.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrAccessDenied
	}

	doc, err := s.docsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, ErrNotFound
	}

	if doc.Public || doc.OwnerLogin == session.Login {
		return doc, nil
	}

	if slices.Contains(doc.Grant, session.Login) {
		return doc, nil
	}

	return nil, ErrAccessDenied
}

func (s *DocsService) Delete(ctx context.Context, id, token string) error {
	session, err := s.sessions.GetByToken(ctx, token)
	if err != nil {
		return err
	}
	if session == nil {
		return ErrAccessDenied
	}

	doc, err := s.docsRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if doc == nil {
		return ErrNotFound
	}

	if doc.OwnerLogin != session.Login {
		return ErrAccessDenied
	}

	if err := s.docsRepo.Delete(ctx, id); err != nil {
		return err
	}

	if doc.File && doc.FilePath != "" {
		_ = s.fileStorage.Delete(doc.FilePath)
	}

	return nil
}
