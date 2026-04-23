package sqlite

import (
	"context"
	"database/sql"
	"time"

	sqlitedb "github.com/MokoGuy/paperclip/internal/repository/sqlite/db"
)

type SyncRepository struct {
	queries *sqlitedb.Queries
}

func NewSyncRepository(conn *sql.DB) *SyncRepository {
	return &SyncRepository{queries: sqlitedb.New(conn)}
}

func (r *SyncRepository) GetLastSync(ctx context.Context) (time.Time, error) {
	val, err := r.queries.GetSyncState(ctx, "last_sync")
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}
	if !val.Valid {
		return time.Time{}, nil
	}
	t, err := time.Parse(time.RFC3339, val.String)
	if err != nil {
		return time.Time{}, nil
	}
	return t, nil
}

func (r *SyncRepository) SetLastSync(ctx context.Context, t time.Time) error {
	return r.queries.SetSyncState(ctx, sqlitedb.SetSyncStateParams{
		Key:   "last_sync",
		Value: toNullString(t.Format(time.RFC3339)),
	})
}
