package sequence

import (
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/do"
)

type (
	// LazyResolver is used to resolve the elements of a lazy Sequence
	LazyResolver func() (data.Value, data.Sequence, bool)

	lazySequence struct {
		once     do.Action
		resolver LazyResolver

		ok     bool
		result data.Value
		rest   data.Sequence
	}
)

// NewLazy creates a new lazy Sequence based on the provided resolver
func NewLazy(r LazyResolver) data.Sequence {
	return &lazySequence{
		once:     do.Once(),
		resolver: r,
		result:   data.Nil,
		rest:     data.EmptyList,
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

func (l *lazySequence) First() data.Value {
	return l.resolve().result
}

func (l *lazySequence) Rest() data.Sequence {
	return l.resolve().rest
}

func (l *lazySequence) Split() (data.Value, data.Sequence, bool) {
	r := l.resolve()
	return r.result, r.rest, l.ok
}

func (l *lazySequence) Car() data.Value {
	return l.resolve().result
}

func (l *lazySequence) Cdr() data.Value {
	return l.resolve().rest
}

func (l *lazySequence) Prepend(v data.Value) data.Sequence {
	return &lazySequence{
		once:   do.Never(),
		ok:     true,
		result: v,
		rest:   l,
	}
}

func (l *lazySequence) Type() data.Name {
	return "lazy-sequence"
}

func (l *lazySequence) Equal(v data.Value) bool {
	if v, ok := v.(*lazySequence); ok {
		return l == v
	}
	return false
}

func (l *lazySequence) String() string {
	return data.DumpString(l)
}
