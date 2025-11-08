package tool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/ninckblokje/csheet/internal/csheet"
)

var CSheetFile string = csheet.GetCSheetDir() + string(os.PathSeparator) + "csheet.md"

func getEntriesHandler(ctx context.Context, request mcp.CallToolRequest, args struct{}) (csheet.Entries, error) {
	entries := csheet.GetEntries(CSheetFile)

	var w bytes.Buffer
	json.NewEncoder(&w).Encode(entries)
	json := w.String()
	fmt.Println(json)

	return entries, nil
}

func getEntriesTool() mcp.Tool {
	return mcp.NewTool(
		"getEntries",
		mcp.WithDescription("Get all entries in the cheat sheet"),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithOutputSchema[csheet.Entries](),
	)
}

func getEntryHandler(ctx context.Context, request mcp.CallToolRequest, args struct{}) ([]string, error) {
	subject, _ := request.RequireString("subject")
	section, _ := request.RequireString("section")

	entry := csheet.GetEntry(CSheetFile, subject, section)

	return entry, nil
}

func getEntryTool() mcp.Tool {
	return mcp.NewTool(
		"getEntryTool",
		mcp.WithDescription("Get a single entry from the cheat sheet"),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("subject",
			mcp.Required(),
			mcp.Description("The subject of the entry to retrieve"),
		),
		mcp.WithString("section",
			mcp.Required(),
			mcp.Description("The section of the entry to retrieve"),
		),
	)
}
