package read

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/stdlib"
)

// FromString converts the raw source into unexpanded data structures
func FromString(src api.String) api.Sequence {
	l := Scan(src)
	return FromScanner(l)
}

// FromScanner returns a Lazy Sequence of scanned data structures
func FromScanner(lexer api.Sequence) api.Sequence {
	var res stdlib.LazyResolver
	r := newReader(lexer)

	res = func() (api.Value, api.Sequence, bool) {
		if f, ok := r.nextValue(); ok {
			return f, stdlib.NewLazySequence(res), true
		}
		return api.Nil, api.EmptyList, false
	}

	return stdlib.NewLazySequence(res)
}
