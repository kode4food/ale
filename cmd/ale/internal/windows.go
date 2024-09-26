//go:build windows

package internal

import "os"

var farewells = getFarewells()

func getFarewells() []string {
	if os.Getenv("ShellId") == "Microsoft.PowerShell" {
		return append(farewellLatin1, farewellUtf8...)
	}
	return farewellLatin1
}
