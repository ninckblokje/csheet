package tool

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func findEntriesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText("Hello World"), nil
}

func findEntriesTool() mcp.Tool {
	return mcp.NewTool(
		"find_entries",
		mcp.WithDescription("Find all entries in the cheat sheet"),
	)
}
