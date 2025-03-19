package lex_test

import (
	"fmt"
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/read"
	rdata "github.com/kode4food/ale/pkg/read/data"
	"github.com/kode4food/ale/pkg/read/lex"
)

func T(t lex.TokenType, v data.Value) *lex.Token {
	return t.FromValue(``, v)
}

func assertToken(t *testing.T, like *lex.Token, value *lex.Token) {
	t.Helper()
	as := assert.New(t)
	as.Equal(like.Type(), value.Type())
}

func assertTokenSequence(t *testing.T, s data.Sequence, tokens []*lex.Token) {
	t.Helper()
	as := assert.New(t)
	var f data.Value
	var r = s
	var ok bool
	for _, l := range tokens {
		f, r, ok = r.Split()
		as.True(ok)
		assertToken(t, l, f.(*lex.Token))
	}
	f, r, ok = r.Split()
	as.False(ok)
	as.Nil(f)
	as.True(r.IsEmpty())
}

func TestCreateLexer(t *testing.T) {
	as := assert.New(t)
	l := read.Tokens("hello")
	as.NotNil(l)
	f, r, ok := l.Split()
	as.True(ok)
	as.True(r.IsEmpty())
	tk, ok := f.(*lex.Token)
	as.NotNil(tk)
	as.True(ok)
	as.String("hello", tk.Value())
	as.Equal(lex.Identifier, tk.Type())
}

func TestKeyword(t *testing.T) {
	l := lex.StripWhitespace(read.Tokens("  :hello  "))
	assertTokenSequence(t, l, []*lex.Token{
		T(lex.Keyword, K("hello")),
	})
}

func TestWhitespace(t *testing.T) {
	l := lex.StripWhitespace(read.Tokens("   \t "))
	assertTokenSequence(t, l, []*lex.Token{})
}

func TestEmptyList(t *testing.T) {
	l := lex.StripWhitespace(read.Tokens(" ( \t ) "))
	assertTokenSequence(t, l, []*lex.Token{
		T(lex.ListStart, S("(")),
		T(lex.ListEnd, S(")")),
	})
}

func TestNumbers(t *testing.T) {
	l := lex.StripWhitespace(
		read.Tokens(
			` 10 12.8 8E+10
				99.598e+10 54e+12 -0xFF
				071 0xf1e9d8c7 2/3`,
		),
	)
	assertTokenSequence(t, l, []*lex.Token{
		T(lex.Number, F(10)),
		T(lex.Number, F(12.8)),
		T(lex.Number, F(8e+10)),
		T(lex.Number, F(99.598e+10)),
		T(lex.Number, F(54e+12)),
		T(lex.Number, F(-255)),
		T(lex.Number, F(57)),
		T(lex.Number, F(4058634439)),
		T(lex.Number, R(2, 3)),
	})
}

func TestBadNumbers(t *testing.T) {
	err := fmt.Sprintf(data.ErrExpectedInteger, S("0xffj-k"))
	l := read.Tokens("0xffj-k")
	assertTokenSequence(t, l, []*lex.Token{
		T(lex.Error, S(err)),
	})
}

func TestStrings(t *testing.T) {
	l := lex.StripWhitespace(
		read.Tokens(` "hello there" "how's \"life\"?"  `),
	)
	assertTokenSequence(t, l, []*lex.Token{
		T(lex.String, S(`hello there`)),
		T(lex.String, S(`how's "life"?`)),
	})
}

func TestMultiLine(t *testing.T) {
	as := assert.New(t)

	l := lex.StripWhitespace(
		read.Tokens("   \"hello there\"\n\"how's life?\"\n\n  99"),
	)
	assertTokenSequence(t, l, []*lex.Token{
		T(lex.String, S(`hello there`)),
		T(lex.String, S(`how's life?`)),
		T(lex.Number, F(99)),
	})

	v := sequence.ToVector(l)

	as.String(`"hello there"`, v[0].(*lex.Token).Input())
	as.Equal(0, v[0].(*lex.Token).Line())
	as.Equal(3, v[0].(*lex.Token).Column())

	as.String(`"how's life?"`, v[1].(*lex.Token).Input())
	as.Equal(1, v[1].(*lex.Token).Line())
	as.Equal(0, v[1].(*lex.Token).Column())

	as.String("99", v[2].(*lex.Token).Input())
	as.Equal(3, v[2].(*lex.Token).Line())
	as.Equal(2, v[2].(*lex.Token).Column())
}

func TestComments(t *testing.T) {
	l1 := lex.StripWhitespace(
		read.Tokens(`
			#| this is a comment |#
			"hello"
			#| nested
				#| comment is |#
               here
            |#  ; with an eol comment
		`),
	)
	assertTokenSequence(t, l1, []*lex.Token{
		T(lex.String, S(`hello`)),
	})

	l2 := lex.StripWhitespace(rdata.Tokens("hello |# there"))
	assertTokenSequence(t, l2, []*lex.Token{
		T(lex.Identifier, S("hello")),
		T(lex.Error, S(lex.ErrUnmatchedComment)),
		T(lex.Identifier, S("there")),
	})
}

func TestIdentifiers(t *testing.T) {
	l := lex.StripWhitespace(
		read.Tokens(`hello th,@re ale/test / ale// /ale/er/ror`),
	)
	assertTokenSequence(t, l, []*lex.Token{
		T(lex.Identifier, S("hello")),
		T(lex.Identifier, S("th,@re")),
		T(lex.Identifier, S("ale/test")),
		T(lex.Identifier, S("/")),
		T(lex.Identifier, S("ale//")),
		T(lex.Identifier, S("/ale/er/ror")),
	})
}

func TestUnexpectedChars(t *testing.T) {
	err := fmt.Sprintf(lex.ErrUnexpectedCharacters, "@")
	l1 := lex.StripWhitespace(read.Tokens("hello @there"))
	assertTokenSequence(t, l1, []*lex.Token{
		T(lex.Identifier, S("hello")),
		T(lex.Error, S(err)),
		T(lex.Identifier, S("there")),
	})

	err = fmt.Sprintf(lex.ErrUnexpectedCharacters, "'")
	l2 := lex.StripWhitespace(rdata.Tokens("hello 'there"))
	assertTokenSequence(t, l2, []*lex.Token{
		T(lex.Identifier, S("hello")),
		T(lex.Error, S(err)),
		T(lex.Identifier, S("there")),
	})
}

func TestUnterminatedString(t *testing.T) {
	l := lex.StripWhitespace(read.Tokens(`"unterminated `))
	assertTokenSequence(t, l, []*lex.Token{
		T(lex.Error, S(lex.ErrStringNotTerminated)),
	})
}

func TestTokenEquality(t *testing.T) {
	as := assert.New(t)

	t1 := T(lex.Identifier, S("hello"))
	t2 := T(lex.Identifier, S("hello")) // Content same
	t3 := T(lex.String, S("hello"))     // Type different
	t4 := T(lex.Number, I(37))
	t5 := T(lex.Number, I(38)) // Value different

	as.True(t1.Equal(t1))
	as.True(t1.Equal(t2))
	as.False(t1.Equal(t4))
	as.False(t1.Equal(t3))
	as.False(t1.Equal(t5))
	as.False(t4.Equal(t5))
	as.False(t1.Equal(I(37)))
}
