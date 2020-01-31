package read_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/read"
)

func makeToken(t read.TokenType, v data.Value) *read.Token {
	return &read.Token{
		Type:  t,
		Value: v,
	}
}

func assertToken(t *testing.T, like *read.Token, value *read.Token) {
	t.Helper()
	as := assert.New(t)
	as.Equal(like.Type, value.Type)
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
		makeToken(read.ListStart, S("(")),
		makeToken(read.ListEnd, S(")")),
	})
}

func TestNumbers(t *testing.T) {
	l := read.Scan(` 10 12.8 8E+10
				99.598e+10 54e+12 -0xFF
				071 0xf1e9d8c7 2/3`)
	assertTokenSequence(t, l, []*read.Token{
		makeToken(read.Number, F(10)),
		makeToken(read.Number, F(12.8)),
		makeToken(read.Number, F(8e+10)),
		makeToken(read.Number, F(99.598e+10)),
		makeToken(read.Number, F(54e+12)),
		makeToken(read.Number, F(-255)),
		makeToken(read.Number, F(57)),
		makeToken(read.Number, F(4058634439)),
		makeToken(read.Number, R(2, 3)),
	})
}

func TestBadNumbers(t *testing.T) {
	err := fmt.Sprintf(data.ErrExpectedInteger, S("0xffj-k"))
	l := read.Scan("0xffj-k")
	assertTokenSequence(t, l, []*read.Token{
		makeToken(read.Error, S(err)),
	})
}

func TestStrings(t *testing.T) {
	l := read.Scan(` "hello there" "how's \"life\"?"  `)
	assertTokenSequence(t, l, []*read.Token{
		makeToken(read.String, S(`hello there`)),
		makeToken(read.String, S(`how's "life"?`)),
	})
}

func TestMultiLine(t *testing.T) {
	l := read.Scan(` "hello there"
  				"how's life?"
				99`)

	assertTokenSequence(t, l, []*read.Token{
		makeToken(read.String, S(`hello there`)),
		makeToken(read.String, S(`how's life?`)),
		makeToken(read.Number, F(99)),
	})
}

func TestComments(t *testing.T) {
	l := read.Scan(`"hello" ; (this is commented)`)
	assertTokenSequence(t, l, []*read.Token{
		makeToken(read.String, S(`hello`)),
	})
}

func TestIdentifiers(t *testing.T) {
	l := read.Scan(`hello th,@re`)
	assertTokenSequence(t, l, []*read.Token{
		makeToken(read.Identifier, S("hello")),
		makeToken(read.Identifier, S("th")),
		makeToken(read.SpliceMarker, S(",@")),
		makeToken(read.Identifier, S("re")),
	})
}

func TestUnexpectedChars(t *testing.T) {
	err := fmt.Sprintf(read.ErrUnexpectedCharacter, "@")
	l := read.Scan("th@re")
	assertTokenSequence(t, l, []*read.Token{
		makeToken(read.Identifier, S("th")),
		makeToken(read.Error, S(err)),
		makeToken(read.Identifier, S("re")),
	})
}
