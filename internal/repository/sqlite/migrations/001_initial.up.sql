CREATE TABLE correspondents (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT,
    document_count INTEGER DEFAULT 0
);

CREATE TABLE document_types (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT,
    document_count INTEGER DEFAULT 0
);

CREATE TABLE tags (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT,
    color TEXT,
    document_count INTEGER DEFAULT 0
);

CREATE TABLE documents (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    correspondent_id INTEGER REFERENCES correspondents(id),
    document_type_id INTEGER REFERENCES document_types(id),
    created TEXT,
    added TEXT,
    modified TEXT,
    archive_serial_number INTEGER,
    original_file_name TEXT,
    page_count INTEGER
);

CREATE TABLE document_tags (
    document_id INTEGER REFERENCES documents(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (document_id, tag_id)
);

CREATE TABLE sync_state (
    key TEXT PRIMARY KEY,
    value TEXT
);
