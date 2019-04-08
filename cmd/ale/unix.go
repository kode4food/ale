// +build !windows

package main

const (
	esc      = "\033["
	cyan     = esc + "36m"
	dgray    = esc + "90m"
	green    = esc + "32m"
	lblue    = esc + "94m"
	lmagenta = esc + "35m"
	red      = esc + "31m"
	yellow   = esc + "33m"
	bold     = esc + "1m"
	italic   = esc + "3m"
	reset    = esc + "0m"
	clear    = esc + "2J" + esc + "f"
	paired   = esc + "7m"
	nlMarker = "␤"

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
	"Αντίο!",
	"До свидания!",
	"अलविदा!",
	"안녕!",
	"じゃあね",
	"再见!",
}

// Paint implements the Painter interface
func (r *REPL) Paint(line []rune, pos int) []rune {
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
	m = append(m, []rune(paired)...)
	m = append(m, line[pos])
	m = append(m, []rune(reset+code)...)
	if pos < len(line)-1 {
		m = append(m, line[pos+1:]...)
	}
	return m
}
