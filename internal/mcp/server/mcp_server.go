package server

import (
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/server"
	"github.com/ninckblokje/csheet/internal/mcp/tool"
)

// TODO : Move version info to file
func CreateServer(version string) *server.MCPServer {
	s := server.NewMCPServer(
		"CSheet MCP Server",
		version,
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	tool.RegisterTools(s)

	return s
}

func StartStdioServer(s *server.MCPServer) {
	if err := server.ServeStdio(s, server.WithErrorLogger(log.Default())); err != nil {
		fmt.Printf("MCP stdio server error: %v\n", err)
	}
}
