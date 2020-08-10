all: install

install: build test
	go install github.com/kode4food/ale/cmd/ale

test: build
	golint ./...
	go vet ./...
	go test ./...

build: generate assets

generate:
	go generate ./...

assets:
	go-snapshot -pkg assets -out core/assets/assets.go core/*.scm
	go-snapshot -pkg assets -out docstring/assets/assets.go docstring/*.md

deps:
	go get -u github.com/kode4food/go-snapshot
	go get -u golang.org/x/tools/cmd/stringer
	go get -u golang.org/x/lint/golint
