# ADR-002: MCP Go SDK

## Status

Accepted (supersedes previous TypeScript SDK decision)

## Context

We need an MCP SDK for Go to handle protocol negotiation, tool registration, and stdio transport.

### Options Considered

1. **`github.com/mark3labs/mcp-go`** — most widely adopted Go MCP SDK
2. **Raw protocol implementation** — implement JSON-RPC over stdio manually

## Decision

Use **`github.com/mark3labs/mcp-go`**.

## Rationale

- Mature, well-maintained, widely used in the Go MCP ecosystem
- Supports stdio transport (required for Claude Code subprocess spawning)
- Clean tool registration API with typed parameters
- Handles protocol negotiation and JSON-RPC automatically

## Consequences

- Tied to the mcp-go SDK's API surface
- Need to follow the SDK's patterns for tool registration
