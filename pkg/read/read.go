package read

import (
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/internal/lang/parse"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

var matcher = lex.ExhaustiveMatcher(
	lex.Ignorable,
	lex.Structure,
	lex.Quoting,
	lex.Values,
	lex.Preprocessors,
)

// FromString converts the raw source into unexpanded data structures
func FromString(ns env.Namespace, src data.String) data.Sequence {
	return parse.FromLexer(ns, Tokenize, MustTokenize(src))
}

// Tokenize creates a new Lexer Sequence of raw Tokenize encompassing the
// entire set of those supported by the language
func Tokenize(src data.String) (data.Sequence, error) {
	return lex.Match(src, matcher), nil
}

func MustTokenize(src data.String) data.Sequence {
	t, err := Tokenize(src)
	if err != nil {
		panic(err)
	}
	return t
}
