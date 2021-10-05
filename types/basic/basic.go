package basic

import (
	"sync"

	"github.com/kode4food/ale/types"
)

type basic struct {
	name string
	kind types.Kind
}

var (
	kindCounter types.Kind
	kindMutex   sync.Mutex
)

// New returns a Basic type that accepts a Value of its own Basic type
func New(name string) types.Basic {
	return &basic{
		name: name,
		kind: nextKind(),
	}
}

func (*basic) basic() {}

func (b *basic) Name() string {
	return b.name
}

func (b *basic) Kind() types.Kind {
	return b.kind
}

func (b *basic) Accepts(other types.Type) bool {
	if b == other {
		return true
	}
	switch other := other.(type) {
	case types.Basic:
		return b.kind == other.Kind()
	case types.Extended:
		return b.Accepts(other.Base())
	default:
		return false
	}
}

func nextKind() types.Kind {
	kindMutex.Lock()
	defer kindMutex.Unlock()
	res := kindCounter
	kindCounter++
	return res
}

// Basic Types
var (
	Bool    = New("boolean")
	Keyword = New("keyword")
	Lambda  = New("lambda")
	List    = New("list")
	Null    = New("null")
	Number  = New("number")
	Object  = New("object")
	Pair    = New("pair")
	String  = New("string")
	Symbol  = New("symbol")
	Vector  = New("vector")
)
