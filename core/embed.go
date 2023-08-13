package core

import (
	"embed"
	"slices"
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
	slices.Sort(res)
	return res
}

// Get exposes the assets FS ReadFile method
var Get = assets.ReadFile
