package docstring

import (
	"fmt"
	"strings"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assets"
)

// Error messages
const (
	DocNotFound = "could not find doc: %s"
)

const (
	prefix = "docstring/"
	names  = data.Name("names")
)

var cache = map[string]string{}

// Get resolves documentation using snapshot assets
func Get(n string) string {
	ensureCache()
	res, ok := cache[n]
	if ok {
		return res
	}
	panic(fmt.Errorf(DocNotFound, n))
}

// Exists returns whether a specific docstring exists
func Exists(n string) bool {
	ensureCache()
	_, ok := cache[n]
	return ok
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
				for _, name := range names.(data.Vector) {
					cache[name.String()] = doc
				}
			} else {
				n := filename[len(prefix) : len(filename)-3]
				cache[n] = doc
			}
		}
	}
}
