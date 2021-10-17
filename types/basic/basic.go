package basic

import (
	"encoding/binary"
	"math/rand"
	"sync"

	"github.com/kode4food/ale/types"
)

type basic struct {
	name string
	kind types.Kind
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
	Cons    = New("cons")
	String  = New("string")
	Symbol  = New("symbol")
	Vector  = New("vector")
)

var (
	kindSequence uint32
	kindMutex    sync.Mutex
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

func (b *basic) Accepts(c types.Checker, other types.Type) bool {
	if b == other {
		return true
	}
	switch other := other.(type) {
	case types.Basic:
		return b.kind == other.Kind()
	case types.Extended:
		return b.Accepts(c, other.Base())
	default:
		return false
	}
}

func nextKind() types.Kind {
	var res types.Kind
	next := nextKindSequence()
	binary.BigEndian.PutUint32(res[0:], next)
	rand.Read(res[4:])
	return res
}

func nextKindSequence() uint32 {
	kindMutex.Lock()
	defer kindMutex.Unlock()
	res := kindSequence
	kindSequence++
	return res
}
