package main

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/ninckblokje/csheet/internal/csheet"
)

type request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

func serve(in *os.File, out *os.File) {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 1024*1024), 16*1024*1024)
	enc := json.NewEncoder(out)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var req request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			enc.Encode(response{JSONRPC: "2.0", ID: json.RawMessage("null"), Error: &rpcError{-32700, "parse error"}})
			continue
		}
		if len(req.ID) == 0 {
			continue
		}
		enc.Encode(handle(req))
	}
}

func handle(req request) response {
	resp := response{JSONRPC: "2.0", ID: req.ID}

	switch req.Method {
	case "initialize":
		resp.Result = map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{"tools": map[string]any{}},
			"serverInfo":      map[string]any{"name": "csheet-mcp", "version": version},
		}
	case "ping":
		resp.Result = map[string]any{}
	case "tools/list":
		resp.Result = map[string]any{"tools": tools()}
	case "tools/call":
		var p struct {
			Name      string            `json:"name"`
			Arguments map[string]string `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &p); err != nil {
			resp.Error = &rpcError{-32602, "invalid params"}
			return resp
		}
		resp.Result = callTool(p.Name, p.Arguments)
	default:
		resp.Error = &rpcError{-32601, "method not found"}
	}
	return resp
}

func tools() []map[string]any {
	str := map[string]any{"type": "string"}
	return []map[string]any{
		{
			"name":        "list_entries",
			"description": "List all cheat sheet entries as 'subject section', optionally filtered by subject",
			"inputSchema": map[string]any{
				"type":       "object",
				"properties": map[string]any{"subject": str},
			},
		},
		{
			"name":        "get_entry",
			"description": "Get the content of a cheat sheet entry",
			"inputSchema": map[string]any{
				"type":       "object",
				"properties": map[string]any{"subject": str, "section": str},
				"required":   []string{"subject", "section"},
			},
		},
	}
}

func textResult(text string, isError bool) map[string]any {
	return map[string]any{
		"content": []map[string]any{{"type": "text", "text": text}},
		"isError": isError,
	}
}

func callTool(name string, args map[string]string) map[string]any {
	fp, err := os.Open(csheetFile)
	if err != nil {
		return textResult(err.Error(), true)
	}
	defer fp.Close()

	switch name {
	case "list_entries":
		entries := csheet.FindEntries(fp)
		if s := args["subject"]; s != "" {
			entries = csheet.FilterEntries(entries, s)
		}
		return textResult(strings.Join(entries, "\n"), false)
	case "get_entry":
		if args["subject"] == "" || args["section"] == "" {
			return textResult("subject and section are required", true)
		}
		return textResult(strings.Join(csheet.FindEntry(fp, args["subject"], args["section"]), "\n"), false)
	default:
		return textResult("unknown tool: "+name, true)
	}
}
