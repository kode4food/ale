package types

import "github.com/kode4food/ale"

// Pair describes a pair of typed Values
type Pair struct {
	Basic
	car  ale.Type
	cdr  ale.Type
	name string
}

// MakeCons declares a new PairType that will only allow a MakeListOf with
// elements of the provided elem Type
func MakeCons(left, right ale.Type) ale.Type {
	return &Pair{
		Basic: BasicCons,
		name:  BasicCons.Name(),
		car:   left,
		cdr:   right,
	}
}

func (p *Pair) Car() ale.Type {
	return p.car
}

func (p *Pair) Cdr() ale.Type {
	return p.cdr
}

func (p *Pair) Name() string {
	return p.name
}

func (p *Pair) Accepts(other ale.Type) bool {
	if other, ok := other.(*Pair); ok {
		return p == other || compoundAccepts(p, other)
	}
	return false
}

func (p *Pair) accepts(c *checker, other ale.Type) bool {
	if other, ok := other.(*Pair); ok {
		return p == other ||
			p.Basic.Accepts(other.Basic) &&
				c.acceptsChild(p.car, other.Car()) &&
				c.acceptsChild(p.cdr, other.Cdr())
	}
	return false
}

func (p *Pair) Equal(other ale.Type) bool {
	if other, ok := other.(*Pair); ok {
		if p == other {
			return true
		}
		return p.name == other.name &&
			p.Basic.Equal(other.Basic) &&
			p.car.Equal(other.car) &&
			p.cdr.Equal(other.cdr)
	}
	return false
}
