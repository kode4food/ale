package docstring

import "embed"

var (
	//go:embed *.md
	assets embed.FS

	// getAsset exposes the assets FS ReadFile method
	getAsset = assets.ReadFile
)

// assetNames returns the names of the available docstring files
func assetNames() []string {
	files, _ := assets.ReadDir(".")
	res := make([]string, 0, len(files))
	for _, f := range files {
		res = append(res, f.Name())
	}
	return res
}
