// +build !windows

package main_test

import (
	"testing"

	main "gitlab.com/kode4food/ale/cmd/ale"
	"gitlab.com/kode4food/ale/internal/assert"
)

func TestREPLPaint(t *testing.T) {
	as := assert.New(t)
	r := main.NewREPL()

	pair := "\033[7m"
	reset := "\033[0m\033[94m"

	src := "this is (hello)"
	p1 := "this is " + pair + "(" + reset + "hello)"
	p2 := "this is (hello" + pair + ")" + reset
	rs := []rune(src)

	s1 := r.Paint(rs, 0)
	as.String(src, string(s1))
	s2 := r.Paint(rs, -1)
	as.String(src, string(s2))
	s3 := r.Paint(rs, len(src))
	as.String(p1, string(s3))
	s4 := r.Paint(rs, len(src)-1)
	as.String(p1, string(s4))
	s5 := r.Paint(rs, 8)
	as.String(p2, string(s5))
}

func TestREPLNonPaint(t *testing.T) {
	as := assert.New(t)
	r := main.NewREPL()

	src1 := "(no match"
	src2 := "no match)"

	s1 := r.Paint([]rune(src1), 0)
	as.String(src1, string(s1))
	s2 := r.Paint([]rune(src2), len(src1)-1)
	as.String(src2, string(s2))
	s3 := r.Paint([]rune{}, 0)
	as.String("", string(s3))
}
