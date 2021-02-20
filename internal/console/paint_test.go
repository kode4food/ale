// +build !windows

package console_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	"github.com/kode4food/ale/internal/console"
)

func TestREPLPaint(t *testing.T) {
	as := assert.New(t)
	p := console.Painter()

	pair := "\033[7m"
	reset := "\033[0m\033[94m"

	src := "this is (hello)"
	p1 := "this is " + pair + "(" + reset + "hello)"
	p2 := "this is (hello" + pair + ")" + reset
	rs := []rune(src)

	s1 := p.Paint(rs, 0)
	as.String(src, string(s1))
	s2 := p.Paint(rs, -1)
	as.String(src, string(s2))
	s3 := p.Paint(rs, len(src))
	as.String(p1, string(s3))
	s4 := p.Paint(rs, len(src)-1)
	as.String(p1, string(s4))
	s5 := p.Paint(rs, 8)
	as.String(p2, string(s5))
}

func TestREPLNonPaint(t *testing.T) {
	as := assert.New(t)
	p := console.Painter()

	src1 := "(no match"
	src2 := "no match)"

	s1 := p.Paint([]rune(src1), 0)
	as.String(src1, string(s1))
	s2 := p.Paint([]rune(src2), len(src1)-1)
	as.String(src2, string(s2))
	s3 := p.Paint([]rune{}, 0)
	as.String("", string(s3))
}
