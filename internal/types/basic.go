package types

import "github.com/google/uuid"

type (
	// Kind uniquely identifies a Type within a process
	Kind uuid.UUID

	basic interface {
		Type
		Kind() Kind
	}

	Basic struct {
		name string
		kind Kind
	}
)

var (
	BasicBoolean   = MakeBasic("boolean")
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
)

func MakeBasic(name string) *Basic {
	return &Basic{
		kind: newKind(),
		name: name,
	}
}

func (b *Basic) Kind() Kind {
	return b.kind
}

func (b *Basic) Name() string {
	return b.name
}

func (b *Basic) Accepts(_ *Checker, other Type) bool {
	if b == other {
		return true
	}
	if other, ok := other.(basic); ok {
		return b.kind == other.Kind()
	}
	return false
}

func (b *Basic) Equal(other Type) bool {
	if b == other {
		return true
	}
	if other, ok := other.(*Basic); ok {
		return b.kind == other.kind && b.name == other.name
	}
	return false
}

func newKind() Kind {
	return Kind(uuid.New())
}
