package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/emiliopalmerini/quintaedizione-data-ita/store"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerEquipment(s *server.MCPServer, data *store.Store) {
	s.AddTool(mcp.NewTool("search_equipment",
		mcp.WithDescription("Search D&D 5e equipment by name or category. Content is in Italian."),
		mcp.WithString("query", mcp.Description("Search term")),
		mcp.WithString("category", mcp.Description("Filter by subcategory")),
	), searchEquipmentHandler(data))

	s.AddTool(mcp.NewTool("get_equipment",
		mcp.WithDescription("Get full details of a D&D 5e equipment item by ID. Content is in Italian."),
		mcp.WithString("id", mcp.Required(), mcp.Description("Equipment ID (slug)")),
	), getEquipmentHandler(data))
}

func searchEquipmentHandler(data *store.Store) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, _ := req.GetArguments()["query"].(string)
		category, _ := req.GetArguments()["category"].(string)

		var results []string
		for _, e := range data.Equipment() {
			if query != "" && !containsI(e.Name, query) {
				continue
			}
			if category != "" && !containsI(e.Subcategory, category) {
				continue
			}
			results = append(results, fmt.Sprintf("- **%s** [ID: %s] — %s", e.Name, e.ID, e.Subcategory))
			if len(results) >= 20 {
				break
			}
		}
		if len(results) == 0 {
			return mcp.NewToolResultText("No equipment found."), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Found %d item(s):\n\n%s", len(results), strings.Join(results, "\n"))), nil
	}
}

func getEquipmentHandler(data *store.Store) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, _ := req.GetArguments()["id"].(string)
		if id == "" {
			return mcp.NewToolResultError("id is required"), nil
		}
		e, err := data.EquipmentItem(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Equipment not found: %s", id)), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "# %s\n\n**Categoria:** %s\n", e.Name, e.Subcategory)
		for k, v := range e.Properties {
			fmt.Fprintf(&sb, "**%s:** %s\n", k, v)
		}
		if e.Description != "" {
			fmt.Fprintf(&sb, "\n%s", e.Description)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}
