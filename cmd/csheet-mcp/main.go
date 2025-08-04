package main

import "github.com/ninckblokje/csheet/internal/mcp/server"

var version = "DEV-BUILD"

func main() {
	mcpServer := server.CreateServer(version)
	server.StartStdioServer(mcpServer)
}
