package data

import (
	"math/rand"
	"strings"
	"sync/atomic"

	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/types"
)

// A List represents a singly linked List
type List struct {
	first Value
	rest  *List
	count int
	hash  atomic.Uint64
}

var (
	// Null represents the absence of a Value (the empty List)
	Null *List

	nullHash = rand.Uint64()

	// compile-time checks for interface implementation
	_ interface {
		Caller
		Hashed
		Prepender
		RandomAccess
		Reverser
		Typed
	} = Null
)

// NewList creates a new List instance
func NewList(v ...Value) *List {
	var res *List
	for i, u := len(v)-1, 1; i >= 0; i, u = i-1, u+1 {
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
	c := 1
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
	for d, u := e.count, 1; d > 0; e, d, u = e.rest, d-1, u+1 {
		res = &List{
			first: e.Car(),
			rest:  res,
			count: u,
		}
	}
	return res
}

func (l *List) Count() int {
	if l == nil {
		return 0
	}
	return l.count
}

func (l *List) ElementAt(index int) (Value, bool) {
	if l == nil || index > l.count-1 || index < 0 {
		return Null, false
	}

	e := l
	for range index {
		e = e.rest
	}
	return e.Car(), true
}

func (l *List) Call(args ...Value) Value {
	return indexedCall(l, args)
}

func (l *List) CheckArity(argc int) error {
	return CheckRangedArity(1, 2, argc)
}

func (l *List) Equal(other Value) bool {
	if other, ok := other.(*List); ok {
		if l == nil || other == nil || l == other {
			return l == other
		}
		if l.count != other.count {
			return false
		}
		for cl, co := l, other; cl != nil; cl, co = cl.rest, co.rest {
			if cl == co {
				return true
			}
			lh := cl.hash.Load()
			rh := co.hash.Load()
			if lh != 0 && rh != 0 && lh != rh {
				return false
			}
			if !cl.first.Equal(co.first) {
				return false
			}
		}
		return true
	}
	return false
}

func (l *List) String() string {
	if l == nil {
		return lang.ListStart + lang.ListEnd
	}
	var b strings.Builder
	b.WriteString(lang.ListStart)
	b.WriteString(ToQuotedString(l.first))
	for r := l.rest; r != nil; r = r.rest {
		b.WriteString(lang.Space)
		b.WriteString(ToQuotedString(r.first))
	}
	b.WriteString(lang.ListEnd)
	return b.String()
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
	if h := l.hash.Load(); h != 0 {
		return h
	}
	var res uint64 = 0
	for c := l; c != nil; c = c.rest {
		if ch := c.hash.Load(); ch != 0 {
			res ^= ch
			l.hash.Store(res)
			return res
		}
		res ^= HashCode(c.first)
		res ^= HashInt(c.count)
	}
	res ^= nullHash
	l.hash.Store(res)
	return res
}
