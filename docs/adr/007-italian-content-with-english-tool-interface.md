# ADR-007: Italian Content with English Tool Interface

## Status

Accepted

## Context

The SRD data from `quintaedizione.online` is in Italian (content text, field names in JSON use English keys but values are Italian). The JSON schema uses English field names (`name`, `level`, `school`) while content values are Italian (`"Palla di fuoco"`, `"Evocazione"`).

We need to decide the language of:
1. Tool names and descriptions (what Claude sees in the tool list)
2. Input parameter names and schema descriptions
3. Filter enum values (e.g., spell schools, monster types)
4. Response field names

### Options Considered

1. **All English** — translate filter values to English, English tool descriptions
2. **All Italian** — Italian tool names, Italian parameter names
3. **English interface, Italian content** — English tool/parameter names, Italian data values

## Decision

Use **English for the tool interface** (names, descriptions, parameter names) and **Italian for content values** (as stored in the JSON data).

## Rationale

- **Claude Code operates in English by default** — English tool names and descriptions give Claude the best understanding of what each tool does and when to use it
- **Data integrity** — the SRD content is Italian; translating filter values would create a mapping layer that adds complexity and potential for errors
- **User expectation** — users of this MCP server are playing D&D in Italian; they expect Italian content in responses
- **Filter values match data** — when Claude filters by `school: "Evocazione"`, it uses the same value that appears in the data. No translation lookup needed.
- **Discoverability** — tool descriptions document what Italian values are valid for each filter (e.g., "School of magic. Values: Abiurazione, Ammaliamento, Divinazione, Evocazione, Illusione, Invocazione, Necromanzia, Trasmutazione")

### Example

```
Tool: search_spells
Description: "Search D&D 5e SRD spells (Italian content from quintaedizione.online)"
Parameters:
  - query (string): "Search term to match against spell name"
  - level (number): "Spell level (0 for cantrips, 1-9)"
  - school (string): "School of magic: Abiurazione, Ammaliamento, Divinazione, Evocazione, Illusione, Invocazione, Necromanzia, Trasmutazione"
  - class (string): "Class name that can cast the spell (e.g., Mago, Chierico, Bardo)"
  - ritual (boolean): "Filter for ritual spells only"
```

## Consequences

- Claude must pass Italian values for filters (tool descriptions list valid values to guide it)
- Users interacting in English may see Italian content; Claude can translate on the fly if needed
- No translation layer to maintain between English filter values and Italian data
