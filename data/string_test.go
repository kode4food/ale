package data_test

import (
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestConstStrings(t *testing.T) {
	as := assert.New(t)

	as.String("true", data.True)
	as.String("false", data.False)
}

func TestStr(t *testing.T) {
	as := assert.New(t)

	s1 := S("hello")
	as.Number(5, data.Size(s1))
	as.String("h", s1.First())
	as.String("ello", s1.Rest())

	as.False(s1.IsEmpty())
	as.True(S("").IsEmpty())

	s2 := s1.Prepend(S("s"))
	as.Number(6, data.Size(s2))
	as.String("shello", s2)

	s3 := s1.Append(S("z"), S("y"))
	as.Number(7, data.Size(s3))
	as.String("hellozy", s3)

	l1 := s1.Prepend(F(99))
	as.Number(6, data.Size(l1))
	as.String(`(99 "h" "e" "l" "l" "o")`, l1)

	v1 := s1.Append(F(99))
	as.Number(6, data.Size(v1))
	as.String(`["h" "e" "l" "l" "o" 99]`, v1)

	s4 := S("thér\\再e")
	as.Number(7, data.Size(s4))

	s5 := data.MaybeQuoteString(s4)
	r1 := []rune(s5)
	as.Number(10, len(r1))
	as.String(`"`, string(r1[0]))

	c, ok := s1.ElementAt(1)
	as.True(ok)
	as.String("e", c)

	c, ok = s1.ElementAt(5)
	as.False(ok)
	as.Nil(c)

	c, ok = s1.ElementAt(6)
	as.False(ok)
	as.Nil(c)

	s6 := S("再见!")
	as.Number(3, data.Size(s6))
	as.String("再", s6.First())
	as.String("见!", s6.Rest())
}

func TestEmptyStr(t *testing.T) {
	as := assert.New(t)

	as.Nil(S("").First())
	as.String("", S("").Rest())

	c, ok := S("").ElementAt(-1)
	as.False(ok)
	as.Nil(c)
}
