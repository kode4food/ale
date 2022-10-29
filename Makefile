all: install

install: build test
	go install github.com/kode4food/ale/cmd/ale

test: build
	go vet ./...
	go test ./...

build: generate

generate:
	go generate ./...

deps:
	go install golang.org/x/tools/cmd/stringer@latest
