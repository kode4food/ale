package main

import (
	"fmt"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/docstring"
	"github.com/kode4food/ale/internal/markdown"
)

// Error messages
const (
	ErrDocNotFound = "could not find doc: %s"
)

const (
	extension = ".md"
	names     = data.Name("names")
)

var docStringCache = map[string][]byte{}

// GetDocString resolves documentation using snapshot assets
func GetDocString(n string) (string, error) {
	ensureDocStringCache()
	res, ok := docStringCache[n]
	if ok {
		return string(res), nil
	}
	return "", fmt.Errorf(ErrDocNotFound, n)
}

func ensureDocStringCache() {
	if len(docStringCache) > 0 {
		return
	}
	for _, filename := range docstring.Names() {
		doc, _ := docstring.Get(filename)
		meta := markdown.ParseHeader(string(doc))
		if names, ok := meta.Get(names); ok {
			for _, name := range names.(data.Vector).Values() {
				docStringCache[name.String()] = doc
			}
		} else {
			n := filename[0 : len(filename)-len(extension)]
			docStringCache[n] = doc
		}
	}
}
