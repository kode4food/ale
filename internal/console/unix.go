// +build !windows

package console

import "github.com/chzyer/readline"

// Terminal codes
const (
	esc = "\033["

	Bold   = esc + "1m"
	Italic = esc + "3m"
	Clear  = esc + "2J" + esc + "f"
	Reset  = esc + "0m"

	Header1 = esc + "35m" + Bold // Light Magenta
	Header2 = esc + "33m"        // Yellow
	Domain  = esc + "36m"        // Cyan
	Code    = esc + "94m"        // Light Blue
	Result  = esc + "32m"        // Green
	Error   = esc + "31m"        // Red
	NewLine = esc + "90m" + "␤"  // Dark Gray
	Paired  = esc + "7m"         // Invert
)

type painter struct{}

var (
	openers = map[rune]rune{')': '(', ']': '[', '}': '{'}
	closers = map[rune]rune{'(': ')', '[': ']', '{': '}'}
)

// Painter implements the Painter interface for readline
func Painter() readline.Painter {
	return new(painter)
}

func (*painter) Paint(line []rune, pos int) []rune {
	if line == nil || len(line) == 0 {
		return line
	}

	l := len(line)
	npos := pos
	if npos < 0 {
		npos = 0
	}
	if npos >= l {
		npos = l - 1
	}
	k := line[npos]
	if _, ok := openers[k]; ok {
		return markOpener(line, npos, k)
	} else if _, ok := closers[k]; ok {
		return markCloser(line, npos, k)
	}
	return line
}

func markOpener(line []rune, pos int, c rune) []rune {
	o := openers[c]
	depth := 0
	for i := pos; i >= 0; i-- {
		if line[i] == o {
			depth--
			if depth == 0 {
				return markMatch(line, i)
			}
		} else if line[i] == c {
			depth++
		}
	}
	return line
}

func markCloser(line []rune, pos int, o rune) []rune {
	c := closers[o]
	depth := 0
	for i := pos; i < len(line); i++ {
		if line[i] == c {
			depth--
			if depth == 0 {
				return markMatch(line, i)
			}
		} else if line[i] == o {
			depth++
		}
	}
	return line
}

func markMatch(line []rune, pos int) []rune {
	var m []rune
	if pos > 0 {
		m = append(m, line[:pos]...)
	}
	m = append(m, []rune(Paired)...)
	m = append(m, line[pos])
	m = append(m, []rune(Reset+Code)...)
	if pos < len(line)-1 {
		m = append(m, line[pos+1:]...)
	}
	return m
}
