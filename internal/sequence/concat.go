package sequence

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
)

func Concat(s ...data.Sequence) data.Sequence {
	switch len(s) {
	case 0:
		return data.Null
	case 1:
		return s[0]
	default:
		var next LazyResolver
		curr := s[0]
		rest := s[1:]

		next = func() (ale.Value, data.Sequence, bool) {
			var f ale.Value
			var ok bool
			if f, curr, ok = curr.Split(); ok {
				return f, NewLazy(next), true
			}
			if len(rest) == 1 {
				return rest[0].Split()
			}
			curr = rest[0]
			rest = rest[1:]
			return next()
		}
		return NewLazy(next)
	}
}
