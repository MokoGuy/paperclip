package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/MokoGuy/paperclip/internal/domain"
	sqlitedb "github.com/MokoGuy/paperclip/internal/repository/sqlite/db"
)

type DocumentRepository struct {
	db      *sql.DB
	queries *sqlitedb.Queries
}

func NewDocumentRepository(conn *sql.DB) *DocumentRepository {
	return &DocumentRepository{db: conn, queries: sqlitedb.New(conn)}
}

func (r *DocumentRepository) UpsertDocument(ctx context.Context, d domain.Document) error {
	err := r.queries.UpsertDocument(ctx, sqlitedb.UpsertDocumentParams{
		ID:                  int64(d.ID),
		Title:               d.Title,
		CorrespondentID:     toNullInt64FromPtr(d.CorrespondentID),
		DocumentTypeID:      toNullInt64FromPtr(d.DocumentTypeID),
		Created:             toNullString(d.Created),
		Added:               toNullString(d.Added),
		Modified:            toNullString(d.Modified),
		ArchiveSerialNumber: toNullInt64FromPtr(d.ArchiveSerialNumber),
		OriginalFileName:    toNullString(d.OriginalFileName),
		PageCount:           toNullInt64(int64(d.PageCount)),
	})
	if err != nil {
		return err
	}

	if err := r.queries.DeleteDocumentTags(ctx, sql.NullInt64{Int64: int64(d.ID), Valid: true}); err != nil {
		return err
	}
	for _, tagID := range d.TagIDs {
		if err := r.queries.InsertDocumentTag(ctx, sqlitedb.InsertDocumentTagParams{
			DocumentID: sql.NullInt64{Int64: int64(d.ID), Valid: true},
			TagID:      sql.NullInt64{Int64: int64(tagID), Valid: true},
		}); err != nil {
			return err
		}
	}

	return nil
}

func (r *DocumentRepository) SearchDocuments(ctx context.Context, filters domain.SearchFilters) ([]domain.Document, error) {
	var conditions []string
	var args []any

	if filters.Query != "" {
		conditions = append(conditions, "d.title LIKE ?")
		args = append(args, "%"+filters.Query+"%")
	}
	if filters.CorrespondentID != nil {
		conditions = append(conditions, "d.correspondent_id = ?")
		args = append(args, *filters.CorrespondentID)
	}
	if filters.DocumentTypeID != nil {
		conditions = append(conditions, "d.document_type_id = ?")
		args = append(args, *filters.DocumentTypeID)
	}
	if filters.TagID != nil {
		conditions = append(conditions, "EXISTS (SELECT 1 FROM document_tags dt WHERE dt.document_id = d.id AND dt.tag_id = ?)")
		args = append(args, *filters.TagID)
	}
	if filters.Year != nil {
		conditions = append(conditions, "d.created LIKE ?")
		args = append(args, fmt.Sprintf("%d-%%", *filters.Year))
	}
	if filters.After != "" {
		conditions = append(conditions, "d.created >= ?")
		args = append(args, filters.After)
	}
	if filters.Before != "" {
		conditions = append(conditions, "d.created <= ?")
		args = append(args, filters.Before)
	}

	query := `SELECT d.id, d.title, d.correspondent_id, d.document_type_id, d.created, d.added, d.page_count
FROM documents d`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY d.created DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []domain.Document
	for rows.Next() {
		var d domain.Document
		var corrID, typeID sql.NullInt64
		var created, added sql.NullString
		var pageCount sql.NullInt64

		if err := rows.Scan(&d.ID, &d.Title, &corrID, &typeID, &created, &added, &pageCount); err != nil {
			return nil, err
		}

		if corrID.Valid {
			id := int(corrID.Int64)
			d.CorrespondentID = &id
		}
		if typeID.Valid {
			id := int(typeID.Int64)
			d.DocumentTypeID = &id
		}
		d.Created = strOrEmpty(created)
		d.Added = strOrEmpty(added)
		d.PageCount = int(intOrZero(pageCount))

		docs = append(docs, d)
	}

	if err := r.enrichDocuments(ctx, docs); err != nil {
		return nil, err
	}

	return docs, nil
}

func (r *DocumentRepository) ListRecentDocuments(ctx context.Context, limit int) ([]domain.Document, error) {
	rows, err := r.queries.ListRecentDocuments(ctx, int64(limit))
	if err != nil {
		return nil, err
	}

	docs := make([]domain.Document, len(rows))
	for i, row := range rows {
		docs[i] = domain.Document{
			ID:        int(row.ID),
			Title:     row.Title,
			Created:   strOrEmpty(row.Created),
			Added:     strOrEmpty(row.Added),
			PageCount: int(intOrZero(row.PageCount)),
		}
		if row.CorrespondentID.Valid {
			id := int(row.CorrespondentID.Int64)
			docs[i].CorrespondentID = &id
		}
		if row.DocumentTypeID.Valid {
			id := int(row.DocumentTypeID.Int64)
			docs[i].DocumentTypeID = &id
		}
	}

	if err := r.enrichDocuments(ctx, docs); err != nil {
		return nil, err
	}

	return docs, nil
}

func (r *DocumentRepository) enrichDocuments(ctx context.Context, docs []domain.Document) error {
	for i := range docs {
		if docs[i].CorrespondentID != nil {
			var name string
			err := r.db.QueryRowContext(ctx, "SELECT name FROM correspondents WHERE id = ?", *docs[i].CorrespondentID).Scan(&name)
			if err == nil {
				docs[i].CorrespondentName = name
			}
		}
		if docs[i].DocumentTypeID != nil {
			var name string
			err := r.db.QueryRowContext(ctx, "SELECT name FROM document_types WHERE id = ?", *docs[i].DocumentTypeID).Scan(&name)
			if err == nil {
				docs[i].DocumentTypeName = name
			}
		}

		var tagIDs []int
		var tagNames []string
		tagRows, err := r.db.QueryContext(ctx, "SELECT t.id, t.name FROM document_tags dt JOIN tags t ON t.id = dt.tag_id WHERE dt.document_id = ?", docs[i].ID)
		if err == nil {
			for tagRows.Next() {
				var id int
				var name string
				if err := tagRows.Scan(&id, &name); err == nil {
					tagIDs = append(tagIDs, id)
					tagNames = append(tagNames, name)
				}
			}
			tagRows.Close()
		}
		docs[i].TagIDs = tagIDs
		docs[i].TagNames = tagNames
	}
	return nil
}

func (r *DocumentRepository) GetDocumentCount(ctx context.Context) (int, error) {
	count, err := r.queries.GetDocumentCount(ctx)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func toNullInt64FromPtr(p *int) sql.NullInt64 {
	if p == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*p), Valid: true}
}
