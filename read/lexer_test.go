package read_test

import (
	"fmt"
	"testing"

	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/internal/assert"
	. "gitlab.com/kode4food/ale/internal/assert/helpers"
	"gitlab.com/kode4food/ale/read"
)

func makeToken(t read.TokenType, v data.Value) *read.Token {
	return &read.Token{
		Type:  t,
		Value: v,
	}
}

func assertToken(as *assert.Wrapper, like *read.Token, value *read.Token) {
	as.Equal(like.Type, value.Type)
}

func assertTokenSequence(as *assert.Wrapper, s data.Sequence, tokens []*read.Token) {
	var f data.Value
	var r = s
	var ok bool
	for _, l := range tokens {
		f, r, ok = r.Split()
		as.True(ok)
		assertToken(as, l, f.(*read.Token))
	}
	f, r, ok = r.Split()
	as.False(ok)
	as.Nil(f)
	as.False(r.IsSequence())
}

func TestCreateLexer(t *testing.T) {
	as := assert.New(t)
	l := read.Scan("hello")
	as.NotNil(l)
	as.String(`([1 "hello"])`, data.MakeSequenceStr(l))
}

func TestWhitespace(t *testing.T) {
	as := assert.New(t)
	l := read.Scan("   \t ")
	assertTokenSequence(as, l, []*read.Token{})
}

func TestEmptyList(t *testing.T) {
	as := assert.New(t)
	l := read.Scan(" ( \t ) ")
	assertTokenSequence(as, l, []*read.Token{
		makeToken(read.ListStart, S("(")),
		makeToken(read.ListEnd, S(")")),
	})
}

func TestNumbers(t *testing.T) {
	as := assert.New(t)
	l := read.Scan(` 10 12.8 8E+10
				99.598e+10 54e+12 -0xFF
				071 0xf1e9d8c7`)
	assertTokenSequence(as, l, []*read.Token{
		makeToken(read.Number, F(10)),
		makeToken(read.Number, F(12.8)),
		makeToken(read.Number, F(8E+10)),
		makeToken(read.Number, F(99.598e+10)),
		makeToken(read.Number, F(54e+12)),
		makeToken(read.Number, F(-255)),
		makeToken(read.Number, F(57)),
		makeToken(read.Number, F(4058634439)),
	})

	defer as.ExpectPanic(fmt.Sprintf(data.ExpectedInteger, data.String("0xffj-k")))
	read.Scan("0xffj-k").First()
}

func TestStrings(t *testing.T) {
	as := assert.New(t)
	l := read.Scan(` "hello there" "how's \"life\"?"  `)
	assertTokenSequence(as, l, []*read.Token{
		makeToken(read.String, S(`hello there`)),
		makeToken(read.String, S(`how's "life"?`)),
	})
}

func TestMultiLine(t *testing.T) {
	as := assert.New(t)
	l := read.Scan(` "hello there"
  				"how's life?"
				99`)

	assertTokenSequence(as, l, []*read.Token{
		makeToken(read.String, S(`hello there`)),
		makeToken(read.String, S(`how's life?`)),
		makeToken(read.Number, F(99)),
	})
}

func TestComments(t *testing.T) {
	as := assert.New(t)
	l := read.Scan(`"hello" ; (this is commented)`)
	assertTokenSequence(as, l, []*read.Token{
		makeToken(read.String, S(`hello`)),
	})
}
