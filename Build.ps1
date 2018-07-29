$ShortRevision = git rev-parse --short HEAD
$Tag = git describe --tags
go build -ldflags "-X main.csheetVersion=$Tag -X main.csheetRevision=$ShortRevision" csheet.go