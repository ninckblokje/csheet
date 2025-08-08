package main

import (
	"flag"
	"log"
	"os"

	"github.com/ninckblokje/csheet/internal/mcp/server"
	"github.com/ninckblokje/csheet/internal/mcp/tool"
)

var version = "DEV-BUILD"

func main() {
	var fileArg = flag.String("f", "", "Cheat sheet Mardown file")

	flag.Parse()

	if *fileArg != "" {
		tool.CSheetFile = *fileArg
	}

	logFile, err := os.OpenFile("csheet-mcp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.Default()

	mcpServer := server.CreateServer(version)
	server.StartStdioServer(mcpServer)
}
