package data

import (
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/read/internal"
)

var (
	Tokenize = internal.MakeTokenizer(lex.ExhaustiveMatcher(
		lex.Ignorable,
		lex.Structure,
		lex.Quoting.Error(),
		lex.Values,
		lex.Symbols,
	))

	MustTokenize   = internal.MakeMustTokenizer(Tokenize)
	FromString     = internal.MakeFromString(Tokenize)
	MustFromString = internal.MakeMustFromString(FromString)
)
