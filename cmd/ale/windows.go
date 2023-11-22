//go:build windows

package main

import "os"

var farewells = getFarewells()

func getFarewells() []string {
	if os.Getenv("ShellId") == "Microsoft.PowerShell" {
		return append(farewell_latin1, farewell_utf8...)
	}
	return farewell_latin1
}
