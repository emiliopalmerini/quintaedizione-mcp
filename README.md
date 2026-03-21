# quintaedizione-mcp

MCP (Model Context Protocol) server for D&D 5e Italian SRD content. Exposes spells, monsters, classes, equipment, magic items, feats, backgrounds, species, rules, glossary, battle maps, and random generators to AI assistants like Claude.

All content is in **Italian** (Quinta Edizione), sourced from the [quintaedizione-data-ita](https://github.com/emiliopalmerini/quintaedizione-data-ita) shared Go module.

## Install

```bash
go install github.com/emiliopalmerini/quintaedizione-mcp@latest
```

## Usage with Claude Code

Add to your Claude Code MCP settings (`~/.claude/settings.json`):

```json
{
  "mcpServers": {
    "quintaedizione": {
      "command": "quintaedizione-mcp"
    }
  }
}
```

Then ask Claude about D&D content in Italian:

> "Cercami tutti gli incantesimi di livello 3 della scuola di evocazione"

## Available Tools

| Tool | Description |
|------|-------------|
| `search` | Cross-collection fuzzy search |
| `search_spells` / `get_spell` | Search and get spells |
| `search_monsters` / `get_monster` | Search and get monsters |
| `list_classes` / `get_class` | List and get classes |
| `search_equipment` / `get_equipment` | Search and get equipment |
| `search_magic_items` / `get_magic_item` | Search and get magic items |
| `search_feats` / `get_feat` | Search and get feats |
| `list_backgrounds` / `get_background` | List and get backgrounds |
| `list_species` / `get_species` | List and get species |
| `search_rules` / `get_rule` | Search and get rules |
| `lookup_glossary` | Look up game terms |
| `search_maps` | Search battle maps |
| `list_generators` / `get_generator` | Random generator tables |

## Build from source

```bash
git clone https://github.com/emiliopalmerini/quintaedizione-mcp.git
cd quintaedizione-mcp
make build
```
