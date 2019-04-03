// +build windows

package main

const (
	red      = ""
	green    = ""
	yellow   = ""
	cyan     = ""
	dgray    = ""
	lyellow  = ""
	lblue    = ""
	bold     = ""
	italic   = ""
	reset    = ""
	clear    = ""
	paired   = ""
	nlMarker = "\\"

	h1     = lyellow
	h2     = yellow
	code   = lblue
	result = green
)

var farewells = []string{
	"¡Adiós!",
	"Au revoir!",
	"Bye for now!",
	"Ciao!",
	"Tchau!",
	"Tschüss!",
	"Hoşçakal!",
}

// Paint implements the Painter interface
func (r *REPL) Paint(line []rune, pos int) []rune {
	return line
}
