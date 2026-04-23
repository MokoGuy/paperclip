package sqlite

import (
	"context"
	"database/sql"

	"github.com/MokoGuy/paperclip/internal/domain"
	sqlitedb "github.com/MokoGuy/paperclip/internal/repository/sqlite/db"
)

type TaxonomyRepository struct {
	queries *sqlitedb.Queries
}

func NewTaxonomyRepository(conn *sql.DB) *TaxonomyRepository {
	return &TaxonomyRepository{queries: sqlitedb.New(conn)}
}

func (r *TaxonomyRepository) ListTags(ctx context.Context) ([]domain.Tag, error) {
	rows, err := r.queries.ListTags(ctx)
	if err != nil {
		return nil, err
	}

	tags := make([]domain.Tag, len(rows))
	for i, row := range rows {
		tags[i] = domain.Tag{
			ID:            int(row.ID),
			Name:          row.Name,
			Slug:          strOrEmpty(row.Slug),
			Color:         strOrEmpty(row.Color),
			DocumentCount: int(intOrZero(row.DocumentCount)),
		}
	}
	return tags, nil
}

func (r *TaxonomyRepository) ListDocumentTypes(ctx context.Context) ([]domain.DocumentType, error) {
	rows, err := r.queries.ListDocumentTypes(ctx)
	if err != nil {
		return nil, err
	}

	types := make([]domain.DocumentType, len(rows))
	for i, row := range rows {
		types[i] = domain.DocumentType{
			ID:            int(row.ID),
			Name:          row.Name,
			Slug:          strOrEmpty(row.Slug),
			DocumentCount: int(intOrZero(row.DocumentCount)),
		}
	}
	return types, nil
}

func (r *TaxonomyRepository) ListCorrespondents(ctx context.Context) ([]domain.Correspondent, error) {
	rows, err := r.queries.ListCorrespondents(ctx)
	if err != nil {
		return nil, err
	}

	correspondents := make([]domain.Correspondent, len(rows))
	for i, row := range rows {
		correspondents[i] = domain.Correspondent{
			ID:            int(row.ID),
			Name:          row.Name,
			Slug:          strOrEmpty(row.Slug),
			DocumentCount: int(intOrZero(row.DocumentCount)),
		}
	}
	return correspondents, nil
}

func (r *TaxonomyRepository) UpsertTag(ctx context.Context, t domain.Tag) error {
	return r.queries.UpsertTag(ctx, sqlitedb.UpsertTagParams{
		ID:            int64(t.ID),
		Name:          t.Name,
		Slug:          toNullString(t.Slug),
		Color:         toNullString(t.Color),
		DocumentCount: toNullInt64(int64(t.DocumentCount)),
	})
}

func (r *TaxonomyRepository) UpsertDocumentType(ctx context.Context, dt domain.DocumentType) error {
	return r.queries.UpsertDocumentType(ctx, sqlitedb.UpsertDocumentTypeParams{
		ID:            int64(dt.ID),
		Name:          dt.Name,
		Slug:          toNullString(dt.Slug),
		DocumentCount: toNullInt64(int64(dt.DocumentCount)),
	})
}

func (r *TaxonomyRepository) UpsertCorrespondent(ctx context.Context, c domain.Correspondent) error {
	return r.queries.UpsertCorrespondent(ctx, sqlitedb.UpsertCorrespondentParams{
		ID:            int64(c.ID),
		Name:          c.Name,
		Slug:          toNullString(c.Slug),
		DocumentCount: toNullInt64(int64(c.DocumentCount)),
	})
}

func (r *TaxonomyRepository) HasData(ctx context.Context) (bool, error) {
	count, err := r.queries.GetTagCount(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func strOrEmpty(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func intOrZero(n sql.NullInt64) int64 {
	if n.Valid {
		return n.Int64
	}
	return 0
}

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func toNullInt64(n int64) sql.NullInt64 {
	return sql.NullInt64{Int64: n, Valid: true}
}
