// +build windows

package console

import "github.com/chzyer/readline"

// Terminal codes
const (
	Bold   = ""
	Italic = ""
	Clear  = ""
	Reset  = ""

	Header1 = ""
	Header2 = ""
	Domain  = ""
	Code    = ""
	Result  = ""
	Error   = ""
	NewLine = ""
	Paired  = ""
)

type painter struct{}

// Painter implements the Painter interface for readline
func Painter() readline.Painter {
	return new(painter)
}

func (*painter) Paint(line []rune, pos int) []rune {
	return line
}
