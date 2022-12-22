package types

import "github.com/google/uuid"

type (
	BasicType interface {
		Type
		Kind() Kind
	}

	// Kind uniquely identifies a Type within a process
	Kind uuid.UUID

	basic struct {
		kind Kind
		name string
	}
)

var (
	Bool      = Basic("boolean")
	Keyword   = Basic("keyword")
	Lambda    = Basic("lambda")
	Null      = Basic("null")
	Number    = Basic("number")
	String    = Basic("string")
	Symbol    = Basic("symbol")
	AnyList   = Basic("list")
	AnyObject = Basic("object")
	AnyCons   = Basic("cons")
	AnyVector = Basic("vector")
	AnyUnion  = Basic("union")
)

func Basic(name string) BasicType {
	return &basic{
		kind: newKind(),
		name: name,
	}
}

func (b *basic) Kind() Kind {
	return b.kind
}

func (b *basic) Name() string {
	return b.name
}

func (b *basic) IsA(other BasicType) bool {
	return b.kind == other.Kind()
}

func (b *basic) Accepts(_ *Checker, other Type) bool {
	return other.IsA(b)
}

func newKind() Kind {
	return Kind(uuid.New())
}
