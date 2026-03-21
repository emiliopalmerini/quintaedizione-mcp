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

func registerSpells(s *server.MCPServer, data *store.Store) {
	s.AddTool(mcp.NewTool("search_spells",
		mcp.WithDescription("Search D&D 5e spells by name, school, level, or class. Returns summaries. Content is in Italian."),
		mcp.WithString("query", mcp.Description("Search term (spell name or keyword)")),
		mcp.WithString("school", mcp.Description("Filter by school (e.g. Evocazione, Abiurazione)")),
		mcp.WithNumber("level", mcp.Description("Filter by spell level (0-9)")),
		mcp.WithString("class", mcp.Description("Filter by class (e.g. Mago, Chierico)")),
	), searchSpellsHandler(data))

	s.AddTool(mcp.NewTool("get_spell",
		mcp.WithDescription("Get full details of a D&D 5e spell by ID. Content is in Italian."),
		mcp.WithString("id", mcp.Required(), mcp.Description("Spell ID (slug)")),
	), getSpellHandler(data))
}

func searchSpellsHandler(data *store.Store) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, _ := req.GetArguments()["query"].(string)
		school, _ := req.GetArguments()["school"].(string)
		levelVal, hasLevel := req.GetArguments()["level"].(float64)
		class, _ := req.GetArguments()["class"].(string)

		var results []srd.Spell
		for _, spell := range data.Spells() {
			if query != "" && !containsI(spell.Name, query) {
				continue
			}
			if school != "" && !containsI(spell.School, school) {
				continue
			}
			if hasLevel && spell.Level != int(levelVal) {
				continue
			}
			if class != "" && !containsI(strings.Join(spell.Classes, ","), class) {
				continue
			}
			results = append(results, spell)
			if len(results) >= 20 {
				break
			}
		}

		if len(results) == 0 {
			return mcp.NewToolResultText("No spells found."), nil
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "Found %d spell(s):\n\n", len(results))
		for _, s := range results {
			ritual := ""
			if s.Ritual {
				ritual = " (rituale)"
			}
			fmt.Fprintf(&sb, "- **%s** [ID: %s] — Livello %d, %s%s, Classi: %s\n",
				s.Name, s.ID, s.Level, s.School, ritual, strings.Join(s.Classes, ", "))
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func getSpellHandler(data *store.Store) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, _ := req.GetArguments()["id"].(string)
		if id == "" {
			return mcp.NewToolResultError("id is required"), nil
		}

		spell, err := data.Spell(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Spell not found: %s", id)), nil
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "# %s\n\n", spell.Name)
		fmt.Fprintf(&sb, "**Livello:** %d | **Scuola:** %s", spell.Level, spell.School)
		if spell.Ritual {
			sb.WriteString(" (rituale)")
		}
		fmt.Fprintf(&sb, "\n**Tempo di lancio:** %s\n", spell.CastingTime)
		fmt.Fprintf(&sb, "**Gittata:** %s\n", spell.Range)
		fmt.Fprintf(&sb, "**Componenti:** %s\n", spell.Components)
		fmt.Fprintf(&sb, "**Durata:** %s\n", spell.Duration)
		fmt.Fprintf(&sb, "**Classi:** %s\n\n", strings.Join(spell.Classes, ", "))
		sb.WriteString(spell.Description)
		if spell.AtHigherLevels != "" {
			fmt.Fprintf(&sb, "\n\n**Ai Livelli Superiori:** %s", spell.AtHigherLevels)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}
