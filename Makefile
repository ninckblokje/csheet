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
	go mod download
	go mod verify

release: clean dependencies $(PLATFORMS)

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
ext = $(word 3, $(temp))

$(PLATFORMS):
	GOARCH=$(arch) GOOS=$(os) go build -ldflags "-X main.csheetVersion=$(GIT_TAG) -X main.csheetRevision=$(GIT_SHORT_REV)" -o bin/$(os)_$(arch)/csheet$(ext) csheet.go

DEB_PKG_ROOT = bin/debian_amd64/packageroot
DEB_INSTALL_DIR = /usr/bin

release-debian: clean dependencies linux/amd64/
	mkdir -p $(DEB_PKG_ROOT)$(DEB_INSTALL_DIR)
	mkdir -p $(DEB_PKG_ROOT)/DEBIAN
	cat pkg/DEBIAN/control | sed s/GIT_TAG/$(GIT_TAG)/g > $(DEB_PKG_ROOT)/DEBIAN/control
	cp bin/linux_amd64/csheet $(DEB_PKG_ROOT)$(DEB_INSTALL_DIR)
	dpkg-deb -b $(DEB_PKG_ROOT) bin/debian_amd64/csheet_$(GIT_TAG)_amd64.deb