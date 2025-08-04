package tool

import "github.com/mark3labs/mcp-go/server"

func RegisterTools(mcpServer *server.MCPServer) {
	mcpServer.AddTool(findEntriesTool(), findEntriesHandler)
}
