package data

import (
	"bytes"
	"cmp"
	"slices"

	"github.com/kode4food/ale/types"
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

	// Cons cells are the standard implementation of a Pair. Unlike
	// other Pairs (ex: List, Vector), it is not treated as a Sequence
	Cons struct {
		car Value
		cdr Value
	}
)

// Sorted returns a sorted set of Pairs
func (p Pairs) Sorted() Pairs {
	res := make(Pairs, len(p))
	copy(res, p)
	slices.SortFunc(res, func(l, r Pair) int {
		return cmp.Compare(l.Car().String(), r.Car().String())
	})
	return res
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
func (c *Cons) Equal(v Value) bool {
	if c == v {
		return true
	}
	if v, ok := v.(*Cons); ok {
		return c.car.Equal(v.car) && c.cdr.Equal(v.cdr)
	}
	return false
}

func (c *Cons) String() string {
	var buf bytes.Buffer
	buf.WriteByte('(')
	var next Pair = c
	for {
		buf.WriteString(MaybeQuoteString(next.Car()))
		cdr := next.Cdr()
		if s, ok := cdr.(String); ok {
			buf.WriteString(" . ")
			buf.WriteString(MaybeQuoteString(s))
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
		buf.WriteString(MaybeQuoteString(cdr))
		break
	}
	buf.WriteByte(')')
	return buf.String()
}

func (*Cons) Type() types.Type {
	return types.AnyCons
}

// HashCode returns the hash code for this Cons
func (c *Cons) HashCode() uint64 {
	return HashCode(c.car) * HashCode(c.cdr)
}
