package docstring

import (
	"embed"
	"io/fs"

	"github.com/kode4food/ale/internal/basics"
)

var (
	//go:embed *.md
	assets embed.FS

	// getAsset exposes the assets FS ReadFile method
	getAsset = assets.ReadFile
)

// assetNames returns the names of the available docstring files
func assetNames() []string {
	files, _ := assets.ReadDir(".")
	return basics.Map(files, func(f fs.DirEntry) string {
		return f.Name()
	})
}
