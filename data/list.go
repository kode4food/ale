package data

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"sync/atomic"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/types"
)

// A List represents a singly linked List
type List struct {
	first ale.Value
	rest  *List
	count int
	hash  atomic.Uint64
}

var (
	// Null represents the absence of a Value (the empty List)
	Null *List

	listSalt = rand.Uint64()

	// compile-time checks for interface implementation
	_ interface {
		Hashed
		Indexed
		Prepender
		Procedure
		Reverser
		fmt.Stringer
	} = Null
)

// NewList creates a new List instance
func NewList(vals ...ale.Value) *List {
	var res *List
	for i, u := len(vals)-1, 1; i >= 0; i, u = i-1, u+1 {
		f := vals[i]
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

func (l *List) Car() ale.Value {
	if l == nil {
		return Null
	}
	return l.first
}

func (l *List) Cdr() ale.Value {
	if l == nil {
		return Null
	}
	return l.rest
}

func (l *List) Split() (ale.Value, Sequence, bool) {
	if l == nil {
		return Null, Null, false
	}
	return l.first, l.rest, true
}

func (l *List) Prepend(v ale.Value) Sequence {
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

func (l *List) ElementAt(index int) (ale.Value, bool) {
	if l == nil || index > l.count-1 || index < 0 {
		return Null, false
	}

	e := l
	for range index {
		e = e.rest
	}
	return e.Car(), true
}

func (l *List) CheckArity(argc int) error {
	return CheckRangedArity(1, 2, argc)
}

func (l *List) Call(args ...ale.Value) ale.Value {
	switch len(args) {
	case 1:
		start := int(args[0].(Integer))
		res, ok := l.from(start)
		if !ok {
			panic(fmt.Errorf(ErrInvalidStartIndex, start))
		}
		return res
	case 2:
		start := int(args[0].(Integer))
		end := int(args[1].(Integer))
		curr, ok := l.from(start)
		if !ok {
			panic(fmt.Errorf(ErrInvalidIndexes, start, end))
		}
		res, ok := curr.take(end - start)
		if !ok {
			panic(fmt.Errorf(ErrInvalidIndexes, start, end))
		}
		return res
	default:
		panic(fmt.Errorf(ErrRangedArity, 1, 2, len(args)))
	}
}

func (l *List) from(idx int) (*List, bool) {
	if l == nil || idx < 0 || idx >= l.count {
		return nil, false
	}

	e := l
	for i := 0; i < idx && e != nil; i++ {
		e = e.rest
	}
	return e, true
}

func (l *List) take(count int) (*List, bool) {
	if l == nil || count < 0 || count > l.count {
		return nil, false
	}
	if count == 0 {
		return Null, true
	}

	res := make(Vector, 0, count)
	curr := l
	for range count {
		res = append(res, curr.first)
		curr = curr.rest
	}
	return NewList(res...), true
}

func (l *List) Equal(other ale.Value) bool {
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

func (l *List) Type() ale.Type {
	if l == nil {
		return types.BasicNull
	}
	return types.MakeLiteral(types.BasicList, l)
}

func (l *List) HashCode() uint64 {
	if l == nil {
		return listSalt
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
	res ^= listSalt
	l.hash.Store(res)
	return res
}
