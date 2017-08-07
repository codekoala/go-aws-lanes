REPO ?= github.com/codekoala/go-aws-lanes
TAG ?= $(shell git rev-parse HEAD)-dev
BUILD_DATE := $(shell date +%FT%T%z)

# build lanes for all platforms, compress the binaries, and produce SHA256 checksums
all: linux osx compress checksums

# build lanes for 64-bit Linux
linux:
	GOOS=linux GOARCH=amd64 $(MAKE) build

# build lanes for 64-bit OSX
osx:
	GOOS=darwin GOARCH=amd64 $(MAKE) build

build: bin
	$(eval SUFFIX=_$(GOOS)_$(GOARCH))
	go build -ldflags "-s -X $(REPO)/version.Commit=$(TAG) -X $(REPO)/version.BuildDate=$(BUILD_DATE)" -o ./bin/lanes$(SUFFIX) ./cmd/lanes

# produce SHA256 checksums of all lanes binaries in the ./bin/ directory
checksums: bin
	cd ./bin/; sha256sum lanes* > SHA256SUMS

# compress all lanes binaries
compress: bin get-upx
	@$(UPX) ./bin/lanes*

# run lanes tests
test:
	go test -race -cover `go list ./... | grep -v vendor`

# clean up any build artifacts
clean:
	rm -rf ./bin/ ./upx/

bin:
	@mkdir -p ./bin/

UPX := $(shell which upx)
upx_version := 3.94

# fetch a local copy of upx for binary compression (only if it's not already installed)
get-upx:
ifeq ($(UPX),)
	@mkdir -p ./upx/
	@curl -Ls https://github.com/upx/upx/releases/download/v$(upx_version)/upx-$(upx_version)-amd64_linux.tar.xz | tar Jx -C ./upx/ --strip-components 1
	$(eval UPX := ./upx/upx)
endif
	@$(UPX) --version
