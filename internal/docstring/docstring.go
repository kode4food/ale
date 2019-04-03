package docstring

import (
	"fmt"
	"regexp"
	"strings"

	"gitlab.com/kode4food/ale/internal/assets"
)

const prefix = "internal/docstring/"

var (
	firstLine = regexp.MustCompile(`^names\: ([^\n]+)\n`)
	cache     = map[string]string{}
)

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
			if sm := firstLine.FindStringSubmatch(doc); sm != nil {
				names := strings.Split(sm[1], " ")
				rest := doc[len(sm[0]):]
				for _, name := range names {
					cache[name] = rest
				}
			} else {
				n := filename[len(prefix) : len(filename)-3]
				cache[n] = doc
			}
		}
	}
}
