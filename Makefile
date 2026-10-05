.PHONY: build

build:
	mkdir -p bin
	go build -o bin/csheet ./cmd/csheet
	go build -o bin/csheet-mcp ./cmd/csheet-mcp
