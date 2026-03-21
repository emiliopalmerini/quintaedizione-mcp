# ADR-005: Project Architecture

## Status

Accepted (supersedes previous TypeScript layered architecture)

## Context

We need to structure the Go MCP server. It is a thin adapter that exposes the shared module's data via MCP tools. There is no domain logic beyond what the shared module provides.

## Decision

Use a **flat, minimal architecture**. The MCP server is a single `main.go` with tool handler files.

```
quintaedizione-mcp/
├── main.go                # Entry point: load store, register tools, serve stdio
├── tools/
│   ├── spells.go          # Spell search/get tool handlers
│   ├── monsters.go        # Monster search/get tool handlers
│   ├── classes.go         # Class list/get tool handlers
│   ├── equipment.go       # Equipment search/get tool handlers
│   ├── magic_items.go     # Magic item search/get tool handlers
│   ├── feats.go           # Feat search/get tool handlers
│   ├── backgrounds.go     # Background list/get tool handlers
│   ├── species.go         # Species list/get tool handlers
│   ├── rules.go           # Rules search tool handlers
│   ├── glossary.go        # Glossary lookup tool handlers
│   ├── maps.go            # Map gallery tool handlers
│   └── generators.go      # Random generator tool handlers
├── docs/
│   └── adr/               # Architecture Decision Records
├── go.mod
├── Makefile
└── CLAUDE.md
```

## Rationale

- **No domain logic** — the shared module owns data, search, and filtering
- **Tool handlers are thin** — parse MCP input, call store methods, format response
- **Single binary** — `go build` produces a self-contained executable
- **No over-abstraction** — no interfaces, no DI, no layers. Direct store access.

## Consequences

- Simple to understand — each tool file maps 1:1 to an MCP tool
- Adding a new tool: one file in `tools/`, register in `main.go`
- All complexity lives in the shared module, not here
