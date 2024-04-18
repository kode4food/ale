package docstring

import "embed"

//go:embed *.md
var assets embed.FS

// assetNames returns the names of the available docstring files
func assetNames() []string {
	files, _ := assets.ReadDir(".")
	res := make([]string, 0, len(files))
	for _, f := range files {
		res = append(res, f.Name())
	}
	return res
}

// getAsset exposes the assets FS ReadFile method
var getAsset = assets.ReadFile
