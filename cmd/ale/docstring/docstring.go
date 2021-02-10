package docstring

import (
	"fmt"
	"strings"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/docstring/assets"
)

// Error messages
const (
	ErrDocNotFound = "could not find doc: %s"
)

const (
	prefix    = "docstring/"
	extension = ".md"
	names     = data.Name("names")
)

var cache = map[string]string{}

// Get resolves documentation using snapshot assets
func Get(n string) (string, error) {
	ensureCache()
	res, ok := cache[n]
	if ok {
		return res, nil
	}
	return "", fmt.Errorf(ErrDocNotFound, n)
}

func ensureCache() {
	if len(cache) > 0 {
		return
	}
	for _, filename := range assets.AssetNames() {
		if strings.HasPrefix(filename, prefix) {
			doc := string(assets.MustGet(filename))
			meta, _ := ParseMarkdown(doc)
			if names, ok := meta.Get(names); ok {
				for _, name := range names.(data.Vector).Values() {
					cache[name.String()] = doc
				}
			} else {
				n := filename[len(prefix) : len(filename)-len(extension)]
				cache[n] = doc
			}
		}
	}
}
