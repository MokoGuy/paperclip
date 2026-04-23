# paperCLIp

[![Release](https://img.shields.io/github/v/release/MokoGuy/paperclip)](https://github.com/MokoGuy/paperclip/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/MokoGuy/paperclip)](https://goreportcard.com/report/github.com/MokoGuy/paperclip)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**A fast CLI for exploring your [Paperless-NGX](https://docs.paperless-ngx.com/) instance.** Built for humans in a terminal and LLM agents in a pipeline.

## Why?

Paperless-NGX's search is limited: you can full-text search **or** filter by metadata, but not both. Finding a specific invoice can take 5+ query reformulations. Extracting content from multiple documents means sequential API calls.

paperCLIp solves this with:

- **Composable filters** -- combine free-text, correspondent, type, tag, and date in one command
- **Local fuzzy search** -- SQLite cache with typo-tolerant matching (`amzn` finds `Amazon`, `bulltin` finds `Bulletin`)
- **Batch content extraction** -- fetch text from multiple documents in parallel
- **Dual output** -- colored tables for humans, structured JSON for LLM agents (auto-detected)

## Install

### Binary (recommended)

Download the latest release from the [releases page](https://github.com/MokoGuy/paperclip/releases):

```bash
# Linux amd64
curl -sL https://github.com/MokoGuy/paperclip/releases/latest/download/paperclip_1.0.0_linux_amd64.tar.gz | tar xz
sudo mv paperclip /usr/local/bin/
```

### From source

```bash
go install github.com/MokoGuy/paperclip/cmd/paperclip@latest
```

## Setup

```bash
paperclip init
# Paperless-NGX URL (e.g. https://paperless.example.com): https://your-instance.example.com
# API token: your-api-token-here
# Testing connection... OK
# Config saved to ~/.config/paperclip/config.toml
# Syncing local cache...
# Ready! Try 'paperclip search --recent 5'
```

One command: config, connection test, and cache sync. Use `--no-sync` to skip the initial sync.

> **Get your API token:** Paperless-NGX web UI → Settings → API tokens, or via `POST /api/token/`.

## Usage

### Search documents

Combine free-text queries with structured filters. All filters are composable.

```bash
# Free-text search on document titles
paperclip search "invoice"

# Filter by correspondent + document type + year
paperclip search --from amazon --type invoice --year 2024

# Fuzzy matching: typos and abbreviations just work
paperclip search --from amzn --type "invce"

# Last 10 documents added
paperclip search --recent 10

# Date range
paperclip search --after 2024-01-01 --before 2024-06-30

# Output only IDs (for piping)
paperclip search --from amazon --ids-only
# 502
# 503
# 510
```

### Extract document content

Fetch the full text content of documents (always live from the API, never cached).

```bash
# Single document
paperclip content 502

# Multiple documents (fetched in parallel)
paperclip content 502 503 505 510

# Pipeline: search then extract
paperclip search --from amazon --year 2024 --ids-only | xargs paperclip content

# Search within content
paperclip content 502 503 | grep "total amount"
```

### Explore taxonomy

See what's in your Paperless instance at a glance.

```bash
# Tags with document counts (sorted by count)
paperclip tags

# Document types
paperclip types

# Correspondents, filtered
paperclip correspondents --filter "bank"
```

### Force JSON output

Output format is auto-detected: **tables** when stdout is a terminal, **JSON** when piped. Use `--json` to force JSON in a terminal.

```bash
# Auto JSON when piped
paperclip search --recent 5 | jq '.'

# Force JSON in terminal
paperclip tags --json
```

<details>
<summary>JSON schema (stable contract for LLM agents)</summary>

**Search results:**
```json
{
  "results": [
    {
      "id": 502,
      "title": "Invoice January 2024",
      "correspondent": {"id": 1, "name": "Amazon"},
      "document_type": {"id": 4, "name": "Invoice"},
      "tags": [{"id": 10, "name": "Equipment"}],
      "created": "2024-01-15",
      "added": "2024-01-20T10:30:00+01:00",
      "page_count": 2,
      "url": "https://your-instance/documents/502/"
    }
  ],
  "count": 1
}
```

**Content extraction:**
```json
{
  "results": [
    {
      "id": 502,
      "title": "Invoice January 2024",
      "content": "INVOICE\nDate: January 15, 2024\n..."
    }
  ],
  "count": 1
}
```

**Taxonomy (tags/types/correspondents):**
```json
{
  "results": [
    {"id": 15, "name": "Finance", "document_count": 425, "color": "#75704E"}
  ],
  "count": 30
}
```

</details>

## How it works

paperCLIp maintains a **local SQLite cache** of your Paperless-NGX metadata (titles, dates, correspondents, types, tags). This enables instant fuzzy search without hitting the API.

| What | Cached locally | Fetched live |
|------|---------------|-------------|
| Document metadata (title, dates, IDs) | Yes | -- |
| Tags, types, correspondents | Yes | -- |
| Document text content | -- | Always |

**Cache sync strategy:**
- First run: automatic full sync
- Cache < 24 hours: instant local queries
- Cache > 24 hours: transparent re-sync on next command
- `paperclip sync`: force refresh anytime
- `--no-cache`: bypass cache entirely for a single query

## Architecture

```
cmd/paperclip/          Entry point
internal/
  domain/               Entities + interfaces (no business logic)
  usecase/              Search, sync, content services
  repository/
    api/                Paperless-NGX REST client
    sqlite/             Local cache (SQLC-generated queries)
  delivery/cli/         Cobra commands + output formatting
```

Clean Architecture. Read-only (no write operations). Segregated interfaces. SQLC for type-safe database access.

## License

MIT
