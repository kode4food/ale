//go:build !windows

package internal

import "slices"

var farewells = slices.Concat(farewellLatin1, farewellUtf8)
