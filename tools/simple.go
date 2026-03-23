package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/emiliopalmerini/quintaedizione-data-ita/store"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// --- Feats ---

func registerFeats(s *server.MCPServer, data *store.Store) {
	s.AddTool(mcp.NewTool("search_feats",
		mcp.WithDescription("Search D&D 5e feats by name or category. Content is in Italian."),
		mcp.WithString("query", mcp.Description("Search term")),
		mcp.WithString("category", mcp.Description("Filter by category")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, _ := req.GetArguments()["query"].(string)
		category, _ := req.GetArguments()["category"].(string)
		var results []string
		for _, f := range data.Feats() {
			if query != "" && !containsI(f.Name, query) {
				continue
			}
			if category != "" && !containsI(f.Category, category) {
				continue
			}
			results = append(results, fmt.Sprintf("- **%s** [ID: %s] — %s", f.Name, f.ID, f.Category))
			if len(results) >= 20 {
				break
			}
		}
		if len(results) == 0 {
			return mcp.NewToolResultText("No feats found."), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Found %d feat(s):\n\n%s", len(results), strings.Join(results, "\n"))), nil
	})

	s.AddTool(mcp.NewTool("get_feat",
		mcp.WithDescription("Get full details of a D&D 5e feat by ID. Content is in Italian."),
		mcp.WithString("id", mcp.Required(), mcp.Description("Feat ID (slug)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, _ := req.GetArguments()["id"].(string)
		if id == "" {
			return mcp.NewToolResultError("id is required"), nil
		}
		f, err := data.Feat(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Feat not found: %s", id)), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "# %s\n\n**Categoria:** %s\n", f.Name, f.Category)
		if len(f.Prerequisite) > 0 {
			fmt.Fprintf(&sb, "**Prerequisito:** %s\n", f.Prerequisite.PlainText())
		}
		if f.Repeatable {
			sb.WriteString("**Ripetibile:** Sì\n")
		}
		fmt.Fprintf(&sb, "\n%s", f.Benefit.PlainText())
		return mcp.NewToolResultText(sb.String()), nil
	})
}

// --- Backgrounds ---

func registerBackgrounds(s *server.MCPServer, data *store.Store) {
	s.AddTool(mcp.NewTool("list_backgrounds",
		mcp.WithDescription("List all D&D 5e backgrounds. Content is in Italian."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		bgs := data.Backgrounds()
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d background(s):\n\n", len(bgs))
		for _, b := range bgs {
			fmt.Fprintf(&sb, "- **%s** [ID: %s]\n", b.Name, b.ID)
		}
		return mcp.NewToolResultText(sb.String()), nil
	})

	s.AddTool(mcp.NewTool("get_background",
		mcp.WithDescription("Get full details of a D&D 5e background by ID. Content is in Italian."),
		mcp.WithString("id", mcp.Required(), mcp.Description("Background ID (slug)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, _ := req.GetArguments()["id"].(string)
		if id == "" {
			return mcp.NewToolResultError("id is required"), nil
		}
		b, err := data.Background(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Background not found: %s", id)), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "# %s\n\n", b.Name)
		if b.AbilityScores != "" {
			fmt.Fprintf(&sb, "**Punteggi di Caratteristica:** %s\n", b.AbilityScores)
		}
		if b.Feat != "" {
			fmt.Fprintf(&sb, "**Talento:** %s\n", b.Feat)
		}
		if b.SkillProficiencies != "" {
			fmt.Fprintf(&sb, "**Competenze:** %s\n", b.SkillProficiencies)
		}
		if b.Equipment != "" {
			fmt.Fprintf(&sb, "**Equipaggiamento:** %s\n", b.Equipment)
		}
		if len(b.Description) > 0 {
			fmt.Fprintf(&sb, "\n%s", b.Description.PlainText())
		}
		return mcp.NewToolResultText(sb.String()), nil
	})
}

// --- Species ---

func registerSpecies(s *server.MCPServer, data *store.Store) {
	s.AddTool(mcp.NewTool("list_species",
		mcp.WithDescription("List all D&D 5e playable species. Content is in Italian."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		species := data.Species()
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d species:\n\n", len(species))
		for _, sp := range species {
			fmt.Fprintf(&sb, "- **%s** [ID: %s] — %s, %s\n", sp.Name, sp.ID, sp.CreatureType, sp.Size)
		}
		return mcp.NewToolResultText(sb.String()), nil
	})

	s.AddTool(mcp.NewTool("get_species",
		mcp.WithDescription("Get full details of a D&D 5e species by ID. Content is in Italian."),
		mcp.WithString("id", mcp.Required(), mcp.Description("Species ID (slug)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, _ := req.GetArguments()["id"].(string)
		if id == "" {
			return mcp.NewToolResultError("id is required"), nil
		}
		sp, err := data.SpeciesEntry(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Species not found: %s", id)), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "# %s\n\n", sp.Name)
		fmt.Fprintf(&sb, "**Tipo:** %s | **Taglia:** %s | **Velocità:** %s\n\n", sp.CreatureType, sp.Size, sp.Speed)
		if len(sp.Description) > 0 {
			sb.WriteString(sp.Description.PlainText())
			sb.WriteString("\n\n")
		}
		for _, t := range sp.Traits {
			fmt.Fprintf(&sb, "**%s.** %s\n\n", t.Name, t.Description.PlainText())
		}
		return mcp.NewToolResultText(sb.String()), nil
	})
}

// --- Rules ---

func registerRules(s *server.MCPServer, data *store.Store) {
	s.AddTool(mcp.NewTool("search_rules",
		mcp.WithDescription("Search D&D 5e rules by keyword. Content is in Italian."),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search term")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, _ := req.GetArguments()["query"].(string)
		if query == "" {
			return mcp.NewToolResultError("query is required"), nil
		}
		var results []string
		for _, r := range data.Rules() {
			if containsI(r.Title, query) || containsI(r.Content.PlainText(), query) {
				results = append(results, fmt.Sprintf("- **%s** [ID: %s]", r.Title, r.ID))
				if len(results) >= 20 {
					break
				}
			}
		}
		if len(results) == 0 {
			return mcp.NewToolResultText("No rules found."), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Found %d rule(s):\n\n%s", len(results), strings.Join(results, "\n"))), nil
	})

	s.AddTool(mcp.NewTool("get_rule",
		mcp.WithDescription("Get full details of a D&D 5e rule by ID. Content is in Italian."),
		mcp.WithString("id", mcp.Required(), mcp.Description("Rule ID (slug)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, _ := req.GetArguments()["id"].(string)
		if id == "" {
			return mcp.NewToolResultError("id is required"), nil
		}
		r, err := data.Rule(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Rule not found: %s", id)), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "# %s\n\n%s", r.Title, r.Content.PlainText())
		return mcp.NewToolResultText(sb.String()), nil
	})
}

// --- Glossary ---

func registerGlossary(s *server.MCPServer, data *store.Store) {
	s.AddTool(mcp.NewTool("lookup_glossary",
		mcp.WithDescription("Look up a D&D 5e game term in the glossary. Content is in Italian."),
		mcp.WithString("term", mcp.Required(), mcp.Description("Term to look up")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		term, _ := req.GetArguments()["term"].(string)
		if term == "" {
			return mcp.NewToolResultError("term is required"), nil
		}
		var results []string
		for _, g := range data.Glossary() {
			if containsI(g.Term, term) {
				entry := fmt.Sprintf("**%s:** %s", g.Term, g.Definition.PlainText())
				if len(g.SeeAlso) > 0 {
					entry += fmt.Sprintf(" (Vedi anche: %s)", strings.Join(g.SeeAlso, ", "))
				}
				results = append(results, entry)
				if len(results) >= 10 {
					break
				}
			}
		}
		if len(results) == 0 {
			return mcp.NewToolResultText("Term not found in glossary."), nil
		}
		return mcp.NewToolResultText(strings.Join(results, "\n\n")), nil
	})
}

// --- Maps ---

func registerMaps(s *server.MCPServer, data *store.Store) {
	s.AddTool(mcp.NewTool("search_maps",
		mcp.WithDescription("Search D&D battle maps by name or tag. Content is in Italian."),
		mcp.WithString("query", mcp.Description("Search term")),
		mcp.WithString("tag", mcp.Description("Filter by tag (e.g. sotterraneo, foresta, città)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, _ := req.GetArguments()["query"].(string)
		tag, _ := req.GetArguments()["tag"].(string)
		var results []string
		for _, m := range data.Mappe() {
			if query != "" && !containsI(m.Nome, query) {
				continue
			}
			if tag != "" {
				found := false
				for _, t := range m.Tag {
					if containsI(t, tag) {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			results = append(results, fmt.Sprintf("- **%s** [%s] — %s, di %s",
				m.Nome, m.Slug, strings.Join(m.Tag, ", "), m.Autore))
			if len(results) >= 20 {
				break
			}
		}
		if len(results) == 0 {
			return mcp.NewToolResultText("No maps found."), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Found %d map(s):\n\n%s", len(results), strings.Join(results, "\n"))), nil
	})
}

// --- Generators ---

func registerGenerators(s *server.MCPServer, data *store.Store) {
	s.AddTool(mcp.NewTool("list_generators",
		mcp.WithDescription("List all available random generator tables. Content is in Italian."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tables := data.GeneratorTables()
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d generator table(s):\n\n", len(tables))
		currentGroup := ""
		for _, t := range tables {
			if t.Group != currentGroup {
				currentGroup = t.Group
				fmt.Fprintf(&sb, "\n### %s\n\n", currentGroup)
			}
			fmt.Fprintf(&sb, "- **%s** [ID: %s] — %s (%s)\n", t.Name, t.ID, t.Description, t.Die)
		}
		return mcp.NewToolResultText(sb.String()), nil
	})

	s.AddTool(mcp.NewTool("get_generator",
		mcp.WithDescription("Get a random generator table by ID. Content is in Italian."),
		mcp.WithString("id", mcp.Required(), mcp.Description("Generator table ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, _ := req.GetArguments()["id"].(string)
		if id == "" {
			return mcp.NewToolResultError("id is required"), nil
		}
		t, err := data.GeneratorTable(id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Generator not found: %s", id)), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "# %s\n\n%s\n**Dado:** %s\n\n", t.Name, t.Description, t.Die)
		if t.IsMultiColumn() {
			for _, col := range t.Columns {
				fmt.Fprintf(&sb, "### %s\n\n", col.Name)
				for i, item := range col.Items {
					fmt.Fprintf(&sb, "%d. %s\n", i+1, item.Text)
				}
				sb.WriteString("\n")
			}
		} else {
			for i, item := range t.Items {
				fmt.Fprintf(&sb, "%d. %s\n", i+1, item.Text)
			}
		}
		return mcp.NewToolResultText(sb.String()), nil
	})
}

// --- Search (cross-collection) ---

func registerSearch(s *server.MCPServer, data *store.Store) {
	s.AddTool(mcp.NewTool("search",
		mcp.WithDescription("Search across all D&D 5e content (spells, monsters, classes, etc.). Content is in Italian."),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search term")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, _ := req.GetArguments()["query"].(string)
		if query == "" {
			return mcp.NewToolResultError("query is required"), nil
		}
		results := data.Search(query, 5)
		if len(results) == 0 {
			return mcp.NewToolResultText("No results found."), nil
		}
		var sb strings.Builder
		for _, set := range results {
			fmt.Fprintf(&sb, "### %s (%d results)\n\n", set.Collection, set.Total)
			for _, r := range set.Results {
				fmt.Fprintf(&sb, "- **%s** [ID: %s]\n", r.Title, r.ID)
			}
			sb.WriteString("\n")
		}
		return mcp.NewToolResultText(sb.String()), nil
	})
}
