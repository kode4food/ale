package data

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/read/lex"
	"github.com/kode4food/ale/read/parse"
)

var matcher = lex.MakeMatcher(
	lex.Ignorable,
	lex.Structure,
	lex.Quoting.Error(),
	lex.Values,
	lex.Symbols,
)

// FromString converts the source into a pure data representation. This means
// no quoting, eval, or macro capabilities will be possible
func FromString(src data.String) data.Sequence {
	return parse.FromLexer(Tokens(src))
}

// Tokens creates a new Lexer Sequence of raw Tokens encompassing the subset of
// those required for data representation
func Tokens(src data.String) data.Sequence {
	return lex.Match(src, matcher)
}
