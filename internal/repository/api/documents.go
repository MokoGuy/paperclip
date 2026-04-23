package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/MokoGuy/paperclip/internal/domain"
)

type apiRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type apiDocument struct {
	ID                  int               `json:"id"`
	Title               string            `json:"title"`
	Correspondent       json.RawMessage   `json:"correspondent"`
	DocumentType        json.RawMessage   `json:"document_type"`
	Tags                json.RawMessage   `json:"tags"`
	Created             string            `json:"created"`
	CreatedDate         string            `json:"created_date"`
	Added               string            `json:"added"`
	Modified            string            `json:"modified"`
	ArchiveSerialNumber *int              `json:"archive_serial_number"`
	OriginalFileName    string            `json:"original_file_name"`
	PageCount           int               `json:"page_count"`
}

func parseRef(raw json.RawMessage) *apiRef {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var ref apiRef
	if err := json.Unmarshal(raw, &ref); err == nil {
		return &ref
	}
	var id int
	if err := json.Unmarshal(raw, &id); err == nil {
		return &apiRef{ID: id}
	}
	return nil
}

func parseTags(raw json.RawMessage) []apiRef {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var refs []apiRef
	if err := json.Unmarshal(raw, &refs); err == nil {
		return refs
	}
	var ids []int
	if err := json.Unmarshal(raw, &ids); err == nil {
		refs := make([]apiRef, len(ids))
		for i, id := range ids {
			refs[i] = apiRef{ID: id}
		}
		return refs
	}
	return nil
}

type DocumentRepository struct {
	client *Client
}

func NewDocumentRepository(client *Client) *DocumentRepository {
	return &DocumentRepository{client: client}
}

func (r *DocumentRepository) FetchAllDocuments(_ context.Context) ([]domain.Document, error) {
	items, err := fetchAllPages[apiDocument](r.client, "/api/documents/", nil)
	if err != nil {
		return nil, err
	}
	return mapDocuments(items), nil
}

func (r *DocumentRepository) SearchDocuments(_ context.Context, query string) ([]domain.Document, error) {
	params := url.Values{}
	params.Set("search", query)
	items, err := fetchAllPages[apiDocument](r.client, "/api/documents/", params)
	if err != nil {
		return nil, err
	}
	return mapDocuments(items), nil
}

func (r *DocumentRepository) ListDocuments(_ context.Context, filters domain.SearchFilters) ([]domain.Document, error) {
	params := url.Values{}
	if filters.CorrespondentID != nil {
		params.Set("correspondent", strconv.Itoa(*filters.CorrespondentID))
	}
	if filters.DocumentTypeID != nil {
		params.Set("document_type", strconv.Itoa(*filters.DocumentTypeID))
	}
	if filters.TagID != nil {
		params.Set("tag", strconv.Itoa(*filters.TagID))
	}
	if filters.After != "" {
		params.Set("created__date__gte", filters.After)
	}
	if filters.Before != "" {
		params.Set("created__date__lte", filters.Before)
	}
	params.Set("ordering", "-created")

	items, err := fetchAllPages[apiDocument](r.client, "/api/documents/", params)
	if err != nil {
		return nil, err
	}
	return mapDocuments(items), nil
}

type DocumentDetail struct {
	Title   string
	Content string
}

func (r *DocumentRepository) GetContent(_ context.Context, id int) (*DocumentDetail, error) {
	body, err := r.client.get(fmt.Sprintf("/api/documents/%d/", id), nil)
	if err != nil {
		return nil, err
	}

	var doc struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	return &DocumentDetail{Title: doc.Title, Content: doc.Content}, nil
}

func mapDocuments(items []apiDocument) []domain.Document {
	docs := make([]domain.Document, len(items))
	for i, d := range items {
		doc := domain.Document{
			ID:               d.ID,
			Title:            d.Title,
			Created:          d.CreatedDate,
			Added:            d.Added,
			Modified:         d.Modified,
			OriginalFileName: d.OriginalFileName,
			PageCount:        d.PageCount,
		}
		if d.ArchiveSerialNumber != nil {
			asn := *d.ArchiveSerialNumber
			doc.ArchiveSerialNumber = &asn
		}
		if ref := parseRef(d.Correspondent); ref != nil {
			doc.CorrespondentID = &ref.ID
			doc.CorrespondentName = ref.Name
		}
		if ref := parseRef(d.DocumentType); ref != nil {
			doc.DocumentTypeID = &ref.ID
			doc.DocumentTypeName = ref.Name
		}
		for _, t := range parseTags(d.Tags) {
			doc.TagIDs = append(doc.TagIDs, t.ID)
			doc.TagNames = append(doc.TagNames, t.Name)
		}
		docs[i] = doc
	}
	return docs
}
