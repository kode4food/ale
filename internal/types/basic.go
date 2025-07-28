package types

import (
	"sync/atomic"

	"github.com/kode4food/ale"
)

type (
	// ID uniquely identifies a Type within a process
	ID uint64

	basic interface {
		ale.Type
		ID() ID
	}

	Basic struct {
		name string
		id   ID
	}
)

var (
	BasicBoolean   = makeBasic("boolean")
	BasicBytes     = makeBasic("bytes")
	BasicKeyword   = makeBasic("keyword")
	BasicProcedure = makeBasic("procedure")
	BasicNull      = makeBasic("null")
	BasicNumber    = makeBasic("number")
	BasicString    = makeBasic("string")
	BasicSymbol    = makeBasic("symbol")
	BasicList      = makeBasic("list")
	BasicObject    = makeBasic("object")
	BasicCons      = makeBasic("cons")
	BasicVector    = makeBasic("vector")
	BasicUnion     = makeBasic("union")

	idCounter atomic.Uint64
)

func MakeBasic(name string) ale.Type {
	return makeBasic(name)
}

func makeBasic(name string) *Basic {
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

func (b *Basic) Accepts(other ale.Type) bool {
	if other, ok := other.(basic); ok {
		return b == other || b.id == other.ID()
	}
	return false
}

func (b *Basic) Equal(other ale.Type) bool {
	if other, ok := other.(*Basic); ok {
		return b == other || b.id == other.id && b.name == other.name
	}
	return false
}
