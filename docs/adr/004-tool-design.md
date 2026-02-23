# ADR-004: MCP Tool Design — Resources vs Tools

## Status

Accepted

## Context

MCP supports two primary primitives for exposing data:

- **Resources** — static or URI-addressable content (like files); the client can read them
- **Tools** — callable functions with typed inputs that return structured results

We need to decide how to expose the SRD content. The data includes ~300 spells, ~300 monsters, 12 classes, ~200 equipment items, ~200 magic items, ~80 feats, 13 backgrounds, 10+ species, ~60 glossary terms, and hierarchical rules.

### Options Considered

1. **Resources only** — expose each content type as a resource URI (e.g., `srd://spells/fireball`)
2. **Tools only** — expose search and lookup tools that return content
3. **Hybrid** — resources for static content, tools for search/filtering

## Decision

Use **tools only**. Expose the SRD content through a focused set of callable tools.

## Rationale

- **Claude Code integration** — Claude Code primarily interacts with MCP servers through tool calls, not resource browsing. Tools give Claude the ability to search and filter, which is how an AI agent naturally consumes reference data.
- **Context window efficiency** — returning full collection dumps would flood Claude's context. Tools with search/filter parameters let Claude request exactly what it needs (e.g., "3rd level evocation spells" instead of all 300 spells).
- **Discoverability** — tool descriptions and input schemas tell Claude *what it can ask for*, acting as a self-documenting API.
- **Consistency** — `due-draghi-combattimenti` uses tools exclusively for its MCP interface; same pattern here.

### Tool Design

The server exposes the following tools:

| Tool | Purpose | Key Inputs |
|------|---------|------------|
| `search_spells` | Search and filter spells | `query?`, `level?`, `school?`, `class?`, `ritual?` |
| `get_spell` | Get a single spell by ID or name | `id` |
| `search_monsters` | Search and filter monsters | `query?`, `type?`, `size?`, `cr?` |
| `get_monster` | Get a single monster by ID or name | `id` |
| `search_equipment` | Search equipment | `query?`, `category?` |
| `get_equipment` | Get a single equipment item | `id` |
| `search_magic_items` | Search magic items | `query?`, `type?`, `rarity?` |
| `get_magic_item` | Get a single magic item | `id` |
| `search_feats` | Search feats | `query?`, `category?` |
| `get_feat` | Get a single feat | `id` |
| `get_class` | Get class details | `id` |
| `list_classes` | List all classes | (none) |
| `get_background` | Get background details | `id` |
| `list_backgrounds` | List all backgrounds | (none) |
| `get_species` | Get species details | `id` |
| `list_species` | List all species | (none) |
| `search_rules` | Search game rules | `query` |
| `lookup_glossary` | Look up a game term | `term` |

### Design Principles

1. **Search tools return summaries** — list results include ID, name, and key metadata (level, CR, rarity, etc.) but not full descriptions, to keep responses compact
2. **Get tools return full detail** — single-item lookups return the complete entry including description text
3. **Fuzzy search** — `query` parameters use case-insensitive substring matching so Claude doesn't need exact names
4. **Filter composition** — filter parameters combine with AND logic (e.g., `level=3 AND school=Evocation`)
5. **Small collections get list tools** — classes (12), backgrounds (13), and species (10+) are small enough to list entirely

## Consequences

- More tools to implement and maintain than a resource-based approach
- Each content type needs search/get handler pairs (mitigated by shared patterns)
- Tool descriptions must be clear enough for Claude to know when to use each tool
- Search results need a reasonable default limit to avoid oversized responses
