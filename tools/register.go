package tools

import (
	"github.com/emiliopalmerini/quintaedizione-data-ita/store"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterAll(s *server.MCPServer, data *store.Store) {
	registerSpells(s, data)
	registerMonsters(s, data)
	registerClasses(s, data)
	registerEquipment(s, data)
	registerMagicItems(s, data)
	registerFeats(s, data)
	registerBackgrounds(s, data)
	registerSpecies(s, data)
	registerRules(s, data)
	registerGlossary(s, data)
	registerMaps(s, data)
	registerGenerators(s, data)
	registerSearch(s, data)
}
