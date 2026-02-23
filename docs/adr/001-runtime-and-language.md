# ADR-001: Bun as Runtime and TypeScript as Language

## Status

Accepted

## Context

We need to build an MCP (Model Context Protocol) server that exposes D&D SRD 2024 content to Claude Code. The existing ecosystem includes:

- **quintaedizione.online** — a Go + Templ web viewer serving Italian SRD content from embedded JSON
- **due-draghi-combattimenti** — a Go encounter calculator with an MCP server using `github.com/modelcontextprotocol/go-sdk`

The MCP protocol originated from Anthropic and the reference SDK is TypeScript-first. We need to choose a runtime and language for the new server.

### Options Considered

1. **Go** — consistent with existing projects; uses `modelcontextprotocol/go-sdk`
2. **Node.js + TypeScript** — reference MCP SDK is TypeScript; large ecosystem
3. **Bun + TypeScript** — fast TypeScript runtime with native TS support, built-in test runner, no transpile step

## Decision

Use **Bun** as the runtime and **TypeScript** as the language.

## Rationale

- **First-class TypeScript** — Bun runs `.ts` files natively with no compilation step, reducing toolchain complexity (no `tsc`, `tsconfig` build paths, `ts-node`)
- **MCP SDK alignment** — the official `@modelcontextprotocol/sdk` is written in TypeScript; using TS gives us the best type safety and SDK integration
- **Fast startup** — Bun's startup time is significantly lower than Node.js, which matters for an MCP server that Claude Code spawns as a subprocess on demand
- **Built-in tooling** — `bun test`, `bun install`, and native JSON imports reduce external dependencies
- **JSON-heavy workload** — the server loads and serves JSON data; Bun's optimized JSON parsing is a natural fit
- **Diversification** — the existing projects are Go; using TypeScript for the MCP server brings variety to the stack and leverages the stronger MCP TypeScript ecosystem

## Consequences

- Team needs TypeScript familiarity (already common in web development)
- Bun is newer than Node.js; some npm packages may have edge-case incompatibilities (mitigated by the small dependency surface of this project)
- Cannot share code directly with the Go projects (acceptable — the data contract is the shared JSON files)
