package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/emiliopalmerini/quintaedizione-data-ita/srd"
	"github.com/emiliopalmerini/quintaedizione-data-ita/store"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerMonsters(s *server.MCPServer, data *store.Store) {
	s.AddTool(mcp.NewTool("search_monsters",
		mcp.WithDescription("Search D&D 5e monsters by name, type, size, or CR. Returns summaries. Content is in Italian."),
		mcp.WithString("query", mcp.Description("Search term (monster name)")),
		mcp.WithString("type", mcp.Description("Filter by type (e.g. Drago, Non morto, Aberrazione)")),
		mcp.WithString("size", mcp.Description("Filter by size (e.g. Grande, Media, Enorme)")),
		mcp.WithString("cr", mcp.Description("Filter by challenge rating (e.g. 5, 1/4)")),
	), searchMonstersHandler(data))

	s.AddTool(mcp.NewTool("get_monster",
		mcp.WithDescription("Get full details of a D&D 5e monster by ID. Content is in Italian."),
		mcp.WithString("id", mcp.Required(), mcp.Description("Monster ID (slug)")),
	), getMonsterHandler(data))
}

func searchMonstersHandler(data *store.Store) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, _ := req.GetArguments()["query"].(string)
		typ, _ := req.GetArguments()["type"].(string)
		size, _ := req.GetArguments()["size"].(string)
		cr, _ := req.GetArguments()["cr"].(string)

		var results []srd.Monster
		for _, m := range data.Monsters() {
			if query != "" && !containsI(m.Name, query) {
				continue
			}
			if typ != "" && !containsI(m.Type, typ) {
				continue
			}
			if size != "" && !containsI(m.Size, size) {
				continue
			}
			if cr != "" && m.CR != cr {
				continue
			}
			results = append(results, m)
			if len(results) >= 20 {
				break
			}
		}

		if len(results) == 0 {
			return mcp.NewToolResultText("No monsters found."), nil
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "Found %d monster(s):\n\n", len(results))
		for _, m := range results {
			fmt.Fprintf(&sb, "- **%s** [ID: %s] — %s, %s, GS %s\n",
				m.Name, m.ID, m.Type, m.Size, m.CR)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func getMonsterHandler(data *store.Store) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, _ := req.GetArguments()["id"].(string)
		if id == "" {
			return mcp.NewToolResultError("id is required"), nil
		}

		m, err := data.Monster(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Monster not found: %s", id)), nil
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "# %s\n\n", m.Name)
		fmt.Fprintf(&sb, "*%s %s, %s*\n\n", m.Size, m.Type, m.Alignment)
		fmt.Fprintf(&sb, "**CA:** %s | **PF:** %s | **Velocità:** %s\n", m.AC, m.HP, m.Speed)
		fmt.Fprintf(&sb, "**GS:** %s", m.CR)
		if m.CRDetail != "" {
			fmt.Fprintf(&sb, " (%s)", m.CRDetail)
		}
		sb.WriteString("\n\n")

		if len(m.AbilityScores) > 0 {
			sb.WriteString("| FOR | DES | COS | INT | SAG | CAR |\n|-----|-----|-----|-----|-----|-----|\n")
			fmt.Fprintf(&sb, "| %d | %d | %d | %d | %d | %d |\n\n",
				m.AbilityScores["strength"], m.AbilityScores["dexterity"],
				m.AbilityScores["constitution"], m.AbilityScores["intelligence"],
				m.AbilityScores["wisdom"], m.AbilityScores["charisma"])
		}

		if m.Skills != "" {
			fmt.Fprintf(&sb, "**Abilità:** %s\n", m.Skills)
		}
		if m.Senses != "" {
			fmt.Fprintf(&sb, "**Sensi:** %s\n", m.Senses)
		}
		if m.Languages != "" {
			fmt.Fprintf(&sb, "**Lingue:** %s\n", m.Languages)
		}
		sb.WriteString("\n")

		writeFeatures(&sb, "Tratti", m.Traits)
		writeFeatures(&sb, "Azioni", m.Actions)
		writeFeatures(&sb, "Azioni Bonus", m.BonusActions)
		writeFeatures(&sb, "Reazioni", m.Reactions)
		writeFeatures(&sb, "Azioni Leggendarie", m.LegendaryActions)

		return mcp.NewToolResultText(sb.String()), nil
	}
}

func writeFeatures(sb *strings.Builder, title string, features []srd.Feature) {
	if len(features) == 0 {
		return
	}
	fmt.Fprintf(sb, "## %s\n\n", title)
	for _, f := range features {
		fmt.Fprintf(sb, "**%s.** %s\n\n", f.Name, f.Description)
	}
}
