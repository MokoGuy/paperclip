package domain

import "context"

type TaxonomyReader interface {
	ListTags(ctx context.Context) ([]Tag, error)
	ListDocumentTypes(ctx context.Context) ([]DocumentType, error)
	ListCorrespondents(ctx context.Context) ([]Correspondent, error)
}

type DocumentReader interface {
	SearchDocuments(ctx context.Context, filters SearchFilters) ([]Document, error)
	ListRecentDocuments(ctx context.Context, limit int) ([]Document, error)
}

type ContentFetcher interface {
	GetContent(ctx context.Context, id int) (string, error)
}

type SearchFilters struct {
	Query             string
	CorrespondentID   *int
	CorrespondentName string
	DocumentTypeID    *int
	DocumentTypeName  string
	TagID             *int
	TagName           string
	Year              *int
	After             string
	Before            string
}
