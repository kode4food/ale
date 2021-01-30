package data_test

import (
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
)

func TestConstStrings(t *testing.T) {
	as := assert.New(t)

	as.String(data.TrueLiteral, data.True)
	as.String(data.FalseLiteral, data.False)
}

func TestStr(t *testing.T) {
	as := assert.New(t)

	s1 := S("hello")
	as.Number(5, s1.Count())
	as.String("h", s1.First())
	as.String("ello", s1.Rest())

	as.False(s1.IsEmpty())
	as.True(S("").IsEmpty())

	s2 := S("thér\\再e")
	as.Number(7, s2.Count())

	s3 := data.MaybeQuoteString(s2)
	r1 := []rune(s3)
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

	s4 := S("再见!")
	as.Number(3, s4.Count())
	as.String("再", s4.First())
	as.String("见!", s4.Rest())
}

func TestEmptyStr(t *testing.T) {
	as := assert.New(t)

	as.Nil(S("").First())
	as.String("", S("").Rest())

	c, ok := S("").ElementAt(-1)
	as.False(ok)
	as.Nil(c)
}

func TestStringEquality(t *testing.T) {
	as := assert.New(t)

	s1 := S("first string")
	s2 := S("first string")
	s3 := S("not the same")

	as.True(s1.Equal(s1))
	as.True(s1.Equal(s2))
	as.False(s1.Equal(s3))
	as.False(s1.Equal(I(32)))
}
