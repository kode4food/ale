all: install

install: build test
	go install gitlab.com/kode4food/ale/cmd/ale

test: build
	golint ./...
	go vet ./...
	go test ./...

build: generate assets

generate:
	go generate ./...

assets:
	go-snapshot -pkg assets -out internal/assets/assets.go \
		docstring/*.md bootstrap/lisp/*.lisp

deps:
	go get -u gitlab.com/kode4food/go-snapshot
	go get -u golang.org/x/tools/cmd/stringer
	go get -u golang.org/x/lint/golint
