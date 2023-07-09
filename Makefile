PLATFORMS = linux/amd64/ linux/arm64/ windows/amd64/.exe windows/arm64/.exe freebsd/amd64/ freebsd/arm64 openbsd/amd64/ openbsd/arm64/

GIT_SHORT_REV := $(shell git rev-parse --short HEAD)
GIT_TAG := $(shell git describe --tags | sed s/v//g)

.DEFAULT_GOAL := build-with-dependencies

ifeq ($(OS),Windows_NT)
	RM_CMD := cmd.exe /c 'if exist "bin" rmdir /S /Q "bin"'
else
	RM_CMD := rm -rf bin/
endif

clean:
	$(RM_CMD)

build:
	go build csheet.go

build-with-dependencies: dependencies build

dependencies:
	go mod download
	go mod verify

docker-release:
	@docker build --no-cache -f Dockerfile.build -t ninckblokje/csheet-build:latest .
	@docker run --rm -v $${PWD}:/app ninckblokje/csheet-build:latest

install:
	@install -m 0755 csheet /usr/local/bin/csheet
	@install -m 0644 docs/csheet.1 /usr/local/man/man1/csheet.1

release: clean dependencies $(PLATFORMS)

uninstall:
	@rm /usr/local/bin/csheet
	@rm /usr/local/man/man1/csheet.1

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
ext = $(word 3, $(temp))

$(PLATFORMS):
	GOARCH=$(arch) GOOS=$(os) go build -trimpath -ldflags "-X main.csheetVersion=$(GIT_TAG) -X main.csheetRevision=$(GIT_SHORT_REV) -buildid=" -o bin/$(os)_$(arch)/csheet$(ext) csheet.go
	zip -j bin/$(os)_$(arch)/csheet_$(os)_$(arch).zip bin/$(os)_$(arch)/csheet$(ext)

DEB_PKG_ROOT_AMD64 = bin/debian_amd64/packageroot
DEB_PKG_ROOT_ARM64 = bin/debian_arm64/packageroot
DEB_INSTALL_DIR = /usr/bin
DEB_MAN_DIR = /usr/share/man

release-debian-amd64: clean dependencies linux/amd64/
	@mkdir -p $(DEB_PKG_ROOT_AMD64)$(DEB_INSTALL_DIR)
	@mkdir -p $(DEB_PKG_ROOT_AMD64)$(DEB_MAN_DIR)/man1
	@mkdir -p $(DEB_PKG_ROOT_AMD64)/DEBIAN
	@cat pkg/DEBIAN/control | sed s/GIT_TAG/$(GIT_TAG)/g | sed s/any/amd64/g > $(DEB_PKG_ROOT_AMD64)/DEBIAN/control
	@cp bin/linux_amd64/csheet $(DEB_PKG_ROOT_AMD64)$(DEB_INSTALL_DIR)
	@cp docs/csheet.1 $(DEB_PKG_ROOT_AMD64)$(DEB_MAN_DIR)/man1
	@gzip $(DEB_PKG_ROOT_AMD64)$(DEB_MAN_DIR)/man1/csheet.1
	@dpkg-deb -Zxz -b $(DEB_PKG_ROOT_AMD64) bin/debian_amd64/csheet_$(GIT_TAG)_amd64.deb

release-debian-arm64: clean dependencies linux/arm64/
	@mkdir -p $(DEB_PKG_ROOT_ARM64)$(DEB_INSTALL_DIR)
	@mkdir -p $(DEB_PKG_ROOT_ARM64)$(DEB_MAN_DIR)/man1
	@mkdir -p $(DEB_PKG_ROOT_ARM64)/DEBIAN
	@cat pkg/DEBIAN/control | sed s/GIT_TAG/$(GIT_TAG)/g | sed s/any/arm64/g > $(DEB_PKG_ROOT_ARM64)/DEBIAN/control
	@cp bin/linux_arm64/csheet $(DEB_PKG_ROOT_ARM64)$(DEB_INSTALL_DIR)
	@cp docs/csheet.1 $(DEB_PKG_ROOT_ARM64)$(DEB_MAN_DIR)/man1
	@gzip $(DEB_PKG_ROOT_ARM64)$(DEB_MAN_DIR)/man1/csheet.1
	@dpkg-deb -Zxz -b $(DEB_PKG_ROOT_ARM64) bin/debian_arm64/csheet_$(GIT_TAG)_arm64.deb
