package types

import (
	"sync/atomic"

	"github.com/kode4food/ale"
)

type (
	// ID uniquely identifies a Type within a process
	ID uint64

	Basic interface {
		ale.Type
		ID() ID
	}

	basic struct {
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

func MakeBasic(name string) Basic {
	return makeBasic(name)
}

func makeBasic(name string) *basic {
	return &basic{
		id:   ID(idCounter.Add(1)),
		name: name,
	}
}

func (b *basic) ID() ID {
	return b.id
}

func (b *basic) Name() string {
	return b.name
}

func (b *basic) Accepts(other ale.Type) bool {
	if other, ok := other.(Basic); ok {
		return b == other || b.id == other.ID()
	}
	return false
}

func (b *basic) Equal(other ale.Type) bool {
	if other, ok := other.(*basic); ok {
		return b == other || b.id == other.id && b.name == other.name
	}
	return false
}
