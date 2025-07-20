package parse

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/internal/sequence"
)

type Tokenizer func(data.String) (data.Sequence, error)

// FromString returns a Lazy Sequence of scanned data structures
func FromString(
	ns env.Namespace, tokenize Tokenizer, str data.String,
) (data.Sequence, error) {
	lexer, err := tokenize(str)
	if err != nil {
		return nil, err
	}

	var res sequence.LazyResolver
	p := &parser{
		ns:       ns,
		tokenize: tokenize,
		seq:      lex.StripWhitespace(lexer),
	}

	res = func() (ale.Value, data.Sequence, bool) {
		f, ok, err := p.nextValue()
		if err != nil {
			panic(err)
		}
		if ok {
			return f, sequence.NewLazy(res), true
		}
		return data.Null, data.Null, false
	}

	return sequence.NewLazy(res), nil
}
