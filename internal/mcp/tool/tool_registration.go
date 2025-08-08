package tool

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterTools(mcpServer *server.MCPServer) {
	mcpServer.AddTool(getEntriesTool(), mcp.NewStructuredToolHandler(getEntriesHandler))
}
