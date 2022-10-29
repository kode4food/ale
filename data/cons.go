package data

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/kode4food/ale/types"
)

type (
	// Pair represents the interface for a binary structure, such as a Cons
	Pair interface {
		Value
		Car() Value
		Cdr() Value
	}

	// Cons represents the most basic implementation of a Pair
	Cons interface {
		cons() // marker
		Pair
	}

	// Pairs represents multiple pairs
	Pairs []Pair

	// Cons cells are the standard implementation of a Pair. Unlike
	// other Pairs (ex: List, Vector), it is not treated as a Sequence
	cons struct {
		car Value
		cdr Value
	}
)

// Sorted returns a sorted set of Pairs
func (p Pairs) Sorted() Pairs {
	res := p[:]
	sort.Slice(res, func(l, r int) bool {
		ls := fmt.Sprintf("%s", res[l].Car().String())
		rs := fmt.Sprintf("%s", res[r].Car().String())
		return ls < rs
	})
	return res
}

// NewCons returns a new Cons cell instance
func NewCons(car, cdr Value) Cons {
	return &cons{
		car: car,
		cdr: cdr,
	}
}

func (*cons) cons() {}

// Car returns the first element of a Pair
func (c *cons) Car() Value {
	return c.car
}

// Cdr returns the second element of a Pair
func (c *cons) Cdr() Value {
	return c.cdr
}

// Equal compares this Cons to another for equality
func (c *cons) Equal(v Value) bool {
	if c == v {
		return true
	}
	if v, ok := v.(*cons); ok {
		return c.car.Equal(v.car) && c.cdr.Equal(v.cdr)
	}
	return false
}

func (c *cons) String() string {
	var buf bytes.Buffer
	buf.WriteByte('(')
	var next Pair = c
	for {
		buf.WriteString(MaybeQuoteString(next.Car()))
		cdr := next.Cdr()
		if s, ok := cdr.(Sequence); ok && s.IsEmpty() {
			break
		} else if p, ok := cdr.(Pair); ok {
			buf.WriteByte(' ')
			next = p
		} else {
			buf.WriteString(" . ")
			buf.WriteString(MaybeQuoteString(cdr))
			break
		}
	}
	buf.WriteByte(')')
	return buf.String()
}

func (*cons) Type() types.Type {
	return types.AnyCons
}

// HashCode returns the hash code for this Cons
func (c *cons) HashCode() uint64 {
	return HashCode(c.car) * HashCode(c.cdr)
}
