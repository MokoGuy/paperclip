-- name: GetSyncState :one
SELECT value FROM sync_state WHERE key = ?;

-- name: SetSyncState :exec
INSERT INTO sync_state (key, value)
VALUES (?, ?)
ON CONFLICT (key) DO UPDATE SET value = excluded.value;
