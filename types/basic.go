package types

import "sync"

type (
	// BasicType accepts a value of its own Type
	BasicType interface {
		Type
		basic() // marker
		Kind() Kind
	}

	// Kind uniquely identifies a BasicType within a process
	Kind uint64

	basic struct {
		name string
		kind Kind
	}
)

var (
	basicKindCounter Kind
	basicKindMutex   sync.Mutex
)

// Basic returns a Type that accepts a Value of its own Type
func Basic(name string) BasicType {
	return &basic{
		name: name,
		kind: nextKind(),
	}
}

func (*basic) basic() {}

func (b *basic) Name() string {
	return b.name
}

func (b *basic) Kind() Kind {
	return b.kind
}

func (b *basic) Accepts(other Type) bool {
	if b == other {
		return true
	}
	if other, ok := other.(BasicType); ok {
		return b.Kind() == other.Kind()
	}
	return false
}

func nextKind() Kind {
	basicKindMutex.Lock()
	defer basicKindMutex.Unlock()
	res := basicKindCounter
	basicKindCounter++
	return res
}
