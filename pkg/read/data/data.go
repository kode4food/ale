package data

import (
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/internal/lang/parse"
	"github.com/kode4food/ale/pkg/data"
)

var matcher = lex.ExhaustiveMatcher(
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

// Tokens create a new Lexer Sequence of raw Tokens encompassing the subset of
// those required for data representation
func Tokens(src data.String) data.Sequence {
	return lex.Match(src, matcher)
}
