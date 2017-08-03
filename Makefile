build: bin
	go build -ldflags '-s' -o ./bin/lanes ./cmd/lanes

clean:
	rm -rf ./bin/

bin:
	mkdir -p ./bin/
