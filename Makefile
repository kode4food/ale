DIST_DIR ?= ./dist
EXE = $(DIST_DIR)/ale
GO ?= go

.PHONY: all install build format check test pre-commit generate clean

all: build

install: test
	$(GO) install github.com/kode4food/ale/cmd/ale

build: test
	@mkdir -p $(DIST_DIR)
	@rm -f $(EXE)
	$(GO) build -o $(EXE) ./cmd/ale

format: generate
	$(GO) run golang.org/x/tools/cmd/goimports@latest -w .
	$(GO) fix ./...

check:
	$(GO) vet ./...
	$(GO) run honnef.co/go/tools/cmd/staticcheck ./...

test: generate check
	$(GO) test ./...

pre-commit: format test

generate:
	$(GO) generate ./...

clean:
	@rm -f $(EXE)
	@rmdir $(DIST_DIR)
