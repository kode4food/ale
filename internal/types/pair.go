package types

// Pair describes a pair of typed Values
type Pair struct {
	basic
	car  Type
	cdr  Type
	name string
}

// MakeCons declares a new PairType that will only allow a MakeListOf with
// elements of the provided elem Type
func MakeCons(left, right Type) Type {
	return &Pair{
		basic: BasicCons,
		name:  BasicCons.Name(),
		car:   left,
		cdr:   right,
	}
}

func (p *Pair) Car() Type {
	return p.car
}

func (p *Pair) Cdr() Type {
	return p.cdr
}

func (p *Pair) Name() string {
	return p.name
}

func (p *Pair) Accepts(c *Checker, other Type) bool {
	if other, ok := other.(*Pair); ok {
		return p == other ||
			p.basic.Accepts(c, other) &&
				c.AcceptsChild(p.car, other.Car()) &&
				c.AcceptsChild(p.cdr, other.Cdr())
	}
	return false
}

func (p *Pair) Equal(other Type) bool {
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
