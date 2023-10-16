package read

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

// FromString converts the raw source into unexpanded data structures
func FromString(src data.String) data.Sequence {
	l := Scan(src)
	return FromScanner(l)
}

// FromScanner returns a Lazy Sequence of scanned data structures
func FromScanner(lexer data.Sequence) data.Sequence {
	var res sequence.LazyResolver
	r := newReader(lexer)

	res = func() (data.Value, data.Sequence, bool) {
		if f, ok := r.nextValue(); ok {
			return f, sequence.NewLazy(res), true
		}
		return data.Null, data.Null, false
	}

	return sequence.NewLazy(res)
}

// Scan creates a filtered Lexer Sequence for the Read function
func Scan(src data.String) data.Sequence {
	return sequence.Filter(Tokens(src), noWhitespace)
}

func noWhitespace(v data.Value) bool {
	return !v.(*Token).isWhitespace()
}
