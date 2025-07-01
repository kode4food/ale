package data

import "github.com/kode4food/ale/internal/debug"

// Number describes a numeric value of some kind
type Number interface {
	Hashed
	Typed
	Cmp(Number) Comparison
	Add(Number) Number
	Sub(Number) Number
	Mul(Number) Number
	Div(Number) Number
	Mod(Number) Number
	IsNaN() bool
	IsPosInf() bool
	IsNegInf() bool
}

// purify performs automatic contagion of operands
func purify(l, r Number) (Number, Number) {
	switch l := l.(type) {
	case Integer:
		switch r.(type) {
		case Float:
			return l.float(), r
		case *BigInt:
			return l.bigInt(), r
		case *Ratio:
			return l.ratio(), r
		}

	case Float:
		switch r := r.(type) {
		case Integer:
			return l, r.float()
		case *BigInt:
			return l, r.float()
		case *Ratio:
			return l, r.float()
		}

	case *BigInt:
		switch r := r.(type) {
		case Integer:
			return l, r.bigInt()
		case Float:
			return l.float(), r
		case *Ratio:
			return l.ratio(), r
		}

	case *Ratio:
		switch r := r.(type) {
		case Integer:
			return l, r.ratio()
		case Float:
			return l.float(), r
		case *BigInt:
			return l, r.ratio()
		}
	}
	panic(debug.ProgrammerErrorf("could not purify: %v and %v", l, r))
}
