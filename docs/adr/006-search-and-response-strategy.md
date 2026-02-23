# ADR-006: Search and Response Strategy

## Status

Accepted

## Context

Claude Code interacts with MCP tools through its context window. Every tool response consumes tokens. The SRD contains substantial content — monster stat blocks can be hundreds of words, spell descriptions vary from one line to a full paragraph, and class features span multiple levels. We need a strategy that lets Claude find what it needs without overwhelming its context.

### Options Considered

1. **Always return full content** — every search result includes complete descriptions
2. **Pagination** — return N results per page with a cursor
3. **Summary + detail pattern** — search returns compact summaries; get returns full content
4. **Truncation** — return full content but truncate long fields

## Decision

Use the **summary + detail pattern** with a configurable result limit.

## Rationale

### Search Results: Compact Summaries

Search tools return a list of matches with key identifying metadata but **no description text**:

```json
// search_spells({ level: 3, school: "Evocazione" })
{
  "results": [
    { "id": "palla-di-fuoco", "name": "Palla di fuoco", "level": 3, "school": "Evocazione", "ritual": false },
    { "id": "fulmine", "name": "Fulmine", "level": 3, "school": "Evocazione", "ritual": false }
  ],
  "total": 2
}
```

This gives Claude enough to:
- Present options to the user
- Decide which item to look up in detail
- Answer simple questions ("What 3rd level evocation spells exist?") without a follow-up call

### Get Results: Full Content

Get tools return the complete entry including all fields and full description text:

```json
// get_spell({ id: "palla-di-fuoco" })
{
  "id": "palla-di-fuoco",
  "name": "Palla di fuoco",
  "level": 3,
  "school": "Evocazione",
  "classes": ["Mago", "Stregone"],
  "casting_time": "1 azione",
  "range": "45 metri",
  "components": "V, S, M (una piccola palla di sterco di pipistrello e zolfo)",
  "duration": "Istantanea",
  "description": "Un lampo luminoso...",
  "at_higher_levels": "Quando lanci questo incantesimo...",
  "ritual": false
}
```

### Search Implementation

- **Fuzzy matching** on `query` parameter — case-insensitive substring match against name (and keywords for monsters/items)
- **Exact filtering** on typed parameters — `level`, `school`, `type`, `rarity`, `cr` use exact match
- **AND composition** — all provided filters must match
- **Default limit: 20 results** — prevents oversized responses; adjustable via optional `limit` parameter (max 50)
- **Total count always included** — so Claude knows if results were truncated

### Special Cases

| Content Type | Search Strategy |
|-------------|-----------------|
| **Monsters** | Search by name, filter by type/size/CR. Summary includes: id, name, type, size, cr, hp, ac |
| **Rules** | Full-text search across rule titles and content. Returns matching section titles with their hierarchy path |
| **Glossary** | Exact term lookup with fuzzy fallback. Returns the definition directly (glossary entries are short) |
| **Classes/Backgrounds/Species** | List tools return all entries (small collections). Get tools return full detail |

## Consequences

- Claude typically needs two tool calls to get full content: search then get (acceptable — this is the standard MCP pattern)
- Summary fields must be chosen carefully per content type to be useful without being verbose
- Result limit prevents accidental context flooding (e.g., "search all spells" without filters)
- Fuzzy search means Claude doesn't need exact Italian spelling, reducing failed lookups
