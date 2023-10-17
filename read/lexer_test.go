package read_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/read"
)

func T(t read.TokenType, v data.Value) *read.Token {
	return read.MakeToken(t, v)
}

func assertToken(t *testing.T, like *read.Token, value *read.Token) {
	t.Helper()
	as := assert.New(t)
	as.Equal(like.Type(), value.Type())
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
	as.String(`([Identifier "hello"])`, data.MakeSequenceStr(l))
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
	as := assert.New(t)

	l := read.Scan("   \"hello there\"\n\"how's life?\"\n\n  99")
	assertTokenSequence(t, l, []*read.Token{
		T(read.String, S(`hello there`)),
		T(read.String, S(`how's life?`)),
		T(read.Number, F(99)),
	})

	v := sequence.ToValues(l)
	as.Equal(0, v[0].(*read.Token).Line())
	as.Equal(3, v[0].(*read.Token).Column())
	as.Equal(1, v[1].(*read.Token).Line())
	as.Equal(0, v[1].(*read.Token).Column())
	as.Equal(3, v[2].(*read.Token).Line())
	as.Equal(2, v[2].(*read.Token).Column())
}

func TestComments(t *testing.T) {
	l := read.Scan(`"hello" ; (this is commented)`)
	assertTokenSequence(t, l, []*read.Token{
		T(read.String, S(`hello`)),
	})
}

func TestIdentifiers(t *testing.T) {
	l := read.Scan(`hello th,@re ale/test / ale// /ale/er/ror`)
	assertTokenSequence(t, l, []*read.Token{
		T(read.Identifier, S("hello")),
		T(read.Identifier, S("th,@re")),
		T(read.Identifier, S("ale/test")),
		T(read.Identifier, S("/")),
		T(read.Identifier, S("ale//")),
		T(read.Identifier, S("/ale/er/ror")),
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

func TestTokenEquality(t *testing.T) {
	as := assert.New(t)

	t1 := T(read.Identifier, S("hello"))
	t2 := T(read.Identifier, S("hello")) // Content same
	t3 := T(read.String, S("hello"))     // Type different
	t4 := T(read.Number, I(37))
	t5 := T(read.Number, I(38)) // Value different

	as.True(t1.Equal(t1))
	as.True(t1.Equal(t2))
	as.False(t1.Equal(t4))
	as.False(t1.Equal(t3))
	as.False(t1.Equal(t5))
	as.False(t4.Equal(t5))
	as.False(t1.Equal(I(37)))
}

func TestStrangeToken(t *testing.T) {
	as := assert.New(t)

	tkn := read.MakeToken(read.TokenType(99), S("hello"))
	as.String(`[TokenType(99) "hello"]`, tkn)
}
