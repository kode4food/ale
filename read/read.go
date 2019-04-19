package read

import (
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/stdlib"
)

// FromString converts the raw source into unexpanded data structures
func FromString(src data.String) data.Sequence {
	l := Scan(src)
	return FromScanner(l)
}

// FromScanner returns a Lazy Sequence of scanned data structures
func FromScanner(lexer data.Sequence) data.Sequence {
	var res stdlib.LazyResolver
	r := newReader(lexer)

	res = func() (data.Value, data.Sequence, bool) {
		if f, ok := r.nextValue(); ok {
			return f, stdlib.NewLazySequence(res), true
		}
		return data.Nil, data.EmptyList, false
	}

	return stdlib.NewLazySequence(res)
}
