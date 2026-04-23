package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"

	"github.com/MokoGuy/paperclip/internal/domain"
)

type Formatter struct {
	forceJSON bool
	baseURL   string
}

func NewFormatter(forceJSON bool, baseURL string) *Formatter {
	return &Formatter{forceJSON: forceJSON, baseURL: baseURL}
}

func (f *Formatter) IsJSON() bool {
	if f.forceJSON {
		return true
	}
	return !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd())
}

// --- JSON response types (stable contract for LLM agents) ---

type TagResponse struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	DocumentCount int    `json:"document_count"`
	Color         string `json:"color,omitempty"`
}

type DocumentTypeResponse struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	DocumentCount int    `json:"document_count"`
}

type CorrespondentResponse struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	DocumentCount int    `json:"document_count"`
}

type RefResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DocumentResponse struct {
	ID               int           `json:"id"`
	Title            string        `json:"title"`
	Correspondent    *RefResponse  `json:"correspondent"`
	DocumentType     *RefResponse  `json:"document_type"`
	Tags             []RefResponse `json:"tags"`
	Created          string        `json:"created"`
	Added            string        `json:"added"`
	PageCount        int           `json:"page_count"`
	URL              string        `json:"url"`
}

type ContentItemResponse struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ListResponse[T any] struct {
	Results []T `json:"results"`
	Count   int `json:"count"`
}

// --- Render methods ---

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	nameStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	countStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	idStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
)

func (f *Formatter) RenderTags(tags []domain.Tag) error {
	sort.Slice(tags, func(i, j int) bool { return tags[i].DocumentCount > tags[j].DocumentCount })

	if f.IsJSON() {
		items := make([]TagResponse, len(tags))
		for i, t := range tags {
			items[i] = TagResponse{ID: t.ID, Name: t.Name, DocumentCount: t.DocumentCount, Color: t.Color}
		}
		return f.renderJSON(ListResponse[TagResponse]{Results: items, Count: len(items)})
	}

	fmt.Println(headerStyle.Render("Tags"))
	fmt.Println()
	for _, t := range tags {
		colorDot := "●"
		if t.Color != "" {
			colorDot = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Color)).Render("●")
		}
		fmt.Printf("  %s %s %s\n",
			colorDot,
			nameStyle.Render(t.Name),
			countStyle.Render(fmt.Sprintf("(%d)", t.DocumentCount)),
		)
	}
	return nil
}

func (f *Formatter) RenderDocumentTypes(types []domain.DocumentType) error {
	sort.Slice(types, func(i, j int) bool { return types[i].DocumentCount > types[j].DocumentCount })

	if f.IsJSON() {
		items := make([]DocumentTypeResponse, len(types))
		for i, t := range types {
			items[i] = DocumentTypeResponse{ID: t.ID, Name: t.Name, DocumentCount: t.DocumentCount}
		}
		return f.renderJSON(ListResponse[DocumentTypeResponse]{Results: items, Count: len(items)})
	}

	fmt.Println(headerStyle.Render("Document Types"))
	fmt.Println()
	for _, t := range types {
		fmt.Printf("  %s %s\n",
			nameStyle.Render(t.Name),
			countStyle.Render(fmt.Sprintf("(%d)", t.DocumentCount)),
		)
	}
	return nil
}

func (f *Formatter) RenderCorrespondents(correspondents []domain.Correspondent) error {
	sort.Slice(correspondents, func(i, j int) bool { return correspondents[i].DocumentCount > correspondents[j].DocumentCount })

	if f.IsJSON() {
		items := make([]CorrespondentResponse, len(correspondents))
		for i, c := range correspondents {
			items[i] = CorrespondentResponse{ID: c.ID, Name: c.Name, DocumentCount: c.DocumentCount}
		}
		return f.renderJSON(ListResponse[CorrespondentResponse]{Results: items, Count: len(items)})
	}

	fmt.Println(headerStyle.Render("Correspondents"))
	fmt.Println()
	for _, c := range correspondents {
		fmt.Printf("  %s %s\n",
			nameStyle.Render(c.Name),
			countStyle.Render(fmt.Sprintf("(%d)", c.DocumentCount)),
		)
	}
	return nil
}

func (f *Formatter) RenderDocuments(docs []domain.Document) error {
	if f.IsJSON() {
		items := make([]DocumentResponse, len(docs))
		for i, d := range docs {
			item := DocumentResponse{
				ID:        d.ID,
				Title:     d.Title,
				Created:   d.Created,
				Added:     d.Added,
				PageCount: d.PageCount,
				URL:       fmt.Sprintf("%s/documents/%d/", strings.TrimRight(f.baseURL, "/"), d.ID),
				Tags:      make([]RefResponse, 0),
			}
			if d.CorrespondentID != nil {
				item.Correspondent = &RefResponse{ID: *d.CorrespondentID, Name: d.CorrespondentName}
			}
			if d.DocumentTypeID != nil {
				item.DocumentType = &RefResponse{ID: *d.DocumentTypeID, Name: d.DocumentTypeName}
			}
			for j, tagID := range d.TagIDs {
				name := ""
				if j < len(d.TagNames) {
					name = d.TagNames[j]
				}
				item.Tags = append(item.Tags, RefResponse{ID: tagID, Name: name})
			}
			items[i] = item
		}
		return f.renderJSON(ListResponse[DocumentResponse]{Results: items, Count: len(items)})
	}

	fmt.Println(headerStyle.Render(fmt.Sprintf("Documents (%d)", len(docs))))
	fmt.Println()
	for _, d := range docs {
		correspondent := ""
		if d.CorrespondentName != "" {
			correspondent = d.CorrespondentName + " — "
		}
		fmt.Printf("  %s %s%s %s\n",
			idStyle.Render(fmt.Sprintf("#%d", d.ID)),
			correspondent,
			nameStyle.Render(d.Title),
			countStyle.Render(fmt.Sprintf("[%s]", d.Created)),
		)
	}
	return nil
}

func (f *Formatter) RenderIDs(docs []domain.Document) {
	for _, d := range docs {
		fmt.Println(d.ID)
	}
}

func (f *Formatter) RenderContent(items []ContentItemResponse) error {
	if f.IsJSON() {
		return f.renderJSON(ListResponse[ContentItemResponse]{Results: items, Count: len(items)})
	}

	for i, item := range items {
		if i > 0 {
			fmt.Println("---")
		}
		fmt.Printf("%s %s\n\n", idStyle.Render(fmt.Sprintf("#%d", item.ID)), headerStyle.Render(item.Title))
		fmt.Println(item.Content)
	}
	return nil
}

func (f *Formatter) renderJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
