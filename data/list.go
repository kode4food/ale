package data

import (
	"math/rand"
	"sync/atomic"

	"github.com/kode4food/ale/internal/types"
)

// A List represents a singly linked List
type List struct {
	first Value
	rest  *List
	count Integer
	hash  uint64
}

var (
	// Null represents the absence of a Value (the empty List)
	Null *List

	nullHash = rand.Uint64()

	// compile-time checks for interface implementation
	_ Caller       = Null
	_ Hashed       = Null
	_ Prepender    = Null
	_ RandomAccess = Null
	_ Reverser     = Null
	_ Typed        = Null
)

// NewList creates a new List instance
func NewList(v ...Value) *List {
	var res *List
	for i, u := len(v)-1, Integer(1); i >= 0; i, u = i-1, u+1 {
		f := v[i]
		res = &List{
			first: f,
			rest:  res,
			count: u,
		}
	}
	return res
}

func (l *List) IsEmpty() bool {
	return l == nil
}

func (l *List) Car() Value {
	if l == nil {
		return Null
	}
	return l.first
}

func (l *List) Cdr() Value {
	if l == nil {
		return Null
	}
	return l.rest
}

func (l *List) Split() (Value, Sequence, bool) {
	if l == nil {
		return Null, Null, false
	}
	return l.first, l.rest, true
}

func (l *List) Prepend(v Value) Sequence {
	c := Integer(1)
	if l != nil {
		c += l.count
	}
	return &List{
		first: v,
		rest:  l,
		count: c,
	}
}

func (l *List) Reverse() Sequence {
	if l == nil || l.count <= 1 {
		return l
	}

	var res *List
	e := l
	for d, u := e.count, Integer(1); d > 0; e, d, u = e.rest, d-1, u+1 {
		res = &List{
			first: e.Car(),
			rest:  res,
			count: u,
		}
	}
	return res
}

func (l *List) Count() Integer {
	if l == nil {
		return 0
	}
	return l.count
}

func (l *List) ElementAt(index Integer) (Value, bool) {
	if l == nil || index > l.count-1 || index < 0 {
		return Null, false
	}

	e := l
	for i := Integer(0); i < index; i++ {
		e = e.rest
	}
	return e.Car(), true
}

func (l *List) Call(args ...Value) Value {
	return indexedCall(l, args)
}

func (l *List) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

func (l *List) Equal(other Value) bool {
	if l == nil {
		return l == other
	}
	cr, ok := other.(*List)
	if !ok || l.count != cr.count {
		return false
	}
	for cl := l; cl != nil; cl, cr = cl.rest, cr.rest {
		if cl == cr {
			return true
		}
		if !cl.first.Equal(cr.first) {
			return false
		}
	}
	return true
}

func (l *List) String() string {
	return MakeSequenceStr(l)
}

func (l *List) Type() types.Type {
	if l == nil {
		return types.BasicNull
	}
	return types.BasicList
}

func (l *List) HashCode() uint64 {
	if l == nil {
		return nullHash
	}
	if h := atomic.LoadUint64(&l.hash); h != 0 {
		return h
	}
	var res uint64 = 0
	for c := l; c != nil; c = c.rest {
		if ch := atomic.LoadUint64(&c.hash); ch != 0 {
			res ^= ch
			atomic.StoreUint64(&l.hash, res)
			return res
		}
		res ^= HashCode(c.first) ^ (uint64(1) << (c.count % 64))
	}
	res ^= nullHash
	atomic.StoreUint64(&l.hash, res)
	return res
}
