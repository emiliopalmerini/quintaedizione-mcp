package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/emiliopalmerini/quintaedizione-data-ita/store"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerMagicItems(s *server.MCPServer, data *store.Store) {
	s.AddTool(mcp.NewTool("search_magic_items",
		mcp.WithDescription("Search D&D 5e magic items by name, type, or rarity. Content is in Italian."),
		mcp.WithString("query", mcp.Description("Search term")),
		mcp.WithString("rarity", mcp.Description("Filter by rarity (e.g. Comune, Non comune, Raro, Molto raro, Leggendario)")),
		mcp.WithString("type", mcp.Description("Filter by type")),
	), searchMagicItemsHandler(data))

	s.AddTool(mcp.NewTool("get_magic_item",
		mcp.WithDescription("Get full details of a D&D 5e magic item by ID. Content is in Italian."),
		mcp.WithString("id", mcp.Required(), mcp.Description("Magic item ID (slug)")),
	), getMagicItemHandler(data))
}

func searchMagicItemsHandler(data *store.Store) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, _ := req.GetArguments()["query"].(string)
		rarity, _ := req.GetArguments()["rarity"].(string)
		typ, _ := req.GetArguments()["type"].(string)

		var results []string
		for _, mi := range data.MagicItems() {
			if query != "" && !containsI(mi.Name, query) {
				continue
			}
			if rarity != "" && !containsI(mi.Rarity, rarity) {
				continue
			}
			if typ != "" && !containsI(mi.Type, typ) {
				continue
			}
			results = append(results, fmt.Sprintf("- **%s** [ID: %s] — %s, %s", mi.Name, mi.ID, mi.Type, mi.Rarity))
			if len(results) >= 20 {
				break
			}
		}
		if len(results) == 0 {
			return mcp.NewToolResultText("No magic items found."), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Found %d item(s):\n\n%s", len(results), strings.Join(results, "\n"))), nil
	}
}

func getMagicItemHandler(data *store.Store) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, _ := req.GetArguments()["id"].(string)
		if id == "" {
			return mcp.NewToolResultError("id is required"), nil
		}
		mi, err := data.MagicItem(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Magic item not found: %s", id)), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "# %s\n\n**Tipo:** %s | **Rarità:** %s\n", mi.Name, mi.Type, mi.Rarity)
		if mi.Attunement {
			att := "Richiede sintonia"
			if mi.AttunementDetails != "" {
				att += " (" + mi.AttunementDetails + ")"
			}
			fmt.Fprintf(&sb, "**%s**\n", att)
		}
		if mi.Description != "" {
			fmt.Fprintf(&sb, "\n%s", mi.Description)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}
