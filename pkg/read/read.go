package read

import (
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/internal/lang/parse"
	"github.com/kode4food/ale/pkg/data"
)

var matcher = lex.ExhaustiveMatcher(
	lex.Ignorable,
	lex.Structure,
	lex.Quoting,
	lex.Values,
	lex.Symbols,
)

// FromString converts the raw source into unexpanded data structures
func FromString(src data.String) data.Sequence {
	return parse.FromLexer(Tokens(src))
}

// Tokens create a new Lexer Sequence of raw Tokens encompassing the entire
// set of those supported by the language
func Tokens(src data.String) data.Sequence {
	return lex.Match(src, matcher)
}
