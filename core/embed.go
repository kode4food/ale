package core

import (
	"embed"
	"io/fs"

	"github.com/kode4food/comb/basics"
)

//go:embed *.ale
var assets embed.FS

// Names returns the names of the embedded core scripts
func Names() []string {
	files, _ := assets.ReadDir(".")
	return basics.SortedMap(files, func(f fs.DirEntry) string {
		return f.Name()
	})
}

// Get exposes the assets FS ReadFile method
var Get = assets.ReadFile
