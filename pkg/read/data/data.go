package data

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
	lex.Quoting.Error(),
	lex.Values,
	lex.Symbols,
)

var (
	MustFromString = internal.MakeMustFromString(FromString)
	MustTokenize   = internal.MakeMustTokenize(Tokenize)
)

// FromString converts the source into a pure data representation. This means
// no quoting, eval, or macro capabilities will be possible
func FromString(ns env.Namespace, src data.String) (data.Sequence, error) {
	return parse.FromString(ns, Tokenize, src)
}

// Tokenize creates a new Lexer Sequence of raw Tokenize encompassing the
// subset of those required for data representation
func Tokenize(src data.String) (data.Sequence, error) {
	return lex.Match(src, matcher), nil
}
