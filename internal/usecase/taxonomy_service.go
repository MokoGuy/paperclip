package usecase

import (
	"context"

	"github.com/MokoGuy/paperclip/internal/domain"
)

type TaxonomyService struct {
	cache   domain.TaxonomyReader
	syncer  *SyncService
}

func NewTaxonomyService(cache domain.TaxonomyReader, syncer *SyncService) *TaxonomyService {
	return &TaxonomyService{cache: cache, syncer: syncer}
}

func (s *TaxonomyService) ListTags(ctx context.Context) ([]domain.Tag, error) {
	if err := s.syncer.SyncIfNeeded(ctx); err != nil {
		return nil, err
	}
	return s.cache.ListTags(ctx)
}

func (s *TaxonomyService) ListDocumentTypes(ctx context.Context) ([]domain.DocumentType, error) {
	if err := s.syncer.SyncIfNeeded(ctx); err != nil {
		return nil, err
	}
	return s.cache.ListDocumentTypes(ctx)
}

func (s *TaxonomyService) ListCorrespondents(ctx context.Context) ([]domain.Correspondent, error) {
	if err := s.syncer.SyncIfNeeded(ctx); err != nil {
		return nil, err
	}
	return s.cache.ListCorrespondents(ctx)
}
