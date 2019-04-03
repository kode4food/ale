package stdlib

import "gitlab.com/kode4food/ale/api"

type (
	// LazyResolver is used to resolve the elements of a lazy Sequence
	LazyResolver func() (api.Value, api.Sequence, bool)

	lazySequence struct {
		once     Do
		resolver LazyResolver

		isSeq  bool
		result api.Value
		rest   api.Sequence
	}
)

// NewLazySequence creates a new lazy Sequence based on the provided resolver
func NewLazySequence(r LazyResolver) api.Sequence {
	return &lazySequence{
		once:     Once(),
		resolver: r,
		result:   api.Nil,
		rest:     api.EmptyList,
	}
}

func (l *lazySequence) resolve() *lazySequence {
	l.once(func() {
		l.result, l.rest, l.isSeq = l.resolver()
		l.resolver = nil
	})
	return l
}

func (l *lazySequence) IsSequence() bool {
	return l.resolve().isSeq
}

func (l *lazySequence) First() api.Value {
	return l.resolve().result
}

func (l *lazySequence) Rest() api.Sequence {
	return l.resolve().rest
}

func (l *lazySequence) Split() (api.Value, api.Sequence, bool) {
	r := l.resolve()
	return r.result, r.rest, l.isSeq
}

func (l *lazySequence) Prepend(v api.Value) api.Sequence {
	return &lazySequence{
		once:   Never(),
		isSeq:  true,
		result: v,
		rest:   l,
	}
}

func (l *lazySequence) Type() api.Name {
	return "lazy-sequence"
}

func (l *lazySequence) String() string {
	return api.DumpString(l)
}
