package source

import (
	"embed"
	"io/fs"

	"github.com/kode4food/comb/basics"
)

var (
	//go:embed *.ale
	assets embed.FS

	// Get exposes the assets FS ReadFile method
	Get = assets.ReadFile
)

// Names returns the names of the embedded core scripts
func Names() []string {
	files, _ := assets.ReadDir(".")
	return basics.SortedMap(files, func(f fs.DirEntry) string {
		return f.Name()
	})
}
