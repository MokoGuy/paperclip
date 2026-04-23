package usecase

import (
	"context"
	"strings"

	"github.com/sahilm/fuzzy"

	"github.com/MokoGuy/paperclip/internal/domain"
	apirepo "github.com/MokoGuy/paperclip/internal/repository/api"
	sqliterepo "github.com/MokoGuy/paperclip/internal/repository/sqlite"
)

type SearchService struct {
	cacheDoc  *sqliterepo.DocumentRepository
	cacheTax  *sqliterepo.TaxonomyRepository
	apiDoc    *apirepo.DocumentRepository
	syncer    *SyncService
}

func NewSearchService(
	cacheDoc *sqliterepo.DocumentRepository,
	cacheTax *sqliterepo.TaxonomyRepository,
	apiDoc *apirepo.DocumentRepository,
	syncer *SyncService,
) *SearchService {
	return &SearchService{
		cacheDoc: cacheDoc,
		cacheTax: cacheTax,
		apiDoc:   apiDoc,
		syncer:   syncer,
	}
}

func (s *SearchService) Search(ctx context.Context, filters domain.SearchFilters, noCache bool) ([]domain.Document, error) {
	if noCache {
		return s.apiDoc.ListDocuments(ctx, filters)
	}

	if err := s.syncer.SyncIfNeeded(ctx); err != nil {
		return nil, err
	}

	resolved, err := s.resolveFilterNames(ctx, filters)
	if err != nil {
		return nil, err
	}

	docs, err := s.cacheDoc.SearchDocuments(ctx, resolved)
	if err != nil {
		return nil, err
	}

	// If SQL LIKE found no results but we have a text query, try fuzzy on titles
	if len(docs) == 0 && resolved.Query != "" {
		fuzzyFilters := resolved
		fuzzyFilters.Query = ""
		allDocs, err := s.cacheDoc.SearchDocuments(ctx, fuzzyFilters)
		if err != nil {
			return nil, err
		}
		docs = fuzzyFilterDocuments(resolved.Query, allDocs)
	}

	return docs, nil
}

func (s *SearchService) ListRecent(ctx context.Context, limit int, noCache bool) ([]domain.Document, error) {
	if noCache {
		return s.apiDoc.ListDocuments(ctx, domain.SearchFilters{})
	}

	if err := s.syncer.SyncIfNeeded(ctx); err != nil {
		return nil, err
	}

	return s.cacheDoc.ListRecentDocuments(ctx, limit)
}

func (s *SearchService) resolveFilterNames(ctx context.Context, filters domain.SearchFilters) (domain.SearchFilters, error) {
	resolved := filters

	if filters.CorrespondentName != "" && filters.CorrespondentID == nil {
		correspondents, err := s.cacheTax.ListCorrespondents(ctx)
		if err != nil {
			return resolved, err
		}
		if id := findByName(filters.CorrespondentName, correspondentsToNamedItems(correspondents)); id != nil {
			resolved.CorrespondentID = id
		}
	}

	if filters.DocumentTypeName != "" && filters.DocumentTypeID == nil {
		types, err := s.cacheTax.ListDocumentTypes(ctx)
		if err != nil {
			return resolved, err
		}
		if id := findByName(filters.DocumentTypeName, typesToNamedItems(types)); id != nil {
			resolved.DocumentTypeID = id
		}
	}

	if filters.TagName != "" && filters.TagID == nil {
		tags, err := s.cacheTax.ListTags(ctx)
		if err != nil {
			return resolved, err
		}
		if id := findByName(filters.TagName, tagsToNamedItems(tags)); id != nil {
			resolved.TagID = id
		}
	}

	return resolved, nil
}

type namedItem struct {
	ID   int
	Name string
}

func findByName(query string, items []namedItem) *int {
	lower := strings.ToLower(query)

	// 1. Exact match (case-insensitive)
	for _, item := range items {
		if strings.ToLower(item.Name) == lower {
			id := item.ID
			return &id
		}
	}

	// 2. Substring match
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Name), lower) {
			id := item.ID
			return &id
		}
	}

	// 3. Fuzzy match
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = item.Name
	}
	matches := fuzzy.Find(query, names)
	if len(matches) > 0 {
		id := items[matches[0].Index].ID
		return &id
	}

	return nil
}

func correspondentsToNamedItems(items []domain.Correspondent) []namedItem {
	result := make([]namedItem, len(items))
	for i, c := range items {
		result[i] = namedItem{ID: c.ID, Name: c.Name}
	}
	return result
}

func typesToNamedItems(items []domain.DocumentType) []namedItem {
	result := make([]namedItem, len(items))
	for i, t := range items {
		result[i] = namedItem{ID: t.ID, Name: t.Name}
	}
	return result
}

func tagsToNamedItems(items []domain.Tag) []namedItem {
	result := make([]namedItem, len(items))
	for i, t := range items {
		result[i] = namedItem{ID: t.ID, Name: t.Name}
	}
	return result
}

func fuzzyFilterDocuments(query string, docs []domain.Document) []domain.Document {
	titles := make([]string, len(docs))
	for i, d := range docs {
		titles[i] = d.Title
	}
	matches := fuzzy.Find(query, titles)
	result := make([]domain.Document, len(matches))
	for i, m := range matches {
		result[i] = docs[m.Index]
	}
	return result
}
