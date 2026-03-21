package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/emiliopalmerini/quintaedizione-data-ita/store"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerClasses(s *server.MCPServer, data *store.Store) {
	s.AddTool(mcp.NewTool("list_classes",
		mcp.WithDescription("List all D&D 5e character classes. Content is in Italian."),
	), listClassesHandler(data))

	s.AddTool(mcp.NewTool("get_class",
		mcp.WithDescription("Get full details of a D&D 5e class by ID. Content is in Italian."),
		mcp.WithString("id", mcp.Required(), mcp.Description("Class ID (slug)")),
	), getClassHandler(data))
}

func listClassesHandler(data *store.Store) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		classes := data.Classes()
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d class(es):\n\n", len(classes))
		for _, c := range classes {
			fmt.Fprintf(&sb, "- **%s** [ID: %s] — Dado Vita: %s\n", c.Name, c.ID, c.HitDie)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func getClassHandler(data *store.Store) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, _ := req.GetArguments()["id"].(string)
		if id == "" {
			return mcp.NewToolResultError("id is required"), nil
		}

		c, err := data.Class(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Class not found: %s", id)), nil
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "# %s\n\n", c.Name)
		fmt.Fprintf(&sb, "**Dado Vita:** %s\n", c.HitDie)
		if c.Proficiencies != "" {
			fmt.Fprintf(&sb, "**Competenze:** %s\n", c.Proficiencies)
		}
		if c.Description != "" {
			fmt.Fprintf(&sb, "\n%s\n", c.Description)
		}
		if len(c.Features) > 0 {
			sb.WriteString("\n## Privilegi\n\n")
			for _, f := range c.Features {
				fmt.Fprintf(&sb, "**%s (Livello %d).** %s\n\n", f.Name, f.Level, f.Description)
			}
		}
		if len(c.Subclasses) > 0 {
			sb.WriteString("## Sottoclassi\n\n")
			for _, sc := range c.Subclasses {
				fmt.Fprintf(&sb, "### %s\n\n%s\n\n", sc.Name, sc.Description)
			}
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}
