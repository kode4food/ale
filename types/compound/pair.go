package compound

import (
	"github.com/kode4food/ale/types"
	"github.com/kode4food/ale/types/basic"
	"github.com/kode4food/ale/types/extended"
)

type (
	// PairType describes a pair of typed Values
	PairType interface {
		types.Extended
		pair() // marker
		Car() types.Type
		Cdr() types.Type
	}

	pair struct {
		types.Extended
		car types.Type
		cdr types.Type
	}

	namedPair struct {
		*pair
		name string
	}
)

func (*pair) pair() {}

// Cons declares a new PairType that will only allow a List with
// elements of the provided elem Type
func Cons(left, right types.Type) types.Type {
	base := basic.Cons
	return &pair{
		Extended: extended.New(base),
		car:      left,
		cdr:      right,
	}
}

func (p *pair) Car() types.Type {
	return p.car
}

func (p *pair) Cdr() types.Type {
	return p.cdr
}

func (p *pair) Accepts(c types.Checker, other types.Type) bool {
	if p == other {
		return true
	}
	if other, ok := other.(PairType); ok {
		return c.Check(p.Extended).Accepts(other) != nil &&
			c.Check(p.car).Accepts(other.Car()) != nil &&
			c.Check(p.cdr).Accepts(other.Cdr()) != nil
	}
	return false
}

func (p *namedPair) Name() string {
	return p.name
}

func (p *namedPair) Accepts(c types.Checker, other types.Type) bool {
	if p == other {
		return true
	}
	if other, ok := other.(*namedPair); ok {
		return p.pair.Accepts(c, other.pair)
	}
	return p.pair.Accepts(c, other)
}
