.PHONY: build clean mod-verify

build: bin/adviseu-backend

clean:
	rm -rf ./bin

mod-verify:
	go mod verify

bin/adviseu-backend: *.go go.mod go.sum
	go build -o bin/adviseu-backend