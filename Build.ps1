$ShortRevision = git rev-parse --short HEAD
go build -ldflags "-X main.csheetVersion=$ShortRevision" csheet.go