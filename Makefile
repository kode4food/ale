all: install

install: build test
	go install github.com/kode4food/ale/cmd/ale

test: build
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck ./...
	go test ./...

build: generate

generate:
	go generate ./...
