package types

type (
	// PairType describes a pair of typed Values
	PairType interface {
		Type
		pair() // marker
		Car() Type
		Cdr() Type
	}

	pair struct {
		BasicType
		car Type
		cdr Type
	}

	namedPair struct {
		*pair
		name string
	}
)

func (*pair) pair() {}

// Cons declares a new PairType that will only allow a ListOf with elements of
// the provided elem Type
func Cons(left, right Type) Type {
	return &pair{
		BasicType: AnyCons,
		car:       left,
		cdr:       right,
	}
}

func (p *pair) Car() Type {
	return p.car
}

func (p *pair) Cdr() Type {
	return p.cdr
}

func (p *pair) Accepts(c *Checker, other Type) bool {
	if other, ok := other.(PairType); ok {
		return other.IsA(p.BasicType) &&
			c.AcceptsChild(p.car, other.Car()) &&
			c.AcceptsChild(p.cdr, other.Cdr())
	}
	return false
}

func (p *namedPair) Name() string {
	return p.name
}

func (p *namedPair) Accepts(c *Checker, other Type) bool {
	if other, ok := other.(*namedPair); ok {
		return p.pair.Accepts(c, other.pair)
	}
	return p.pair.Accepts(c, other)
}
