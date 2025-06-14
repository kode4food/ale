package parse

import (
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
)

type Tokenizer func(data.String) (data.Sequence, error)

// FromLexer returns a Lazy Sequence of scanned data structures
func FromLexer(
	ns env.Namespace, t Tokenizer, lexer data.Sequence,
) data.Sequence {
	var res sequence.LazyResolver
	r := newParser(ns, t, lex.StripWhitespace(lexer))

	res = func() (data.Value, data.Sequence, bool) {
		f, ok, err := r.nextValue()
		if err != nil {
			panic(err)
		}
		if ok {
			return f, sequence.NewLazy(res), true
		}
		return data.Null, data.Null, false
	}

	return sequence.NewLazy(res)
}
