package docstring

import (
	"fmt"
	"sync"

	"github.com/kode4food/ale/data"
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

var (
	docStringCache     = map[string][]byte{}
	docStringCacheOnce sync.Once
)

// Get resolves a registered docstring entry by name
func Get(n string) (string, error) {
	ensureDocStringCache()
	res, ok := docStringCache[n]
	if ok {
		return string(res), nil
	}
	return "", fmt.Errorf(ErrDocNotFound, n)
}

// Names returns the registered names of available docstring entries
func Names() []string {
	ensureDocStringCache()
	res := make([]string, 0, len(docStringCache))
	for k := range docStringCache {
		res = append(res, k)
	}
	return res
}

func ensureDocStringCache() {
	docStringCacheOnce.Do(func() {
		for _, filename := range assetNames() {
			doc, _ := getAsset(filename)
			meta := markdown.ParseHeader(string(doc))
			if names := meta.Names; len(names) > 0 {
				for _, name := range names {
					docStringCache[name] = doc
				}
			} else {
				n := filename[0 : len(filename)-len(extension)]
				docStringCache[n] = doc
			}
		}
	})
}
