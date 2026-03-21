# ADR-001: Go as Runtime and Language

## Status

Accepted (supersedes previous Bun/TypeScript decision)

## Context

We need to build an MCP server that exposes D&D SRD 2024 content to Claude Code. The `quintaedizione-data-ita` shared Go module provides typed structs, embedded JSON, search, and filtering — ready to consume.

### Options Considered

1. **Go** — direct import of shared module, single binary, consistent ecosystem
2. **Bun + TypeScript** — reference MCP SDK, but requires copying JSON data

## Decision

Use **Go** as the language and compile to a single binary.

## Rationale

- **Direct import** of `quintaedizione-data-ita` — no data duplication, typed access to all SRD entities
- **Single static binary** — trivial deployment, no runtime dependencies
- **Consistent ecosystem** — same language as quintaedizione.online, same team knowledge
- **MCP Go SDK** available (`github.com/mark3labs/mcp-go`)
- **Fast startup**, low memory footprint

## Consequences

- Cannot use the official TypeScript MCP SDK (acceptable — Go SDK is mature)
- Shared module updates propagate via `go get -u`
