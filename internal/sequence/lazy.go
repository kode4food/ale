package sequence

import (
	"github.com/kode4food/ale/internal/sync"
	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/ale/pkg/data"
)

type (
	// LazyResolver is used to resolve the elements of a lazy Sequence
	LazyResolver func() (data.Value, data.Sequence, bool)

	lazySequence struct {
		once     sync.Action
		resolver LazyResolver

		result data.Value
		rest   data.Sequence
		ok     bool
	}
)

var (
	lazySequenceType = types.MakeBasic("lazy-sequence")

	// compile-time check for interface implementation
	_ data.Prepender = (*lazySequence)(nil)
)

// NewLazy creates a new lazy Sequence based on the provided resolver
func NewLazy(r LazyResolver) data.Sequence {
	return &lazySequence{
		once:     sync.Once(),
		resolver: r,
	}
}

func MakeLazyResolver(p data.Procedure) LazyResolver {
	return func() (data.Value, data.Sequence, bool) {
		r := p.Call()
		if r != data.Null {
			s := r.(data.Sequence)
			if sf, sr, ok := s.Split(); ok {
				return sf, sr, true
			}
		}
		return data.Null, data.Null, false
	}
}

func (l *lazySequence) resolve() *lazySequence {
	l.once(func() {
		l.result, l.rest, l.ok = l.resolver()
		l.resolver = nil
	})
	return l
}

func (l *lazySequence) IsEmpty() bool {
	return !l.resolve().ok
}

func (l *lazySequence) Car() data.Value {
	return l.resolve().result
}

func (l *lazySequence) Cdr() data.Value {
	return l.resolve().rest
}

func (l *lazySequence) Split() (data.Value, data.Sequence, bool) {
	r := l.resolve()
	return r.result, r.rest, l.ok
}

func (l *lazySequence) Prepend(v data.Value) data.Sequence {
	return &lazySequence{
		once:   sync.Never(),
		ok:     true,
		result: v,
		rest:   l,
	}
}

func (l *lazySequence) Type() types.Type {
	return lazySequenceType
}

func (l *lazySequence) Equal(other data.Value) bool {
	return l == other
}

func (l *lazySequence) Get(key data.Value) (data.Value, bool) {
	return data.DumpMapped(l).Get(key)
}
