package data

import (
	"cmp"
	"fmt"
	"math"
	"math/big"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/types"
)

type (
	// Float represents a 64-bit floating point number
	Float float64

	// Ratio represents a number having a numerator and denominator
	Ratio big.Rat
)

const (
	// ErrExpectedFloat is raised when a call to ParseFloat can't properly
	// interpret its input as a floating point number
	ErrExpectedFloat = "value is not a float: %s"

	// ErrExpectedRatio is raised when a call to ParseRatio can't properly
	// interpret its input as a ratio
	ErrExpectedRatio = "value is not a ratio: %s"
)

var (
	ratSalt = rand.Uint64()
	one     = big.NewInt(1)

	// compile-time checks for interface implementation
	_ interface {
		Hashed
		Number
		fmt.Stringer
	} = Float(0)

	_ interface {
		Hashed
		Number
		fmt.Stringer
	} = (*Ratio)(nil)
)

// ParseFloat attempts to parse a string representing a float
func ParseFloat(s string) (Number, error) {
	res, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, fmt.Errorf(ErrExpectedFloat, s)
	}
	return Float(res), nil
}

func (l Float) Cmp(r Number) Comparison {
	if math.IsNaN(float64(l)) {
		return Incomparable
	}
	rf, ok := r.(Float)
	if !ok {
		pl, pr := purify(l, r)
		return pl.Cmp(pr)
	}
	if math.IsNaN(float64(rf)) {
		return Incomparable
	}
	return Comparison(cmp.Compare(l, rf))
}

func (l Float) Add(r Number) Number {
	if rf, ok := r.(Float); ok {
		return l + rf
	}
	pl, pr := purify(l, r)
	return pl.Add(pr)
}

func (l Float) Sub(r Number) Number {
	if rf, ok := r.(Float); ok {
		return l - rf
	}
	pl, pr := purify(l, r)
	return pl.Sub(pr)
}

func (l Float) Mul(r Number) Number {
	if rf, ok := r.(Float); ok {
		return l * rf
	}
	pl, pr := purify(l, r)
	return pl.Mul(pr)
}

func (l Float) Div(r Number) Number {
	if rf, ok := r.(Float); ok {
		return l / rf
	}
	pl, pr := purify(l, r)
	return pl.Div(pr)
}

func (l Float) Mod(r Number) Number {
	if rf, ok := r.(Float); ok {
		res := Float(math.Mod(float64(l), float64(rf)))
		if (res < 0 && rf > 0) || (res > 0 && rf < 0) {
			return res + rf
		}
		return res
	}
	pl, pr := purify(l, r)
	return pl.Mod(pr)
}

func (l Float) IsNaN() bool {
	return math.IsNaN(float64(l))
}

func (l Float) IsPosInf() bool {
	return math.IsInf(float64(l), 1)
}

func (l Float) IsNegInf() bool {
	return math.IsInf(float64(l), -1)
}

func (l Float) Equal(r ale.Value) bool {
	return l == r
}

func (l Float) String() string {
	i := int64(l)
	if Float(i) == l {
		return fmt.Sprintf("%d.0", i)
	}
	return strings.ToLower(fmt.Sprintf("%g", l))
}

func (l Float) Type() ale.Type {
	return types.MakeLiteral(types.BasicNumber, l)
}

func (l Float) HashCode() uint64 {
	return ratSalt ^ uint64(l)
}

// ParseRatio attempts to parse a string representing a ratio
func ParseRatio(s string) (Number, error) {
	if res, ok := new(big.Rat).SetString(s); ok {
		return maybeWhole(res), nil
	}
	return nil, fmt.Errorf(ErrExpectedRatio, s)
}

func (l *Ratio) Cmp(r Number) Comparison {
	if rr, ok := r.(*Ratio); ok {
		lb := (*big.Rat)(l)
		rb := (*big.Rat)(rr)
		return Comparison(lb.Cmp(rb))
	}
	pl, pr := purify(l, r)
	return pl.Cmp(pr)
}

func (l *Ratio) Add(r Number) Number {
	if rr, ok := r.(*Ratio); ok {
		lb := (*big.Rat)(l)
		rb := (*big.Rat)(rr)
		res := new(big.Rat).Add(lb, rb)
		return maybeWhole(res)
	}
	pl, pr := purify(l, r)
	return pl.Add(pr)
}

func (l *Ratio) Sub(r Number) Number {
	if rr, ok := r.(*Ratio); ok {
		lb := (*big.Rat)(l)
		rb := (*big.Rat)(rr)
		res := new(big.Rat).Sub(lb, rb)
		return maybeWhole(res)
	}
	pl, pr := purify(l, r)
	return pl.Sub(pr)
}

func (l *Ratio) Mul(r Number) Number {
	if rr, ok := r.(*Ratio); ok {
		lb := (*big.Rat)(l)
		rb := (*big.Rat)(rr)
		res := new(big.Rat).Mul(lb, rb)
		return maybeWhole(res)
	}
	pl, pr := purify(l, r)
	return pl.Mul(pr)
}

func (l *Ratio) Div(r Number) Number {
	if rr, ok := r.(*Ratio); ok {
		lb := (*big.Rat)(l)
		rb := (*big.Rat)(rr)
		res := new(big.Rat).Quo(lb, rb)
		return maybeWhole(res)
	}
	pl, pr := purify(l, r)
	return pl.Div(pr)
}

func (l *Ratio) Mod(r Number) Number {
	if rr, ok := r.(*Ratio); ok {
		lb := (*big.Rat)(l)
		rb := (*big.Rat)(rr)
		n := new(big.Int).Mul(lb.Num(), rb.Denom())
		d := new(big.Int).Mul(lb.Denom(), rb.Num())
		res := new(big.Rat).SetFrac(new(big.Int).Div(n, d), one)
		res = res.Mul(res, rb)
		res = res.Sub(lb, res)
		return maybeWhole(res)
	}
	pl, pr := purify(l, r)
	return pl.Mod(pr)
}

func (*Ratio) IsNaN() bool {
	return false
}

func (*Ratio) IsPosInf() bool {
	return false
}

func (*Ratio) IsNegInf() bool {
	return false
}

func (l *Ratio) Equal(r ale.Value) bool {
	if r, ok := r.(*Ratio); ok {
		if l == r {
			return true
		}
		lb := (*big.Rat)(l)
		rb := (*big.Rat)(r)
		return lb.Cmp(rb) == 0
	}
	return false
}

func (l *Ratio) String() string {
	return (*big.Rat)(l).String()
}

func (l *Ratio) Type() ale.Type {
	return types.MakeLiteral(types.BasicNumber, l)
}

func (l *Ratio) HashCode() uint64 {
	br := (*big.Rat)(l)
	return ratSalt ^ br.Num().Uint64() ^ br.Denom().Uint64()
}

func (l *Ratio) float() Float {
	f, _ := (*big.Rat)(l).Float64()
	return Float(f)
}

func maybeWhole(r *big.Rat) Number {
	if r.IsInt() {
		return maybeInteger(r.Num())
	}
	return (*Ratio)(r)
}
