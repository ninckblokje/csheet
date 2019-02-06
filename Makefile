PLATFORMS = linux/amd64/ windows/amd64/.exe

GIT_SHORT_REV := $(shell git rev-parse --short HEAD)
GIT_TAG := $(shell git describe --tags)

ifeq ($(OS),Windows_NT)
	RM_CMD := cmd.exe /c 'if exist "bin" rmdir /S /Q "bin"'
else
	RM_CMD := rm -rf bin/
endif

clean:
	$(RM_CMD)

build:
	go build csheet.go

dependencies:
	go get -d

release: clean dependencies $(PLATFORMS)

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
ext = $(word 3, $(temp))

$(PLATFORMS):
	GOARCH=$(arch) GOOS=$(os) go build -ldflags "-X main.csheetVersion=$(GIT_TAG) -X main.csheetRevision=$(GIT_SHORT_REV)" -o bin/$(os)_$(arch)/csheet$(ext) csheet.go
