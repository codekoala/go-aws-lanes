build: bin
	go build -race -ldflags '-s' -o ./bin/lanes ./cmd/lanes

compress:
	@upx ./bin/lanes

test:
	go test -race -cover `go list ./... | grep -v vendor`

clean:
	rm -rf ./bin/

bin:
	mkdir -p ./bin/
