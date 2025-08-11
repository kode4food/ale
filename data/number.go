package data

import (
	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/debug"
)

// Number describes a numeric value of some kind
type Number interface {
	ale.Typed

	// Cmp compares this Number to another Number
	Cmp(Number) Comparison

	// Add adds this Number to another Number
	Add(Number) Number

	// Sub subtracts another Number from this Number
	Sub(Number) Number

	// Mul multiplies this Number by another Number
	Mul(Number) Number

	// Div divides this Number by another Number
	Div(Number) Number

	// Mod calculates the remainder of dividing this Number by another Number
	Mod(Number) Number

	// IsNaN returns whether this Number is not a number (NaN)
	IsNaN() bool

	// IsPosInf returns whether this Number is positive infinity
	IsPosInf() bool

	// IsNegInf returns whether this Number is negative infinity
	IsNegInf() bool
}

// errCouldNotPurify is raised when the purify function cannot convert the two
// operands to a common type. This is a programmer error
const errCouldNotPurify = "could not purify: %v and %v"

// purify performs automatic contagion of operands
func purify(l, r Number) (Number, Number) {
	switch l := l.(type) {
	case Integer:
		return purifyInteger(l, r)
	case Float:
		return purifyFloat(l, r)
	case *BigInt:
		return purifyBigInt(l, r)
	case *Ratio:
		return purifyRatio(l, r)
	default:
		panic(debug.ProgrammerErrorf(errCouldNotPurify, l, r))
	}
}

func purifyInteger(l Integer, r Number) (Number, Number) {
	switch r.(type) {
	case Float:
		return l.float(), r
	case *BigInt:
		return l.bigInt(), r
	case *Ratio:
		return l.ratio(), r
	default:
		panic(debug.ProgrammerErrorf(errCouldNotPurify, l, r))
	}
}

func purifyFloat(l Float, r Number) (Number, Number) {
	switch r := r.(type) {
	case Integer:
		return l, r.float()
	case *BigInt:
		return l, r.float()
	case *Ratio:
		return l, r.float()
	default:
		panic(debug.ProgrammerErrorf(errCouldNotPurify, l, r))
	}
}

func purifyBigInt(l *BigInt, r Number) (Number, Number) {
	switch r := r.(type) {
	case Integer:
		return l, r.bigInt()
	case Float:
		return l.float(), r
	case *Ratio:
		return l.ratio(), r
	default:
		panic(debug.ProgrammerErrorf(errCouldNotPurify, l, r))
	}
}

func purifyRatio(l *Ratio, r Number) (Number, Number) {
	switch r := r.(type) {
	case Integer:
		return l, r.ratio()
	case Float:
		return l.float(), r
	case *BigInt:
		return l, r.ratio()
	default:
		panic(debug.ProgrammerErrorf(errCouldNotPurify, l, r))
	}
}
