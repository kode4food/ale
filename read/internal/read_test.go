package internal_test

import (
	"testing"

	"github.com/kode4food/ale/internal/assert"
	. "github.com/kode4food/ale/internal/assert/helpers"
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/read/internal"
)

var (
	Tokenize = internal.MakeTokenizer(lex.ExhaustiveMatcher(
		lex.Ignorable,
		lex.Structure,
		lex.Quoting,
		lex.Values,
		lex.Symbols,
	))

	MustTokenize   = internal.MakeMustTokenizer(Tokenize)
	FromString     = internal.MakeFromString(Tokenize)
	MustFromString = internal.MakeMustFromString(FromString)
)

func TestFromString(t *testing.T) {
	as := assert.New(t)
	ns := assert.GetTestNamespace()
	tr := MustFromString(ns, "99")
	as.NotNil(tr)
	as.Equal(I(99), tr.Car())
}

func TestTokenize(t *testing.T) {
	as := assert.New(t)
	tr := MustTokenize("99")
	as.NotNil(tr)
	as.Equal(lex.Number.FromValue("99", I(99)), tr.Car())
}
