package docstring

import (
	"fmt"
	"strings"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/assets"
)

const (
	prefix = "docstring/"
	names  = api.Name("names")
)

var cache = map[string]string{}

// Get resolves documentation using snapshot assets
func Get(n string) string {
	ensureCache()
	res, ok := cache[n]
	if ok {
		return res
	}
	panic(fmt.Errorf("could not find doc: %s", n))
}

// Exists returns whether or not a specific docstring exists
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
				for _, name := range names.(api.Vector) {
					cache[name.String()] = doc
				}
			} else {
				n := filename[len(prefix) : len(filename)-3]
				cache[n] = doc
			}
		}
	}
}
