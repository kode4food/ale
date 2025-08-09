//go:build tools
// +build tools

package ale

import (
	_ "github.com/kode4food/gen-maxkind"
	_ "golang.org/x/tools/cmd/stringer"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
