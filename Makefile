build: bin
	go build -o ./bin/lanes ./cmd/lanes

clean:
	rm -rf ./bin/

bin:
	mkdir -p ./bin/
