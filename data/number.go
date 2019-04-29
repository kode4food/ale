package data

import (
	"fmt"
)

// Number describes a numeric value of some kind
type Number interface {
	Value
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

// Error messages
const (
	CouldNotPurify = "could not purify: %v and %v"
)

// purify performs automatic contagion of operands
func purify(l, r Number) (Number, Number) {
	switch lt := l.(type) {
	case Integer:
		switch r.(type) {
		case Float:
			return lt.float(), r
		case *BigInt:
			return lt.bigInt(), r
		case *Ratio:
			return lt.ratio(), r
		}

	case Float:
		switch rt := r.(type) {
		case Integer:
			return l, rt.float()
		case *BigInt:
			return l, rt.float()
		case *Ratio:
			return l, rt.float()
		}

	case *BigInt:
		switch rt := r.(type) {
		case Integer:
			return l, rt.bigInt()
		case Float:
			return lt.float(), r
		case *Ratio:
			return lt.ratio(), r
		}

	case *Ratio:
		switch rt := r.(type) {
		case Integer:
			return l, rt.ratio()
		case Float:
			return lt.float(), r
		case *BigInt:
			return l, rt.ratio()
		}
	}

	panic(fmt.Errorf(CouldNotPurify, l, r))
}
