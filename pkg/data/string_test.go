package data_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/pkg/data"
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
	as.String("h", s1.Car())
	as.String("ello", s1.Cdr())

	as.False(s1.IsEmpty())
	as.True(S("").IsEmpty())

	s2 := S("thér\\再e")
	as.Number(7, s2.Count())

	s3 := data.ToQuotedString(s2)
	r1 := []rune(s3)
	as.Number(10, len(r1))
	as.String(lang.StringQuote, string(r1[0]))

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
	as.String("再", s4.Car())
	as.String("见!", s4.Cdr())

	as.True(types.BasicString.Equal(s4.Type()))
}

func TestEmptyStr(t *testing.T) {
	as := assert.New(t)

	as.Nil(S("").Car())
	as.String("", S("").Cdr())

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

func TestSubstringCall(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`("hello" 0)`, S("hello"))
	as.MustEvalTo(`("hello" 1)`, S("ello"))
	as.MustEvalTo(`("hello" 4)`, S("o"))
	as.MustEvalTo(`("hello" 5)`, S(""))
	as.MustEvalTo(`("hello" 0 1)`, S("h"))
	as.MustEvalTo(`("hello" 0 2)`, S("he"))
	as.MustEvalTo(`("hello" 0 5)`, S("hello"))
	as.MustEvalTo(`("hello" 3 5)`, S("lo"))

	as.PanicWith(`("hello" -1)`, fmt.Errorf(data.ErrInvalidStartIndex, -1))
	as.PanicWith(`("hello" 6)`, fmt.Errorf(data.ErrInvalidStartIndex, 6))
	as.PanicWith(`("hello" 0 6)`, fmt.Errorf(data.ErrInvalidIndexes, 0, 6))
	as.PanicWith(`("hello" 3 2)`, fmt.Errorf(data.ErrInvalidIndexes, 3, 2))
}

func TestReverse(t *testing.T) {
	as := assert.New(t)
	as.MustEvalTo(`(reverse "hello")`, S("olleh"))
	as.MustEvalTo(`(reverse "")`, S(""))
	as.MustEvalTo(`(reverse "X")`, S("X"))
	as.MustEvalTo(`(reverse "😎⚽")`, S("⚽😎"))
	as.MustEvalTo(
		`(reverse "The quick bròwn 狐 jumped over the lazy 犬")`,
		S("犬 yzal eht revo depmuj 狐 nwòrb kciuq ehT"),
	)
}

func TestStringAppend(t *testing.T) {
	as := assert.New(t)

	s1 := S("hello, ")
	s2 := S("there!")
	s3 := s1.Append(s2).(data.String)
	as.String("hello, there!", s3)

	as.MustEvalTo(`(append "😎" "⚽")`, S("😎⚽"))
}
