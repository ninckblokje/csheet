package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const testSheet = "# csheet\n\n" +
	"## go\n\n" +
	"### build\n\n" +
	"````\ngo build ./...\n````\n\n" +
	"### test\n\n" +
	"````\ngo test ./...\n````\n\n" +
	"## git\n\n" +
	"### status\n\n" +
	"````\ngit status\n````\n"

func TestHandle(t *testing.T) {
	tests := []struct {
		name      string
		request   request
		wantError int
		wantTools int
	}{
		{
			name:    "initialize",
			request: request{ID: json.RawMessage("1"), Method: "initialize"},
		},
		{
			name:    "ping",
			request: request{ID: json.RawMessage("2"), Method: "ping"},
		},
		{
			name:      "tools list",
			request:   request{ID: json.RawMessage("3"), Method: "tools/list"},
			wantTools: 2,
		},
		{
			name:      "invalid tool params",
			request:   request{ID: json.RawMessage("4"), Method: "tools/call", Params: json.RawMessage("{")},
			wantError: -32602,
		},
		{
			name:      "unknown method",
			request:   request{ID: json.RawMessage("5"), Method: "unknown"},
			wantError: -32601,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handle(tt.request)
			if got.JSONRPC != "2.0" {
				t.Errorf("JSONRPC = %q, want %q", got.JSONRPC, "2.0")
			}
			if !reflect.DeepEqual(got.ID, tt.request.ID) {
				t.Errorf("ID = %s, want %s", got.ID, tt.request.ID)
			}
			if tt.wantError != 0 {
				if got.Error == nil || got.Error.Code != tt.wantError {
					t.Fatalf("Error = %v, want code %d", got.Error, tt.wantError)
				}
				return
			}
			if got.Error != nil {
				t.Fatalf("unexpected error: %+v", got.Error)
			}
			if tt.wantTools > 0 {
				result, ok := got.Result.(map[string]any)
				if !ok {
					t.Fatalf("Result has type %T, want map[string]any", got.Result)
				}
				tools, ok := result["tools"].([]map[string]any)
				if !ok {
					t.Fatalf("tools has type %T, want []map[string]any", result["tools"])
				}
				if len(tools) != tt.wantTools {
					t.Errorf("got %d tools, want %d", len(tools), tt.wantTools)
				}
			}
		})
	}
}

func TestServe(t *testing.T) {
	dir := t.TempDir()
	inPath := filepath.Join(dir, "input")
	outPath := filepath.Join(dir, "output")
	input := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"ping"}`,
		`{"jsonrpc":"2.0","method":"ping"}`,
		`invalid json`,
	}, "\n")
	if err := os.WriteFile(inPath, []byte(input), 0600); err != nil {
		t.Fatal(err)
	}

	in, err := os.Open(inPath)
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	out, err := os.Create(outPath)
	if err != nil {
		t.Fatal(err)
	}
	serve(in, out)
	if err := out.Close(); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d responses, want 2: %s", len(lines), data)
	}

	var ping response
	if err := json.Unmarshal([]byte(lines[0]), &ping); err != nil {
		t.Fatal(err)
	}
	if string(ping.ID) != "1" || ping.Error != nil {
		t.Errorf("ping response = %+v, want id 1 and no error", ping)
	}

	var parseError response
	if err := json.Unmarshal([]byte(lines[1]), &parseError); err != nil {
		t.Fatal(err)
	}
	if parseError.Error == nil || parseError.Error.Code != -32700 {
		t.Errorf("parse error response = %+v, want error code -32700", parseError)
	}
}

func TestCallTool(t *testing.T) {
	path := filepath.Join(t.TempDir(), "csheet.md")
	if err := os.WriteFile(path, []byte(testSheet), 0600); err != nil {
		t.Fatal(err)
	}
	oldPath := csheetFile
	csheetFile = path
	t.Cleanup(func() { csheetFile = oldPath })

	tests := []struct {
		name    string
		tool    string
		args    map[string]string
		want    string
		wantErr bool
	}{
		{
			name: "list all entries",
			tool: "list_entries",
			want: "go build\ngo test\ngit status",
		},
		{
			name: "filter entries",
			tool: "list_entries",
			args: map[string]string{"subject": "go"},
			want: "go build\ngo test",
		},
		{
			name: "get entry",
			tool: "get_entry",
			args: map[string]string{"subject": "go", "section": "build"},
			want: "go build ./...",
		},
		{
			name:    "require entry arguments",
			tool:    "get_entry",
			args:    map[string]string{"subject": "go"},
			want:    "subject and section are required",
			wantErr: true,
		},
		{
			name:    "unknown tool",
			tool:    "not_a_tool",
			want:    "unknown tool: not_a_tool",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := callTool(tt.tool, tt.args)
			if got["isError"] != tt.wantErr {
				t.Errorf("isError = %v, want %v", got["isError"], tt.wantErr)
			}
			content, ok := got["content"].([]map[string]any)
			if !ok || len(content) != 1 {
				t.Fatalf("content = %#v, want one text item", got["content"])
			}
			if text := content[0]["text"]; text != tt.want {
				t.Errorf("text = %q, want %q", text, tt.want)
			}
		})
	}
}

func TestCallToolMissingSheet(t *testing.T) {
	oldPath := csheetFile
	csheetFile = filepath.Join(t.TempDir(), "missing.md")
	t.Cleanup(func() { csheetFile = oldPath })

	got := callTool("list_entries", nil)
	if got["isError"] != true {
		t.Errorf("isError = %v, want true", got["isError"])
	}
}
