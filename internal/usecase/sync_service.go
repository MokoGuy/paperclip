package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/MokoGuy/paperclip/internal/domain"
	apirepo "github.com/MokoGuy/paperclip/internal/repository/api"
	sqliterepo "github.com/MokoGuy/paperclip/internal/repository/sqlite"
)

const syncThreshold = 24 * time.Hour

type SyncService struct {
	apiTaxonomy domain.TaxonomyReader
	apiDoc      *apirepo.DocumentRepository
	cacheTax    *sqliterepo.TaxonomyRepository
	cacheDoc    *sqliterepo.DocumentRepository
	syncRepo    *sqliterepo.SyncRepository
}

func NewSyncService(
	apiTaxonomy domain.TaxonomyReader,
	apiDoc *apirepo.DocumentRepository,
	cacheTax *sqliterepo.TaxonomyRepository,
	cacheDoc *sqliterepo.DocumentRepository,
	syncRepo *sqliterepo.SyncRepository,
) *SyncService {
	return &SyncService{
		apiTaxonomy: apiTaxonomy,
		apiDoc:      apiDoc,
		cacheTax:    cacheTax,
		cacheDoc:    cacheDoc,
		syncRepo:    syncRepo,
	}
}

func (s *SyncService) SyncIfNeeded(ctx context.Context) error {
	lastSync, err := s.syncRepo.GetLastSync(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sync state: %w", err)
	}

	if lastSync.IsZero() || time.Since(lastSync) > syncThreshold {
		return s.Sync(ctx)
	}
	return nil
}

func (s *SyncService) NeedsSync(ctx context.Context) (bool, error) {
	lastSync, err := s.syncRepo.GetLastSync(ctx)
	if err != nil {
		return true, nil
	}
	return lastSync.IsZero() || time.Since(lastSync) > syncThreshold, nil
}

func (s *SyncService) Sync(ctx context.Context) error {
	if err := s.syncTaxonomy(ctx); err != nil {
		return err
	}

	if err := s.syncDocuments(ctx); err != nil {
		return err
	}

	return s.syncRepo.SetLastSync(ctx, time.Now())
}

func (s *SyncService) syncDocuments(ctx context.Context) error {
	docs, err := s.apiDoc.FetchAllDocuments(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch documents: %w", err)
	}
	for _, d := range docs {
		if err := s.cacheDoc.UpsertDocument(ctx, d); err != nil {
			return fmt.Errorf("failed to cache document %d: %w", d.ID, err)
		}
	}
	return nil
}

func (s *SyncService) syncTaxonomy(ctx context.Context) error {
	tags, err := s.apiTaxonomy.ListTags(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch tags: %w", err)
	}
	for _, t := range tags {
		if err := s.cacheTax.UpsertTag(ctx, t); err != nil {
			return fmt.Errorf("failed to cache tag %q: %w", t.Name, err)
		}
	}

	types, err := s.apiTaxonomy.ListDocumentTypes(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch document types: %w", err)
	}
	for _, dt := range types {
		if err := s.cacheTax.UpsertDocumentType(ctx, dt); err != nil {
			return fmt.Errorf("failed to cache document type %q: %w", dt.Name, err)
		}
	}

	correspondents, err := s.apiTaxonomy.ListCorrespondents(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch correspondents: %w", err)
	}
	for _, c := range correspondents {
		if err := s.cacheTax.UpsertCorrespondent(ctx, c); err != nil {
			return fmt.Errorf("failed to cache correspondent %q: %w", c.Name, err)
		}
	}

	return nil
}
