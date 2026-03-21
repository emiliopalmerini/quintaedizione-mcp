# quintaedizione-mcp

MCP server exposing D&D 5e Italian SRD data to AI assistants.

## Build & Run

```bash
make build    # Build binary
make test     # Run tests
make install  # Install to $GOPATH/bin
```

## Prerequisites

- Go 1.25.2

## Architecture

Thin MCP adapter over the `quintaedizione-data-ita` shared Go module. All data, search, and filtering logic lives in the shared module. This server only defines MCP tool schemas and formats responses.

```
main.go         Entry point: load store, register tools, serve stdio
tools/          MCP tool handlers (one file per content type)
docs/adr/       Architecture Decision Records
```

## Tools

| Tool | Description |
|------|-------------|
| `search` | Cross-collection fuzzy search |
| `search_spells` / `get_spell` | Spell search and detail |
| `search_monsters` / `get_monster` | Monster search and detail |
| `list_classes` / `get_class` | Class listing and detail |
| `search_equipment` / `get_equipment` | Equipment search and detail |
| `search_magic_items` / `get_magic_item` | Magic item search and detail |
| `search_feats` / `get_feat` | Feat search and detail |
| `list_backgrounds` / `get_background` | Background listing and detail |
| `list_species` / `get_species` | Species listing and detail |
| `search_rules` / `get_rule` | Rules search and detail |
| `lookup_glossary` | Glossary term lookup |
| `search_maps` | Battle map search |
| `list_generators` / `get_generator` | Random generator tables |

## Claude Code Configuration

```json
{
  "mcpServers": {
    "quintaedizione": {
      "command": "quintaedizione-mcp"
    }
  }
}
```
