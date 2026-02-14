//go:build !windows

package console

import (
	"github.com/chzyer/readline"

	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/lang"
)

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
	Comment = esc + "90m"        // Comment
	Code    = esc + "94m"        // Light Blue
	Result  = esc + "32m"        // Green
	Error   = esc + "31m"        // Red
	NewLine = esc + "90m" + "␤"  // Dark Gray
	Paired  = esc + "7m"         // Invert
)

type painter struct{}

var (
	openers = makeRuneMap(map[string]string{
		lang.ListEnd:   lang.ListStart,
		lang.ObjectEnd: lang.ObjectStart,
		lang.VectorEnd: lang.VectorStart,
	})

	closers = invertMap(openers)
)

// Painter implements the Painter interface for readline
func Painter() readline.Painter {
	return painter{}
}

func (painter) Paint(line []rune, pos int) []rune {
	if len(line) == 0 {
		return line
	}

	l := len(line)
	npos := max(pos, 0)
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
		switch line[i] {
		case o:
			depth--
			if depth == 0 {
				return markMatch(line, i)
			}
		case c:
			depth++
		}
	}
	return line
}

func markCloser(line []rune, pos int, o rune) []rune {
	c := closers[o]
	depth := 0
	for i := pos; i < len(line); i++ {
		switch line[i] {
		case c:
			depth--
			if depth == 0 {
				return markMatch(line, i)
			}
		case o:
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

func makeRuneMap(m map[string]string) map[rune]rune {
	res := make(map[rune]rune, len(m))
	for k, v := range m {
		if len(k) != 1 || len(v) != 1 {
			panic(debug.ProgrammerError("input strings must be single runes"))
		}
		res[rune(k[0])] = rune(v[0])
	}
	return res
}

func invertMap[K, V comparable](m map[K]V) map[V]K {
	res := make(map[V]K, len(m))
	for k, v := range m {
		res[v] = k
	}
	return res
}
