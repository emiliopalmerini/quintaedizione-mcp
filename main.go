package main

import (
	"fmt"
	"log"

	"github.com/emiliopalmerini/quintaedizione-data-ita/store"
	"github.com/emiliopalmerini/quintaedizione-mcp/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s, err := store.Load()
	if err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}

	mcpServer := server.NewMCPServer(
		"quintaedizione-mcp",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	tools.RegisterAll(mcpServer, s)

	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
