REPO ?= github.com/codekoala/go-aws-lanes
TAG ?= dev
BUILD_DATE := $(shell date +%FT%T%z)

all: linux osx checksums

linux:
	GOOS=linux GOARCH=amd64 $(MAKE) build

osx:
	GOOS=darwin GOARCH=amd64 $(MAKE) build

build: bin
	$(eval SUFFIX=_$(GOOS)_$(GOARCH))
	go build -ldflags "-s -X $(REPO)/version.Tag=$(TAG) -X $(REPO)/version.BuildDate=$(BUILD_DATE)" -o ./bin/lanes$(SUFFIX) ./cmd/lanes

checksums:
	cd ./bin/; sha256sum lanes* > SHA256SUMS

compress:
	@upx ./bin/lanes

test:
	go test -race -cover `go list ./... | grep -v vendor`

clean:
	rm -rf ./bin/

bin:
	mkdir -p ./bin/
