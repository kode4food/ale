package core

import (
	"embed"
	"sort"
)

//go:embed *.ale
var assets embed.FS

// Names returns the names of the embedded core scripts
func Names() []string {
	files, _ := assets.ReadDir(".")
	res := make([]string, 0, len(files))
	for _, f := range files {
		res = append(res, f.Name())
	}
	sort.Strings(res)
	return res
}

// Get exposes the assets FS ReadFile method
var Get = assets.ReadFile
