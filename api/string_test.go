package api_test

import (
	"testing"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
)

func TestConstStrings(t *testing.T) {
	as := assert.New(t)

	as.String("true", api.True)
	as.String("false", api.False)
}

func TestStr(t *testing.T) {
	as := assert.New(t)

	s1 := S("hello")
	as.Integer(5, api.Count(s1))
	as.String("h", s1.First())
	as.String("ello", s1.Rest())

	as.True(s1.IsSequence())
	as.False(S("").IsSequence())

	s2 := s1.Prepend(S("s"))
	as.Integer(6, api.Count(s2))
	as.String("shello", s2)

	s3 := s1.Conjoin(S("z"))
	as.Integer(6, api.Count(s3))
	as.String("helloz", s3)

	l1 := s1.Prepend(F(99))
	as.Integer(6, api.Count(l1))
	as.String(`(99 "h" "e" "l" "l" "o")`, l1)

	v1 := s1.Conjoin(F(99))
	as.Integer(6, api.Count(v1))
	as.String(`["h" "e" "l" "l" "o" 99]`, v1)

	s4 := S("thér\\再e")
	as.Integer(7, api.Count(s4))

	s5 := api.MaybeQuoteString(s4)
	r1 := []rune(s5)
	as.Integer(10, len(r1))
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
	as.Integer(3, api.Count(s6))
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
