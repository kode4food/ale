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
	Ratio struct{ *big.Rat }
)

// Error messages
const (
	ExpectedFloat = "value is not a float: %s"
	ExpectedRatio = "value is not a ratio: %s"
)

// ParseFloat attempts to parse a string representing an float
func ParseFloat(s String) Number {
	ns := string(s)
	if res, err := strconv.ParseFloat(ns, 64); err == nil {
		return Float(res)
	}
	panic(fmt.Errorf(ExpectedFloat, s))
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

// IsNaN returns whether or not this Float is not a number
func (l Float) IsNaN() bool {
	return math.IsNaN(float64(l))
}

// IsPosInf returns whether or not this Float represents positive infinity
func (l Float) IsPosInf() bool {
	return math.IsInf(float64(l), 1)
}

// IsNegInf returns whether or not this Float represents negative infinity
func (l Float) IsNegInf() bool {
	return math.IsInf(float64(l), -1)
}

// String converts this Float to a string
func (l Float) String() string {
	return fmt.Sprintf("%g", l)
}

// ParseRatio attempts to parse a string representing a ratio
func ParseRatio(s String) Number {
	ns := string(s)
	if res, ok := new(big.Rat).SetString(ns); ok {
		return &Ratio{
			Rat: res,
		}
	}
	panic(fmt.Sprintf(ExpectedRatio, s))
}

// Cmp compares this *Ratio to another Number
func (l *Ratio) Cmp(r Number) Comparison {
	if rf, ok := r.(*Ratio); ok {
		return Comparison(l.Rat.Cmp(rf.Rat))
	}
	pl, pr := purify(l, r)
	return pl.Cmp(pr)
}

// Add adds this *Ratio to another Number
func (l *Ratio) Add(r Number) Number {
	if rf, ok := r.(*Ratio); ok {
		return &Ratio{
			Rat: new(big.Rat).Add(l.Rat, rf.Rat),
		}
	}
	pl, pr := purify(l, r)
	return pl.Add(pr)
}

// Sub subtracts another Number from this *Ratio
func (l *Ratio) Sub(r Number) Number {
	if rf, ok := r.(*Ratio); ok {
		return &Ratio{
			Rat: new(big.Rat).Sub(l.Rat, rf.Rat),
		}
	}
	pl, pr := purify(l, r)
	return pl.Sub(pr)
}

// Mul multiplies this *Ratio by another Number
func (l *Ratio) Mul(r Number) Number {
	if rf, ok := r.(*Ratio); ok {
		return &Ratio{
			Rat: new(big.Rat).Mul(l.Rat, rf.Rat),
		}
	}
	pl, pr := purify(l, r)
	return pl.Mul(pr)
}

// Div divides this *Ratio by another Number
func (l *Ratio) Div(r Number) Number {
	if rf, ok := r.(*Ratio); ok {
		return &Ratio{
			Rat: new(big.Rat).Quo(l.Rat, rf.Rat),
		}
	}
	pl, pr := purify(l, r)
	return pl.Div(pr)
}

// Mod calculates the remainder of dividing this *Ratio by another Number
func (l *Ratio) Mod(r Number) Number {
	if _, ok := r.(*Ratio); ok {
		panic("unsupported")
	}
	pl, pr := purify(l, r)
	return pl.Mod(pr)
}

// IsNaN returns whether or not this *Ratio is not a number
func (*Ratio) IsNaN() bool {
	return false
}

// IsPosInf returns whether or not this *Ratio represents positive infinity
func (*Ratio) IsPosInf() bool {
	return false
}

// IsNegInf returns whether or not this *Ratio represents negative infinity
func (*Ratio) IsNegInf() bool {
	return false
}

// String converts this *Ratio to a string
func (l *Ratio) String() string {
	return l.Rat.String()
}
