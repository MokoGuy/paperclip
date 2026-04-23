-- name: UpsertDocument :exec
INSERT INTO documents (id, title, correspondent_id, document_type_id, created, added, modified, archive_serial_number, original_file_name, page_count)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (id) DO UPDATE SET
    title = excluded.title,
    correspondent_id = excluded.correspondent_id,
    document_type_id = excluded.document_type_id,
    created = excluded.created,
    added = excluded.added,
    modified = excluded.modified,
    archive_serial_number = excluded.archive_serial_number,
    original_file_name = excluded.original_file_name,
    page_count = excluded.page_count;

-- name: DeleteDocumentTags :exec
DELETE FROM document_tags WHERE document_id = ?;

-- name: InsertDocumentTag :exec
INSERT OR IGNORE INTO document_tags (document_id, tag_id) VALUES (?, ?);

-- name: ListRecentDocuments :many
SELECT d.id, d.title, d.correspondent_id, d.document_type_id, d.created, d.added, d.modified, d.archive_serial_number, d.original_file_name, d.page_count
FROM documents d
ORDER BY d.added DESC
LIMIT ?;

-- name: GetDocumentCount :one
SELECT COUNT(*) FROM documents;
