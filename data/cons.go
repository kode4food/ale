package data

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/lang"
	"github.com/kode4food/ale/internal/types"
)

type (
	// Pair represents the interface for a binary structure, such as a Cons
	Pair interface {
		ale.Typed
		Car() ale.Value
		Cdr() ale.Value
	}

	// Pairs represents multiple pairs
	Pairs []Pair

	// Cons cells are the standard implementation of a Pair. Unlike other Pairs
	// (ex: List, Vector), it is not treated as a Sequence
	Cons struct {
		car ale.Value
		cdr ale.Value
	}
)

var (
	consSalt = rand.Uint64()

	// compile-time checks for interface implementation
	_ interface {
		Hashed
		Pair
		fmt.Stringer
	} = (*Cons)(nil)
)

// NewCons returns a new Cons cell instance
func NewCons(car, cdr ale.Value) *Cons {
	return &Cons{
		car: car,
		cdr: cdr,
	}
}

// Car returns the first element of a Pair
func (c *Cons) Car() ale.Value {
	return c.car
}

// Cdr returns the second element of a Pair
func (c *Cons) Cdr() ale.Value {
	return c.cdr
}

// Equal compares this Cons to another for equality
func (c *Cons) Equal(other ale.Value) bool {
	if other, ok := other.(*Cons); ok {
		return c == other || c.car.Equal(other.car) && c.cdr.Equal(other.cdr)
	}
	return false
}

func (c *Cons) String() string {
	var buf strings.Builder
	buf.WriteString(lang.ListStart)
	var next Pair = c
	for {
		buf.WriteString(ToQuotedString(next.Car()))
		cdr := next.Cdr()
		if s, ok := cdr.(String); ok {
			buf.WriteString(lang.Space + lang.Dot + lang.Space)
			buf.WriteString(ToQuotedString(s))
			break
		}
		if s, ok := cdr.(Sequence); ok && s.IsEmpty() {
			break
		}
		if p, ok := cdr.(Pair); ok {
			buf.WriteString(lang.Space)
			next = p
			continue
		}
		buf.WriteString(lang.Space + lang.Dot + lang.Space)
		buf.WriteString(ToQuotedString(cdr))
		break
	}
	buf.WriteString(lang.ListEnd)
	return buf.String()
}

func (*Cons) Type() ale.Type {
	return types.BasicCons
}

// HashCode returns the hash code for this Cons
func (c *Cons) HashCode() uint64 {
	return consSalt ^ HashCode(c.car) ^ HashCode(c.cdr)
}
