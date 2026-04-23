package domain

type Tag struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Slug          string `json:"slug,omitempty"`
	Color         string `json:"color,omitempty"`
	DocumentCount int    `json:"document_count"`
}

type DocumentType struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Slug          string `json:"slug,omitempty"`
	DocumentCount int    `json:"document_count"`
}

type Correspondent struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Slug          string `json:"slug,omitempty"`
	DocumentCount int    `json:"document_count"`
}
