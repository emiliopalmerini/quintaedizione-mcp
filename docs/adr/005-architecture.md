# ADR-005: Project Architecture

## Status

Accepted

## Context

We need to decide how to structure the Bun TypeScript MCP server. The existing Go projects follow hexagonal architecture with clear domain/application/adapter separation. We need to balance architectural rigor with the simpler nature of this project — it's a read-only data server with no mutations, no external APIs, and no persistence layer beyond static JSON files.

### Options Considered

1. **Full hexagonal architecture** — ports, adapters, domain, application layers (mirrors Go projects)
2. **Flat structure** — single directory with all files
3. **Layered by responsibility** — data loading, content types, tools, search

## Decision

Use a **lightweight layered architecture** organized by responsibility, not full hexagonal.

## Rationale

This server is fundamentally a read-only data query layer. There are no mutations, no external service calls, no complex domain logic. Full hexagonal architecture would add ceremony without value. Instead, we organize by what the code does:

```
quintaedizione-mcp/
├── src/
│   ├── index.ts              # Entry point: create server, register tools, serve stdio
│   ├── server.ts             # MCP server setup and tool registration
│   ├── data/
│   │   └── loader.ts         # Load and index JSON files at startup
│   ├── content/
│   │   ├── types.ts          # TypeScript interfaces for all content types
│   │   ├── spells.ts         # Spell search/get logic
│   │   ├── monsters.ts       # Monster search/get logic
│   │   ├── equipment.ts      # Equipment search/get logic
│   │   ├── magic-items.ts    # Magic item search/get logic
│   │   ├── feats.ts          # Feat search/get logic
│   │   ├── classes.ts        # Class list/get logic
│   │   ├── backgrounds.ts    # Background list/get logic
│   │   ├── species.ts        # Species list/get logic
│   │   ├── rules.ts          # Rules search logic
│   │   └── glossary.ts       # Glossary lookup logic
│   ├── search/
│   │   └── fuzzy.ts          # Shared fuzzy search utilities
│   └── tools/
│       ├── spells.ts         # Spell tool definitions (schema + handler)
│       ├── monsters.ts       # Monster tool definitions
│       ├── equipment.ts      # Equipment tool definitions
│       ├── magic-items.ts    # Magic item tool definitions
│       ├── feats.ts          # Feat tool definitions
│       ├── classes.ts        # Class tool definitions
│       ├── backgrounds.ts    # Background tool definitions
│       ├── species.ts        # Species tool definitions
│       ├── rules.ts          # Rules tool definitions
│       └── glossary.ts       # Glossary tool definitions
├── data/                     # Copied SRD JSON files (see ADR-003)
│   ├── spells.json
│   ├── monsters.json
│   ├── classes.json
│   ├── equipment.json
│   ├── magic_items.json
│   ├── feats.json
│   ├── backgrounds.json
│   ├── species.json
│   ├── glossary.json
│   ├── rules_gameplay.json
│   ├── rules_creation.json
│   └── rules_tools.json
├── scripts/
│   └── sync-data.sh          # Copy JSON from quintaedizione.online
├── tests/
│   ├── content/              # Unit tests for content modules
│   └── tools/                # Integration tests for tool handlers
├── docs/
│   └── adr/                  # Architecture Decision Records
├── package.json
├── tsconfig.json
└── CLAUDE.md                 # Project instructions for Claude Code
```

### Layer Responsibilities

| Layer | Purpose |
|-------|---------|
| `data/loader.ts` | Reads JSON files at startup, builds in-memory indexes (by ID, by name). Single load, immutable after init. |
| `content/*.ts` | Pure functions that query the in-memory data. Each module owns search/filter/get logic for its content type. No framework dependencies. |
| `search/fuzzy.ts` | Shared search utilities: case-insensitive matching, substring search, result limiting. |
| `tools/*.ts` | MCP tool definitions: Zod schemas for inputs, handler functions that call content modules and format MCP responses. |
| `server.ts` | Wires tools to the MCP server instance. |
| `index.ts` | Composition root: loads data, creates server, starts stdio transport. |

### Key Principles

1. **Content modules are framework-free** — they take typed inputs and return typed outputs; no MCP SDK dependency. This makes them independently testable.
2. **Tool modules are thin adapters** — they define Zod schemas, call content functions, and format `CallToolResult` responses.
3. **Data is loaded once** — all JSON is read at startup into memory. No lazy loading, no file I/O during tool calls.
4. **No over-abstraction** — no repository interfaces, no dependency injection containers. Direct imports are fine for a read-only server.

## Consequences

- Simple to understand and navigate; each file has a clear purpose
- Content logic is testable without MCP SDK mocking
- Adding a new content type requires: type definition, content module, tool module, loader entry
- Not as rigorous as the Go projects' hexagonal architecture (acceptable — this is a simpler system)
