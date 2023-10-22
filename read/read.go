package read

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/read/lex"
	"github.com/kode4food/ale/read/parse"
)

var langMatcher = lex.MakeMatcher(
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

// Tokens creates a new Lexer Sequence of raw Tokens encompassing the entire
// set of those supported by the language
func Tokens(src data.String) data.Sequence {
	return lex.Match(src, langMatcher)
}
