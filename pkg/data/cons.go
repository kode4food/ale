package data

import (
	"cmp"
	"strings"

	"github.com/kode4food/ale/internal/types"
	"github.com/kode4food/comb/basics"
)

type (
	// Pair represents the interface for a binary structure, such as a Cons
	Pair interface {
		Value
		Car() Value
		Cdr() Value
	}

	// Pairs represents multiple pairs
	Pairs []Pair

	// Cons cells are the standard implementation of a Pair. Unlike other Pairs
	// (ex: List, Vector), it is not treated as a Sequence
	Cons struct {
		car Value
		cdr Value
	}
)

// compile-time checks for interface implementation
var _ interface {
	Hashed
	Pair
	Typed
} = (*Cons)(nil)

// sorted returns a sorted set of Pairs
func (p Pairs) sorted() Pairs {
	return basics.SortFunc(p, func(l, r Pair) int {
		ls := ToString(l.Car())
		rs := ToString(r.Car())
		return cmp.Compare(ls, rs)
	})
}

// NewCons returns a new Cons cell instance
func NewCons(car, cdr Value) *Cons {
	return &Cons{
		car: car,
		cdr: cdr,
	}
}

// Car returns the first element of a Pair
func (c *Cons) Car() Value {
	return c.car
}

// Cdr returns the second element of a Pair
func (c *Cons) Cdr() Value {
	return c.cdr
}

// Equal compares this Cons to another for equality
func (c *Cons) Equal(other Value) bool {
	if other, ok := other.(*Cons); ok {
		return c == other || c.car.Equal(other.car) && c.cdr.Equal(other.cdr)
	}
	return false
}

func (c *Cons) String() string {
	var buf strings.Builder
	buf.WriteByte('(')
	var next Pair = c
	for {
		buf.WriteString(ToQuotedString(next.Car()))
		cdr := next.Cdr()
		if s, ok := cdr.(String); ok {
			buf.WriteString(" . ")
			buf.WriteString(ToQuotedString(s))
			break
		}
		if s, ok := cdr.(Sequence); ok && s.IsEmpty() {
			break
		}
		if p, ok := cdr.(Pair); ok {
			buf.WriteByte(' ')
			next = p
			continue
		}
		buf.WriteString(" . ")
		buf.WriteString(ToQuotedString(cdr))
		break
	}
	buf.WriteByte(')')
	return buf.String()
}

func (*Cons) Type() types.Type {
	return types.BasicCons
}

// HashCode returns the hash code for this Cons
func (c *Cons) HashCode() uint64 {
	return HashCode(c.car) ^ HashCode(c.cdr)
}
