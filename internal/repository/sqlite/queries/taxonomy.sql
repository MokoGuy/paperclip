-- name: UpsertTag :exec
INSERT INTO tags (id, name, slug, color, document_count)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT (id) DO UPDATE SET
    name = excluded.name,
    slug = excluded.slug,
    color = excluded.color,
    document_count = excluded.document_count;

-- name: UpsertDocumentType :exec
INSERT INTO document_types (id, name, slug, document_count)
VALUES (?, ?, ?, ?)
ON CONFLICT (id) DO UPDATE SET
    name = excluded.name,
    slug = excluded.slug,
    document_count = excluded.document_count;

-- name: UpsertCorrespondent :exec
INSERT INTO correspondents (id, name, slug, document_count)
VALUES (?, ?, ?, ?)
ON CONFLICT (id) DO UPDATE SET
    name = excluded.name,
    slug = excluded.slug,
    document_count = excluded.document_count;

-- name: ListTags :many
SELECT id, name, slug, color, document_count
FROM tags
ORDER BY document_count DESC;

-- name: ListDocumentTypes :many
SELECT id, name, slug, document_count
FROM document_types
ORDER BY document_count DESC;

-- name: ListCorrespondents :many
SELECT id, name, slug, document_count
FROM correspondents
ORDER BY document_count DESC;

-- name: GetTagCount :one
SELECT COUNT(*) FROM tags;

-- name: GetDocumentTypeCount :one
SELECT COUNT(*) FROM document_types;

-- name: GetCorrespondentCount :one
SELECT COUNT(*) FROM correspondents;
