package data

import (
	"fmt"
	"math"
	"strconv"
)

type (
	// Number describes a numeric value of some kind
	Number interface {
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

	// Integer represents a 64-bit integer
	Integer int64

	// Float represents a 64-bit floating point number
	Float float64
)

// Error messages
const (
	ExpectedInteger = "value is not an integer: %s"
	ExpectedFloat   = "value is not a float: %s"
	UnknownType     = "unknown type: %v"
)

// ParseInteger attempts to parse a string representing a 64-bit integer
func ParseInteger(s String) Integer {
	ns := string(s)
	if res, err := strconv.ParseInt(ns, 0, 64); err == nil {
		return Integer(res)
	}
	panic(fmt.Errorf(ExpectedInteger, s))
}

// Cmp compares this Integer to another Number
func (i Integer) Cmp(n Number) Comparison {
	switch typed := n.(type) {
	case Integer:
		if i > typed {
			return GreaterThan
		}
		if i < typed {
			return LessThan
		}
		return EqualTo
	case Float:
		if math.IsNaN(float64(typed)) {
			return Incomparable
		}
		l := Float(i)
		if l > typed {
			return GreaterThan
		}
		if l < typed {
			return LessThan
		}
		return EqualTo
	default:
		panic(fmt.Errorf(UnknownType, n))
	}
}

// Add adds this Integer to another Number
func (i Integer) Add(n Number) Number {
	switch typed := n.(type) {
	case Integer:
		return i + typed
	case Float:
		return Float(i) + typed
	default:
		panic(fmt.Errorf(UnknownType, n))
	}
}

// Sub subtracts another Number from this Integer
func (i Integer) Sub(n Number) Number {
	switch typed := n.(type) {
	case Integer:
		return i - typed
	case Float:
		return Float(i) - typed
	default:
		panic(fmt.Errorf(UnknownType, n))
	}
}

// Mul multiples this Integer by another Number
func (i Integer) Mul(n Number) Number {
	switch typed := n.(type) {
	case Integer:
		return i * typed
	case Float:
		return Float(i) * typed
	default:
		panic(fmt.Errorf(UnknownType, n))
	}
}

// Div divides this Integer by another Number
func (i Integer) Div(n Number) Number {
	switch typed := n.(type) {
	case Integer:
		return i / typed
	case Float:
		return Float(i) / typed
	default:
		panic(fmt.Errorf(UnknownType, n))
	}
}

// Mod calculates the remainder of dividing this Integer by another Number
func (i Integer) Mod(n Number) Number {
	switch typed := n.(type) {
	case Integer:
		return i % typed
	case Float:
		return Float(math.Mod(float64(i), float64(typed)))
	default:
		panic(fmt.Errorf(UnknownType, n))
	}
}

// IsNaN tells you that this Integer is, in fact, a Number
func (i Integer) IsNaN() bool { return false }

// IsPosInf tells you that this Integer is not positive infinity
func (i Integer) IsPosInf() bool { return false }

// IsNegInf tells you that this Integer is not negative infinity
func (i Integer) IsNegInf() bool { return false }

// String converts this Integer to a string
func (i Integer) String() string { return fmt.Sprintf("%d", i) }

// ParseFloat attempts to parse a string representing a 64-bit float
func ParseFloat(s String) Float {
	ns := string(s)
	if res, err := strconv.ParseFloat(ns, 64); err == nil {
		return Float(res)
	}
	panic(fmt.Errorf(ExpectedFloat, s))
}

// Cmp compares this Float to another Number
func (f Float) Cmp(n Number) Comparison {
	if math.IsNaN(float64(f)) {
		return Incomparable
	}
	switch typed := n.(type) {
	case Float:
		if math.IsNaN(float64(typed)) {
			return Incomparable
		}
		if f > typed {
			return GreaterThan
		}
		if f < typed {
			return LessThan
		}
		return EqualTo
	case Integer:
		r := Float(typed)
		if f > r {
			return GreaterThan
		}
		if f < r {
			return LessThan
		}
		return EqualTo
	default:
		panic(fmt.Errorf(UnknownType, n))
	}
}

// Add adds this Float to another Number
func (f Float) Add(n Number) Number {
	switch typed := n.(type) {
	case Float:
		return f + typed
	case Integer:
		return f + Float(typed)
	default:
		panic(fmt.Errorf(UnknownType, n))
	}
}

// Sub subtracts another Number from this Float
func (f Float) Sub(n Number) Number {
	switch typed := n.(type) {
	case Float:
		return f - typed
	case Integer:
		return f - Float(typed)
	default:
		panic(fmt.Errorf(UnknownType, n))
	}
}

// Mul multiplies this Float by another Number
func (f Float) Mul(n Number) Number {
	switch typed := n.(type) {
	case Float:
		return f * typed
	case Integer:
		return f * Float(typed)
	default:
		panic(fmt.Errorf(UnknownType, n))
	}
}

// Div divides this Float by another Number
func (f Float) Div(n Number) Number {
	switch typed := n.(type) {
	case Float:
		return f / typed
	case Integer:
		return f / Float(typed)
	default:
		panic(fmt.Errorf(UnknownType, n))
	}
}

// Mod calculates the remainder of dividing this Float by another Number
func (f Float) Mod(n Number) Number {
	switch typed := n.(type) {
	case Float:
		return Float(math.Mod(float64(f), float64(typed)))
	case Integer:
		return Float(math.Mod(float64(f), float64(typed)))
	default:
		panic(fmt.Errorf(UnknownType, n))
	}
}

// IsNaN returns whether or not this Float is not a number
func (f Float) IsNaN() bool {
	return math.IsNaN(float64(f))
}

// IsPosInf returns whether or not this Float represents positive infinity
func (f Float) IsPosInf() bool {
	return math.IsInf(float64(f), 1)
}

// IsNegInf returns whether or not this Float represents negative infinity
func (f Float) IsNegInf() bool {
	return math.IsInf(float64(f), -1)
}

// String converts this Float to a string
func (f Float) String() string {
	return fmt.Sprintf("%g", f)
}
