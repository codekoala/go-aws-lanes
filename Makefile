REPO ?= github.com/codekoala/go-aws-lanes
TAG ?= $(shell git rev-parse HEAD)-dev
BUILD_DATE := $(shell date +%FT%T%z)

all: linux osx checksums

linux:
	GOOS=linux GOARCH=amd64 $(MAKE) build

osx:
	GOOS=darwin GOARCH=amd64 $(MAKE) build

build: bin
	$(eval SUFFIX=_$(GOOS)_$(GOARCH))
	go build -ldflags "-s -X $(REPO)/version.Commit=$(TAG) -X $(REPO)/version.BuildDate=$(BUILD_DATE)" -o ./bin/lanes$(SUFFIX) ./cmd/lanes

checksums:
	cd ./bin/; sha256sum lanes* > SHA256SUMS

compress: get-upx
	@$(UPX) ./bin/lanes*

test:
	go test -race -cover `go list ./... | grep -v vendor`

clean:
	rm -rf ./bin/ ./upx/

bin:
	mkdir -p ./bin/

UPX := $(shell which upx)
upx_version := 3.94
get-upx:
ifeq ($(UPX),)
	@mkdir -p ./upx/
	@curl -Ls https://github.com/upx/upx/releases/download/v$(upx_version)/upx-$(upx_version)-amd64_linux.tar.xz | tar Jx -C ./upx/ --strip-components 1
	$(eval UPX := ./upx/upx)
endif
	@$(UPX) --version
