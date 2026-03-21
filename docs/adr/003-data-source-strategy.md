# ADR-003: Shared Go Module as Data Source

## Status

Accepted (supersedes previous JSON copy strategy)

## Context

The MCP server needs access to D&D 5e Italian SRD data. The `quintaedizione-data-ita` Go module embeds all JSON data and provides typed structs, search, and filtering.

### Options Considered

1. **Import shared Go module** — call `store.Load()`, get typed access
2. **Copy JSON files** — duplicate data into MCP project
3. **Git submodule** — reference shared module at build time

## Decision

Import `github.com/emiliopalmerini/quintaedizione-data-ita` as a Go module dependency.

## Rationale

- **No data duplication** — single source of truth in the shared module
- **Typed access** — `Store.Spell(id)` instead of parsing raw JSON
- **Search and filtering included** — reuse fuzzy search and predicate builder
- **Data updates** propagate via `go get -u`

## Consequences

- Go module dependency on quintaedizione-data-ita
- Data updates require bumping the module version
