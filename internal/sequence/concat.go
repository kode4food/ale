package sequence

import "github.com/kode4food/ale/pkg/data"

func Concat(s ...data.Sequence) data.Sequence {
	switch len(s) {
	case 0:
		return data.Null
	case 1:
		return s[0]
	}

	var next LazyResolver
	curr := s[0]
	rest := s[1:]

	next = func() (data.Value, data.Sequence, bool) {
		var f data.Value
		var ok bool
		f, curr, ok = curr.Split()
		if ok {
			return f, NewLazy(next), true
		}
		switch len(rest) {
		case 0:
			return data.Null, data.Null, false
		case 1:
			return rest[0].Split()
		default:
			curr = rest[0]
			rest = rest[1:]
			return next()
		}
	}
	return NewLazy(next)
}
