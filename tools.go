//go:build tools
// +build tools

package ale

import (
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/stringer"
)
