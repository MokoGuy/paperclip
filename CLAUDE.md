# CLAUDE.md

## What is this project?

paperCLIp — CLI Go d'exploration Paperless-NGX, dual humain/agent LLM.
Cache local SQLite + fuzzy-match, recherche composable, extraction batch.

## Development Commands

- `make build` — Build binary to `bin/paperclip`
- `make run` — Build and run
- `make test` — Run all tests
- `make generate` — Generate SQLC code from SQL files
- `make deps` — Install/update Go dependencies
- `make fmt` — Format Go code
- `make dev-setup` — Full setup (deps + generate)

## Architecture

Clean Architecture with 4 layers:

1. **Domain** (`internal/domain/`) — Entities, interfaces, config. No business logic.
2. **Use Case** (`internal/usecase/`) — Business logic orchestration (search, sync, content).
3. **Repository** (`internal/repository/`) — Data access:
   - `api/` — Paperless-NGX REST API client
   - `sqlite/` — Local cache (SQLC-generated, embedded migrations)
4. **Delivery** (`internal/delivery/cli/`) — Cobra commands + output formatting.

### Key Patterns

- **Segregated interfaces** (ISP): `TaxonomyReader`, `DocumentReader`, `ContentFetcher`
- **Dual output**: auto-detect TTY → table (lipgloss) / non-TTY → JSON (stable schema for LLM agents)
- **Lazy sync**: cache SQLite auto-refreshed if >24h, `--no-cache` to bypass
- **Fuzzy resolution**: filter names (`--from`, `--type`, `--tag`) are fuzzy-matched against cached taxonomy

### Config

`~/.config/paperclip/config.toml` with `url` and `token` (hardcoded, chmod 600 enforced).

### Database

- SQLite with embedded migrations in `internal/repository/sqlite/migrations/`
- SQLC generates type-safe Go code from queries in `queries/`
- Schema: documents, tags, document_types, correspondents, document_tags (M2M), sync_state

## When modifying

- Business logic → use case layer, not repositories or CLI commands
- New SQL → add migration + queries, run `make generate`
- CLI commands → `internal/delivery/cli/`, delegate to use case services
- JSON output schema is a stable contract — don't change field names or structure without versioning
