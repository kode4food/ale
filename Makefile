all: install

install: build test
	go install github.com/kode4food/ale/cmd/ale

test: build
	go vet ./...
	go run golang.org/x/lint/golint ./...
	go test ./...

build: generate

generate:
	go generate ./...
