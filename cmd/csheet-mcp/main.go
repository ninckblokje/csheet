// Command csheet-mcp is a Model Context Protocol server (stdio transport)
// exposing the csheet cheat sheets as tools.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

var commit = "test"
var version = "DEV-BUILD"

var csheetFile string

func main() {
	fileArg := flag.String("f", "", "Cheat sheet Markdown file")
	versionArg := flag.Bool("v", false, "Display version")
	flag.Parse()

	if *versionArg {
		fmt.Printf("csheet-mcp version v%s, revision %s\n", version, commit)
		return
	}

	if *fileArg != "" {
		csheetFile = *fileArg
	} else {
		usr, err := user.Current()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		csheetFile = filepath.Join(usr.HomeDir, "csheet.md")
	}

	serve(os.Stdin, os.Stdout)
}
