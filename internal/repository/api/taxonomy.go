package api

import (
	"context"

	"github.com/MokoGuy/paperclip/internal/domain"
)

type apiTag struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Color         string `json:"color"`
	DocumentCount int    `json:"document_count"`
}

type apiDocumentType struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	DocumentCount int    `json:"document_count"`
}

type apiCorrespondent struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	DocumentCount int    `json:"document_count"`
}

type TaxonomyRepository struct {
	client *Client
}

func NewTaxonomyRepository(client *Client) *TaxonomyRepository {
	return &TaxonomyRepository{client: client}
}

func (r *TaxonomyRepository) ListTags(_ context.Context) ([]domain.Tag, error) {
	items, err := fetchAllPages[apiTag](r.client, "/api/tags/", nil)
	if err != nil {
		return nil, err
	}

	tags := make([]domain.Tag, len(items))
	for i, t := range items {
		tags[i] = domain.Tag{
			ID:            t.ID,
			Name:          t.Name,
			Slug:          t.Slug,
			Color:         t.Color,
			DocumentCount: t.DocumentCount,
		}
	}
	return tags, nil
}

func (r *TaxonomyRepository) ListDocumentTypes(_ context.Context) ([]domain.DocumentType, error) {
	items, err := fetchAllPages[apiDocumentType](r.client, "/api/document_types/", nil)
	if err != nil {
		return nil, err
	}

	types := make([]domain.DocumentType, len(items))
	for i, t := range items {
		types[i] = domain.DocumentType{
			ID:            t.ID,
			Name:          t.Name,
			Slug:          t.Slug,
			DocumentCount: t.DocumentCount,
		}
	}
	return types, nil
}

func (r *TaxonomyRepository) ListCorrespondents(_ context.Context) ([]domain.Correspondent, error) {
	items, err := fetchAllPages[apiCorrespondent](r.client, "/api/correspondents/", nil)
	if err != nil {
		return nil, err
	}

	correspondents := make([]domain.Correspondent, len(items))
	for i, c := range items {
		correspondents[i] = domain.Correspondent{
			ID:            c.ID,
			Name:          c.Name,
			Slug:          c.Slug,
			DocumentCount: c.DocumentCount,
		}
	}
	return correspondents, nil
}
