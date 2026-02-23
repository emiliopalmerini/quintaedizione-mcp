# ADR-003: Reuse JSON Data from quintaedizione.online

## Status

Accepted

## Context

The D&D SRD 2024 content already exists as structured JSON files in the `quintaedizione.online` project under `data/ita/json/`. These files contain Italian translations of the SRD covering 12 content types:

| File | Content |
|------|---------|
| `spells.json` | Spells (level, school, classes, components, description) |
| `monsters.json` | Monsters (stats, traits, actions, legendary actions) |
| `classes.json` | Character classes (features, hit die, proficiencies) |
| `equipment.json` | Equipment (weapons, armor, gear with properties) |
| `magic_items.json` | Magic items (type, rarity, attunement, description) |
| `feats.json` | Feats (category, prerequisite, benefit) |
| `backgrounds.json` | Backgrounds (ability scores, feat, skills) |
| `species.json` | Playable species (traits, creature type, size) |
| `glossary.json` | Game term definitions |
| `rules_gameplay.json` | Gameplay rules (hierarchical) |
| `rules_creation.json` | Character creation rules (hierarchical) |
| `rules_tools.json` | Tool-related rules (hierarchical) |

We need to decide how the MCP server accesses this data.

### Options Considered

1. **Copy JSON files into MCP project** — duplicate the data, bundle it with the server
2. **Symlink to quintaedizione.online data** — reference files at build/run time
3. **Git submodule** — include quintaedizione.online as a submodule
4. **Shared data package** — extract JSON into a standalone npm/Go package

## Decision

**Copy the JSON files** into the MCP project under `data/`.

## Rationale

- **Self-contained deployment** — the MCP server must work as a standalone Claude Code tool; it cannot depend on another project being cloned at a specific path
- **Stability** — the SRD content is static reference material that changes infrequently (only when upstream PDF is re-parsed); copying is a one-time operation with rare updates
- **Simplicity** — no submodule complexity, no symlink fragility, no extra package to maintain
- **Data contract** — the JSON schema is the contract between projects; as long as the schema is respected, the MCP server is decoupled from the Go project's internals
- **Update workflow** — when `quintaedizione.online` re-parses the SRD, copy the updated JSON files and commit; a simple script (`scripts/sync-data.sh`) can automate this

## Consequences

- Data is duplicated across repositories (acceptable given SRD content stability)
- Must manually sync when upstream JSON changes (mitigated by sync script)
- JSON schema changes in `quintaedizione.online` require corresponding type updates in the MCP project
