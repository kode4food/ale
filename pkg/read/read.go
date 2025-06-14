package read

import (
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/internal/lang/parse"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/read/internal"
)

var matcher = lex.ExhaustiveMatcher(
	lex.Ignorable,
	lex.Structure,
	lex.Quoting,
	lex.Values,
	lex.Preprocessors,
)

var (
	MustFromString = internal.MakeMustFromString(FromString)
	MustTokenize   = internal.MakeMustTokenizer(Tokenize)
)

// FromString converts the raw source into unexpanded data structures
func FromString(ns env.Namespace, src data.String) (data.Sequence, error) {
	return parse.FromString(ns, Tokenize, src)
}

// Tokenize creates a new Lexer Sequence of raw Tokenize encompassing the
// entire set of those supported by the language
func Tokenize(src data.String) (data.Sequence, error) {
	return lex.Match(src, matcher), nil
}
