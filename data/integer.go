package data

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand/v2"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/internal/types"
)

type (
	// Integer represents a 64-bit integer
	Integer int64

	// BigInt represents a multi-precision integer
	BigInt big.Int
)

const (
	// ErrExpectedInteger is raised when the value provided to ParseInteger
	// can't be interpreted as an integer
	ErrExpectedInteger = "value is not an integer: %s"

	// ErrDivideByZero is raised when an attempt is made to perform integer
	// division by zero
	ErrDivideByZero = "divide by zero"
)

var (
	intSalt = rand.Uint64()

	// compile-time checks for interface implementation
	_ interface {
		Hashed
		Number
		Procedure
		fmt.Stringer
	} = Integer(0)

	_ interface {
		Hashed
		Number
		fmt.Stringer
	} = (*BigInt)(nil)
)

// ParseInteger attempts to parse a string representing an integer
func ParseInteger(s string) (Number, error) {
	if res, ok := new(big.Int).SetString(s, 0); ok {
		if res.IsInt64() {
			return Integer(res.Int64()), nil
		}
		return (*BigInt)(res), nil
	}
	return nil, fmt.Errorf(ErrExpectedInteger, s)
}

func (l Integer) Cmp(r Number) Comparison {
	if ri, ok := r.(Integer); ok {
		return Comparison(cmp.Compare(l, ri))
	}
	pl, pr := purify(l, r)
	return pl.Cmp(pr)
}

func (l Integer) Add(r Number) Number {
	ri, ok := r.(Integer)
	if !ok {
		pl, pr := purify(l, r)
		return pl.Add(pr)
	}
	res := l + ri
	if (res^l) >= 0 || (res^ri) >= 0 {
		return res
	}
	lb := big.NewInt(int64(l))
	rb := big.NewInt(int64(ri))
	lb.Add(lb, rb)
	return (*BigInt)(lb)
}

func (l Integer) Sub(r Number) Number {
	ri, ok := r.(Integer)
	if !ok {
		pl, pr := purify(l, r)
		return pl.Sub(pr)
	}
	res := l - ri
	if (res^l) >= 0 || (res^^ri) >= 0 {
		return res
	}
	lb := big.NewInt(int64(l))
	rb := big.NewInt(int64(ri))
	lb.Sub(lb, rb)
	return (*BigInt)(lb)
}

func (l Integer) Mul(r Number) Number {
	ri, ok := r.(Integer)
	if !ok {
		pl, pr := purify(l, r)
		return pl.Mul(pr)
	}
	res := l * ri
	if (l != math.MinInt64 || ri >= 0) && (ri == 0 || res/ri == l) {
		return res
	}
	lb := big.NewInt(int64(l))
	rb := big.NewInt(int64(ri))
	lb.Mul(lb, rb)
	return (*BigInt)(lb)
}

func (l Integer) Div(r Number) Number {
	ri, ok := r.(Integer)
	if !ok {
		pl, pr := purify(l, r)
		return pl.Div(pr)
	}
	if ri == 0 {
		panic(errors.New(ErrDivideByZero))
	}
	res := big.NewRat(int64(l), int64(ri))
	return maybeWhole(res)
}

func (l Integer) Mod(r Number) Number {
	ri, ok := r.(Integer)
	if !ok {
		pl, pr := purify(l, r)
		return pl.Mod(pr)
	}
	if ri == 0 {
		panic(errors.New(ErrDivideByZero))
	}
	res := l % ri
	if (res < 0 && ri > 0) || (res > 0 && ri < 0) {
		return res + ri
	}
	return res
}

func (Integer) IsNaN() bool {
	return false
}

func (Integer) IsPosInf() bool {
	return false
}

func (Integer) IsNegInf() bool {
	return false
}

func (l Integer) Call(args ...ale.Value) ale.Value {
	m := args[0].(Indexed)
	res, ok := m.ElementAt(int(l))
	if !ok && len(args) > 1 {
		return args[1]
	}
	return res
}

func (l Integer) CheckArity(argc int) error {
	return CheckRangedArity(1, 2, argc)
}

func (l Integer) Equal(r ale.Value) bool {
	return l == r
}

func (l Integer) HashCode() uint64 {
	return intSalt ^ HashInt64(int64(l))
}

func (l Integer) String() string {
	return fmt.Sprintf("%d", l)
}

func (l Integer) Type() ale.Type {
	return types.MakeLiteral(types.BasicNumber, l)
}

func (l Integer) float() Float {
	return Float(l)
}

func (l Integer) bigInt() *BigInt {
	bi := big.NewInt(int64(l))
	return (*BigInt)(bi)
}

func (l Integer) ratio() *Ratio {
	r := new(big.Rat).SetFrac64(int64(l), 1)
	return (*Ratio)(r)
}

func (l *BigInt) Cmp(r Number) Comparison {
	if ri, ok := r.(*BigInt); ok {
		lb := (*big.Int)(l)
		rb := (*big.Int)(ri)
		return Comparison(lb.Cmp(rb))
	}
	lp, rp := purify(l, r)
	return lp.Cmp(rp)
}

func (l *BigInt) Add(r Number) Number {
	if ri, ok := r.(*BigInt); ok {
		lb := (*big.Int)(l)
		rb := (*big.Int)(ri)
		res := new(big.Int).Add(lb, rb)
		return maybeInteger(res)
	}
	lp, rp := purify(l, r)
	return lp.Add(rp)
}

func (l *BigInt) Sub(r Number) Number {
	if ri, ok := r.(*BigInt); ok {
		lb := (*big.Int)(l)
		rb := (*big.Int)(ri)
		res := new(big.Int).Sub(lb, rb)
		return maybeInteger(res)
	}
	lp, rp := purify(l, r)
	return lp.Sub(rp)
}

func (l *BigInt) Mul(r Number) Number {
	if ri, ok := r.(*BigInt); ok {
		lb := (*big.Int)(l)
		rb := (*big.Int)(ri)
		res := new(big.Int).Mul(lb, rb)
		return maybeInteger(res)
	}
	lp, rp := purify(l, r)
	return lp.Mul(rp)
}

func (l *BigInt) Div(r Number) Number {
	if ri, ok := r.(*BigInt); ok {
		lb := (*big.Int)(l)
		rb := (*big.Int)(ri)
		if rb.IsInt64() && rb.Int64() == 0 {
			panic(errors.New(ErrDivideByZero))
		}
		res := new(big.Int).Quo(lb, rb)
		return maybeInteger(res)
	}
	lp, rp := purify(l, r)
	return lp.Div(rp)
}

func (l *BigInt) Mod(r Number) Number {
	if ri, ok := r.(*BigInt); ok {
		lb := (*big.Int)(l)
		rb := (*big.Int)(ri)
		if rb.IsInt64() && rb.Int64() == 0 {
			panic(errors.New(ErrDivideByZero))
		}
		res := new(big.Int).Rem(lb, rb)
		return maybeInteger(res)
	}
	lp, rp := purify(l, r)
	return lp.Mod(rp)
}

func (*BigInt) IsNaN() bool {
	return false
}

func (*BigInt) IsPosInf() bool {
	return false
}

func (*BigInt) IsNegInf() bool {
	return false
}

func (l *BigInt) Equal(r ale.Value) bool {
	if r, ok := r.(*BigInt); ok {
		if l == r {
			return true
		}
		lb := (*big.Int)(l)
		rb := (*big.Int)(r)
		return lb.Cmp(rb) == 0
	}
	return false
}

func (l *BigInt) String() string {
	return (*big.Int)(l).String()
}

func (l *BigInt) Type() ale.Type {
	return types.MakeLiteral(types.BasicNumber, l)
}

func (l *BigInt) HashCode() uint64 {
	return intSalt ^ HashBytes((*big.Int)(l).Bytes())
}

func (l *BigInt) float() Float {
	bf := new(big.Float).SetInt((*big.Int)(l))
	f, _ := bf.Float64()
	return Float(f)
}

func (l *BigInt) ratio() *Ratio {
	r := new(big.Rat).SetInt((*big.Int)(l))
	return (*Ratio)(r)
}

func maybeInteger(bi *big.Int) Number {
	if bi.IsInt64() {
		return Integer(bi.Int64())
	}
	return (*BigInt)(bi)
}
