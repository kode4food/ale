DIST_DIR ?= ./dist
EXE = $(DIST_DIR)/ale
GO ?= go

.PHONY: all install build test generate

all: build

install: test
	$(GO) install github.com/kode4food/ale/cmd/ale

build: test
	@mkdir -p $(DIST_DIR)
	@rm -f $(EXE)
	$(GO) build -o $(EXE) ./cmd/ale

test: generate
	$(GO) test ./...
	$(GO) vet ./...
	$(GO) run honnef.co/go/tools/cmd/staticcheck ./...

generate:
	$(GO) generate ./...

clean:
	@rm -f $(EXE)
	@rmdir $(DIST_DIR)
