# build lanes for all platforms, compress the binaries, and produce SHA256 checksums
all:
	goreleaser release --snapshot --clean

# run lanes tests
test:
	go test -race -cover `go list ./... | grep -v vendor`

# clean up any build artifacts
clean:
	rm -rf ./bin/ ./dist/
