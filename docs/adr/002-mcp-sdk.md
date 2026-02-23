# ADR-002: Official MCP TypeScript SDK

## Status

Accepted

## Context

We need an SDK to implement the MCP server protocol. The server must communicate over stdio with Claude Code, register tools, and handle tool call requests.

### Options Considered

1. **`@modelcontextprotocol/sdk`** — official Anthropic-maintained TypeScript SDK
2. **Raw protocol implementation** — implement JSON-RPC over stdio manually
3. **Third-party wrappers** — community libraries like `fastmcp`

## Decision

Use **`@modelcontextprotocol/sdk`** (the official MCP TypeScript SDK).

## Rationale

- **Official and maintained** — backed by Anthropic, guaranteed to stay compatible with Claude Code
- **Battle-tested** — the reference implementation for the MCP protocol; used in Anthropic's own examples
- **Type-safe** — provides TypeScript types for tool definitions, request/response schemas, and content types
- **Stdio transport built-in** — `StdioServerTransport` handles the Claude Code communication channel out of the box
- **Consistent with ecosystem** — `due-draghi-combattimenti` uses the Go equivalent (`modelcontextprotocol/go-sdk`); same patterns, different language
- **Zod integration** — input schemas can be defined with Zod for runtime validation and automatic JSON Schema generation

## Consequences

- Tied to the official SDK's API surface; breaking changes in the SDK require updates
- Need to follow the SDK's patterns for tool registration (which are well-documented)
- Zod becomes a transitive dependency (acceptable — it's lightweight and useful for validation)
