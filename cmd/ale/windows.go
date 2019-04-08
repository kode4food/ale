// +build windows

package main

const (
	cyan     = ""
	dgray    = ""
	green    = ""
	lblue    = ""
	lmagenta = ""
	red      = ""
	yellow   = ""
	bold     = ""
	italic   = ""
	reset    = ""
	clear    = ""
	paired   = ""
	nlMarker = "\\"

	h1     = lmagenta + bold
	h2     = yellow
	code   = lblue
	result = green
)

var farewells = []string{
	"¡Adiós!",
	"Au revoir!",
	"Bye for now!",
	"Ciao!",
	"Hoşçakal!",
	"Tchau!",
	"Tschüss!",
}

// Paint implements the Painter interface
func (r *REPL) Paint(line []rune, pos int) []rune {
	return line
}
