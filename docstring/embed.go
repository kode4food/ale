package docstring

import "embed"

//go:embed *.md
var assets embed.FS

// Names returns the names of the available docstring files
func Names() []string {
	files, _ := assets.ReadDir(".")
	res := make([]string, 0, len(files))
	for _, f := range files {
		res = append(res, f.Name())
	}
	return res
}

// Get exposes the assets FS ReadFile method
var Get = assets.ReadFile
