package stdlib

import "github.com/kode4food/ale/data"

type (
	// LazyResolver is used to resolve the elements of a lazy Sequence
	LazyResolver func() (data.Value, data.Sequence, bool)

	lazySequence struct {
		once     Do
		resolver LazyResolver

		ok     bool
		result data.Value
		rest   data.Sequence
	}
)

// NewLazySequence creates a new lazy Sequence based on the provided resolver
func NewLazySequence(r LazyResolver) data.Sequence {
	return &lazySequence{
		once:     Once(),
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
	return data.SequenceCar(l)
}

func (l *lazySequence) Cdr() data.Value {
	return data.SequenceCdr(l)
}

func (l *lazySequence) Prepend(v data.Value) data.Sequence {
	return &lazySequence{
		once:   Never(),
		ok:     true,
		result: v,
		rest:   l,
	}
}

func (l *lazySequence) Type() data.Name {
	return "lazy-sequence"
}

func (l *lazySequence) String() string {
	return data.DumpString(l)
}
