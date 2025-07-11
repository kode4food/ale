package types

import "sync/atomic"

type (
	// ID uniquely identifies a Type within a process
	ID uint64

	basic interface {
		Type
		ID() ID
	}

	Basic struct {
		name string
		id   ID
	}
)

var (
	BasicBoolean   = MakeBasic("boolean")
	BasicBytes     = MakeBasic("bytes")
	BasicKeyword   = MakeBasic("keyword")
	BasicProcedure = MakeBasic("procedure")
	BasicNull      = MakeBasic("null")
	BasicNumber    = MakeBasic("number")
	BasicString    = MakeBasic("string")
	BasicSymbol    = MakeBasic("symbol")
	BasicList      = MakeBasic("list")
	BasicObject    = MakeBasic("object")
	BasicCons      = MakeBasic("cons")
	BasicVector    = MakeBasic("vector")
	BasicUnion     = MakeBasic("union")

	idCounter atomic.Uint64
)

func MakeBasic(name string) *Basic {
	return &Basic{
		id:   ID(idCounter.Add(1)),
		name: name,
	}
}

func (b *Basic) ID() ID {
	return b.id
}

func (b *Basic) Name() string {
	return b.name
}

func (b *Basic) Accepts(_ *Checker, other Type) bool {
	if other, ok := other.(basic); ok {
		return b == other || b.id == other.ID()
	}
	return false
}

func (b *Basic) Equal(other Type) bool {
	if other, ok := other.(*Basic); ok {
		return b == other || b.id == other.id && b.name == other.name
	}
	return false
}
