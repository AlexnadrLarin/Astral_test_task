package service

import (
	"context"
	"errors"
	"fmt"
	"slices"
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

type cache interface {
	Get(ctx context.Context, key string) (any, bool)
	Set(ctx context.Context, key string, value any)
	Delete(ctx context.Context, key string)
	DeletePrefix(ctx context.Context, prefix string)
}

type DocsService struct {
	docsRepo    docsRepository
	fileStorage fileStorage
	sessions    sessionRepo
	cache       cache
}

func NewDocsService(docRepo docsRepository, fileStorage fileStorage, sessions sessionRepo, c cache) *DocsService {
	return &DocsService{
		docsRepo:    docRepo,
		fileStorage: fileStorage,
		sessions:    sessions,
		cache:       c,
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

	s.cache.DeletePrefix(ctx, fmt.Sprintf("list:%s", session.Login))
	s.cache.Set(ctx, fmt.Sprintf("doc:%s", doc.ID), doc)

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

	cacheKey := fmt.Sprintf("list:%s:%s:%s:%s:%d", session.Login, login, key, value, limit)
	if cached, ok := s.cache.Get(ctx, cacheKey); ok {
		if docs, ok := cached.([]models.Document); ok {
			return docs, nil
		}
	}

	docs, err := s.docsRepo.List(ctx, session.Login, login, key, value, limit)
	if err != nil {
		return nil, err
	}

	s.cache.Set(ctx, cacheKey, docs)
	return docs, nil
}

func (s *DocsService) GetByID(ctx context.Context, id, token string) (*models.Document, error) {
	session, err := s.sessions.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrAccessDenied
	}

	cacheKey := fmt.Sprintf("doc:%s", id)
	if cached, ok := s.cache.Get(ctx, cacheKey); ok {
		if doc, ok := cached.(*models.Document); ok {
			if doc.Public || doc.OwnerLogin == session.Login || slices.Contains(doc.Grant, session.Login) {
				return doc, nil
			}
			return nil, ErrAccessDenied
		}
	}

	doc, err := s.docsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, ErrNotFound
	}

	if doc.Public || doc.OwnerLogin == session.Login || slices.Contains(doc.Grant, session.Login) {
		s.cache.Set(ctx, cacheKey, doc)
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

	s.cache.Delete(ctx, fmt.Sprintf("doc:%s", id))
	s.cache.DeletePrefix(ctx, fmt.Sprintf("list:%s", session.Login))

	return nil
}
