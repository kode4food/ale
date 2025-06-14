package read

import (
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/pkg/read/internal"
)

var (
	Tokenize = internal.MakeTokenizer(lex.ExhaustiveMatcher(
		lex.Ignorable,
		lex.Structure,
		lex.Quoting,
		lex.Values,
		lex.Preprocessors,
	))

	MustTokenize   = internal.MakeMustTokenizer(Tokenize)
	FromString     = internal.MakeFromString(Tokenize)
	MustFromString = internal.MakeMustFromString(FromString)
)
