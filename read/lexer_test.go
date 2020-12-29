package read_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/read"
)

func T(t read.TokenType, v data.Value) *read.Token {
	return &read.Token{
		Type:   t,
		Value:  v,
		Line:   -1,
		Column: -1,
	}
}

func TL(t read.TokenType, v data.Value, line int, col int) *read.Token {
	return &read.Token{
		Type:   t,
		Value:  v,
		Line:   line,
		Column: col,
	}
}

func assertToken(t *testing.T, like *read.Token, value *read.Token) {
	t.Helper()
	as := assert.New(t)
	as.Equal(like.Type, value.Type)

	if like.Line >= 0 {
		as.Equal(like.Line, value.Line)
	}

	if like.Column >= 0 {
		as.Equal(like.Column, value.Column)
	}
}

func assertTokenSequence(t *testing.T, s data.Sequence, tokens []*read.Token) {
	t.Helper()
	as := assert.New(t)
	var f data.Value
	var r = s
	var ok bool
	for _, l := range tokens {
		f, r, ok = r.Split()
		as.True(ok)
		assertToken(t, l, f.(*read.Token))
	}
	f, r, ok = r.Split()
	as.False(ok)
	as.Nil(f)
	as.True(r.IsEmpty())
}

func TestCreateLexer(t *testing.T) {
	as := assert.New(t)
	l := read.Scan("hello")
	as.NotNil(l)
	as.String(`([1 "hello"])`, data.MakeSequenceStr(l))
}

func TestWhitespace(t *testing.T) {
	l := read.Scan("   \t ")
	assertTokenSequence(t, l, []*read.Token{})
}

func TestEmptyList(t *testing.T) {
	l := read.Scan(" ( \t ) ")
	assertTokenSequence(t, l, []*read.Token{
		T(read.ListStart, S("(")),
		T(read.ListEnd, S(")")),
	})
}

func TestNumbers(t *testing.T) {
	l := read.Scan(` 10 12.8 8E+10
				99.598e+10 54e+12 -0xFF
				071 0xf1e9d8c7 2/3`)
	assertTokenSequence(t, l, []*read.Token{
		T(read.Number, F(10)),
		T(read.Number, F(12.8)),
		T(read.Number, F(8e+10)),
		T(read.Number, F(99.598e+10)),
		T(read.Number, F(54e+12)),
		T(read.Number, F(-255)),
		T(read.Number, F(57)),
		T(read.Number, F(4058634439)),
		T(read.Number, R(2, 3)),
	})
}

func TestBadNumbers(t *testing.T) {
	err := fmt.Sprintf(data.ErrExpectedInteger, S("0xffj-k"))
	l := read.Scan("0xffj-k")
	assertTokenSequence(t, l, []*read.Token{
		T(read.Error, S(err)),
	})
}

func TestStrings(t *testing.T) {
	l := read.Scan(` "hello there" "how's \"life\"?"  `)
	assertTokenSequence(t, l, []*read.Token{
		T(read.String, S(`hello there`)),
		T(read.String, S(`how's "life"?`)),
	})
}

func TestMultiLine(t *testing.T) {
	l := read.Scan(` "hello there"
  				"how's life?"
				99`)

	assertTokenSequence(t, l, []*read.Token{
		T(read.String, S(`hello there`)),
		T(read.String, S(`how's life?`)),
		T(read.Number, F(99)),
	})
}

func TestComments(t *testing.T) {
	l := read.Scan(`"hello" ; (this is commented)`)
	assertTokenSequence(t, l, []*read.Token{
		T(read.String, S(`hello`)),
	})
}

func TestIdentifiers(t *testing.T) {
	l := read.Scan(`hello th,@re`)
	assertTokenSequence(t, l, []*read.Token{
		T(read.Identifier, S("hello")),
		T(read.Identifier, S("th,@re")),
	})
}

func TestUnexpectedChars(t *testing.T) {
	err := fmt.Sprintf(read.ErrUnexpectedCharacter, "@")
	l := read.Scan("hello @there")
	assertTokenSequence(t, l, []*read.Token{
		T(read.Identifier, S("hello")),
		T(read.Error, S(err)),
		T(read.Identifier, S("there")),
	})
}

func TestNewLine(t *testing.T) {
	l := read.Scan("1\n2\n 3")

	assertTokenSequence(t, l, []*read.Token{
		TL(read.Number, S("1"), 1, 1),
		TL(read.Number, S("2"), 2, 1),
		TL(read.Number, S("3"), 3, 2),
	})
}

func TestComment(t *testing.T) {
	l := read.Scan("; 1\n2\n 3")

	assertTokenSequence(t, l, []*read.Token{
		TL(read.Number, S("2"), 2, 1),
		TL(read.Number, S("3"), 3, 2),
	})
}
