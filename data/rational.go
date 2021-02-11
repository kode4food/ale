package data

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
)

type (
	// Float represents a 64-bit floating point number
	Float float64

	// Ratio represents a number having a numerator and denominator
	Ratio big.Rat
)

// Error messages
const (
	ErrExpectedFloat = "value is not a float: %s"
	ErrExpectedRatio = "value is not a ratio: %s"
)

// ParseFloat attempts to parse a string representing a float
func ParseFloat(s string) (Number, error) {
	res, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, fmt.Errorf(ErrExpectedFloat, s)
	}
	return Float(res), nil
}

// MustParseFloat forcefully parses a string representing a float
func MustParseFloat(s string) Number {
	if res, err := ParseFloat(s); err != nil {
		panic(err)
	} else {
		return res
	}
}

// Cmp compares this Float to another Number
func (l Float) Cmp(r Number) Comparison {
	if math.IsNaN(float64(l)) {
		return Incomparable
	}
	if rf, ok := r.(Float); ok {
		if math.IsNaN(float64(rf)) {
			return Incomparable
		}
		if l > rf {
			return GreaterThan
		}
		if l < rf {
			return LessThan
		}
		return EqualTo
	}
	pl, pr := purify(l, r)
	return pl.Cmp(pr)
}

// Add adds this Float to another Number
func (l Float) Add(r Number) Number {
	if rf, ok := r.(Float); ok {
		return l + rf
	}
	pl, pr := purify(l, r)
	return pl.Add(pr)
}

// Sub subtracts another Number from this Float
func (l Float) Sub(r Number) Number {
	if rf, ok := r.(Float); ok {
		return l - rf
	}
	pl, pr := purify(l, r)
	return pl.Sub(pr)
}

// Mul multiplies this Float by another Number
func (l Float) Mul(r Number) Number {
	if rf, ok := r.(Float); ok {
		return l * rf
	}
	pl, pr := purify(l, r)
	return pl.Mul(pr)
}

// Div divides this Float by another Number
func (l Float) Div(r Number) Number {
	if rf, ok := r.(Float); ok {
		return l / rf
	}
	pl, pr := purify(l, r)
	return pl.Div(pr)
}

// Mod calculates the remainder of dividing this Float by another Number
func (l Float) Mod(r Number) Number {
	if rf, ok := r.(Float); ok {
		return Float(math.Mod(float64(l), float64(rf)))
	}
	pl, pr := purify(l, r)
	return pl.Mod(pr)
}

// IsNaN returns whether this Float is not a number
func (l Float) IsNaN() bool {
	return math.IsNaN(float64(l))
}

// IsPosInf returns whether this Float represents positive infinity
func (l Float) IsPosInf() bool {
	return math.IsInf(float64(l), 1)
}

// IsNegInf returns whether this Float represents negative infinity
func (l Float) IsNegInf() bool {
	return math.IsInf(float64(l), -1)
}

// Equal compares this Float to another for equality
func (l Float) Equal(r Value) bool {
	if r, ok := r.(Float); ok {
		return l == r
	}
	return false
}

// String converts this Float to a string
func (l Float) String() string {
	return fmt.Sprintf("%g", l)
}

// HashCode returns a hash code for this Float
func (l Float) HashCode() uint64 {
	return uint64(l)
}

// ParseRatio attempts to parse a string representing a ratio
func ParseRatio(s string) (Number, error) {
	if res, ok := new(big.Rat).SetString(s); ok {
		return (*Ratio)(res), nil
	}
	return nil, fmt.Errorf(ErrExpectedRatio, s)
}

// MustParseRatio forcefully parses a string representing a ratio
func MustParseRatio(s string) Number {
	if res, err := ParseRatio(s); err != nil {
		panic(err)
	} else {
		return res
	}
}

// Cmp compares this Ratio to another Number
func (l *Ratio) Cmp(r Number) Comparison {
	if rr, ok := r.(*Ratio); ok {
		lb := (*big.Rat)(l)
		rb := (*big.Rat)(rr)
		return Comparison(lb.Cmp(rb))
	}
	pl, pr := purify(l, r)
	return pl.Cmp(pr)
}

// Add adds this Ratio to another Number
func (l *Ratio) Add(r Number) Number {
	if rr, ok := r.(*Ratio); ok {
		lb := (*big.Rat)(l)
		rb := (*big.Rat)(rr)
		res := new(big.Rat).Add(lb, rb)
		return (*Ratio)(res)
	}
	pl, pr := purify(l, r)
	return pl.Add(pr)
}

// Sub subtracts another Number from this Ratio
func (l *Ratio) Sub(r Number) Number {
	if rr, ok := r.(*Ratio); ok {
		lb := (*big.Rat)(l)
		rb := (*big.Rat)(rr)
		res := new(big.Rat).Sub(lb, rb)
		return (*Ratio)(res)
	}
	pl, pr := purify(l, r)
	return pl.Sub(pr)
}

// Mul multiplies this Ratio by another Number
func (l *Ratio) Mul(r Number) Number {
	if rr, ok := r.(*Ratio); ok {
		lb := (*big.Rat)(l)
		rb := (*big.Rat)(rr)
		res := new(big.Rat).Mul(lb, rb)
		return (*Ratio)(res)
	}
	pl, pr := purify(l, r)
	return pl.Mul(pr)
}

// Div divides this Ratio by another Number
func (l *Ratio) Div(r Number) Number {
	if rr, ok := r.(*Ratio); ok {
		lb := (*big.Rat)(l)
		rb := (*big.Rat)(rr)
		res := new(big.Rat).Quo(lb, rb)
		return (*Ratio)(res)
	}
	pl, pr := purify(l, r)
	return pl.Div(pr)
}

// Mod calculates the remainder of dividing this Ratio by another Number
func (l *Ratio) Mod(r Number) Number {
	if rr, ok := r.(*Ratio); ok {
		lb := (*big.Rat)(l)
		rb := (*big.Rat)(rr)
		lf, _ := lb.Float64()
		rf, _ := rb.Float64()
		return Float(math.Mod(lf, rf))
	}
	pl, pr := purify(l, r)
	return pl.Mod(pr)
}

// IsNaN returns whether this Ratio is not a number
func (*Ratio) IsNaN() bool {
	return false
}

// IsPosInf returns whether this Ratio represents positive infinity
func (*Ratio) IsPosInf() bool {
	return false
}

// IsNegInf returns whether this Ratio represents negative infinity
func (*Ratio) IsNegInf() bool {
	return false
}

// Equal compares this Ratio to another for equality
func (l *Ratio) Equal(r Value) bool {
	if r, ok := r.(*Ratio); ok {
		lb := (*big.Rat)(l)
		rb := (*big.Rat)(r)
		return lb.Cmp(rb) == 0
	}
	return false
}

// String converts this Ratio to a string
func (l *Ratio) String() string {
	return (*big.Rat)(l).String()
}

// HashCode returns a hash code for this Ratio
func (l *Ratio) HashCode() uint64 {
	br := (*big.Rat)(l)
	return br.Num().Uint64() * br.Denom().Uint64()
}

func (l *Ratio) float() Float {
	f, _ := (*big.Rat)(l).Float64()
	return Float(f)
}
