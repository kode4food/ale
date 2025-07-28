package types

import "github.com/kode4food/ale"

// Pair describes a pair of typed Values
type Pair struct {
	basic
	car  ale.Type
	cdr  ale.Type
	name string
}

// MakeCons declares a new PairType that will only allow a MakeListOf with
// elements of the provided elem Type
func MakeCons(left, right ale.Type) ale.Type {
	return &Pair{
		basic: BasicCons,
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
			p.basic.Accepts(other.basic) &&
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
			p.basic.Equal(other.basic) &&
			p.car.Equal(other.car) &&
			p.cdr.Equal(other.cdr)
	}
	return false
}
