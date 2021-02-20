all: install

install: build test
	go install github.com/kode4food/ale/cmd/ale

test: build
	golint ./...
	go vet ./...
	go test ./...

build: generate

generate:
	go generate ./...

deps:
	go get -u golang.org/x/tools/cmd/stringer
	go get -u golang.org/x/lint/golint
