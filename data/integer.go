package data

import (
	"fmt"
	"math/big"
)

type (
	// Integer represents a 64-bit integer
	Integer int64

	// BigInt represents a multi-precision integer
	BigInt struct{ *big.Int }
)

// Error messages
const (
	ExpectedInteger = "value is not an integer: %s"
)

// ParseInteger attempts to parse a string representing an integer
func ParseInteger(s String) Number {
	ns := string(s)
	if res, ok := new(big.Int).SetString(ns, 0); ok {
		if res.IsInt64() {
			return Integer(res.Int64())
		}
		return &BigInt{Int: res}
	}
	panic(fmt.Errorf(ExpectedInteger, s))
}

// Cmp compares this Integer to another Number
func (l Integer) Cmp(r Number) Comparison {
	if ri, ok := r.(Integer); ok {
		if l > ri {
			return GreaterThan
		}
		if l < ri {
			return LessThan
		}
		return EqualTo
	}
	pl, pr := purify(l, r)
	return pl.Cmp(pr)
}

// Add adds this Integer to another Number
func (l Integer) Add(r Number) Number {
	if ri, ok := r.(Integer); ok {
		return l + ri
	}
	pl, pr := purify(l, r)
	return pl.Add(pr)
}

// Sub subtracts another Number from this Integer
func (l Integer) Sub(r Number) Number {
	if ri, ok := r.(Integer); ok {
		return l - ri
	}
	pl, pr := purify(l, r)
	return pl.Sub(pr)
}

// Mul multiples this Integer by another Number
func (l Integer) Mul(r Number) Number {
	if ri, ok := r.(Integer); ok {
		return l * ri
	}
	pl, pr := purify(l, r)
	return pl.Mul(pr)
}

// Div divides this Integer by another Number
func (l Integer) Div(r Number) Number {
	if ri, ok := r.(Integer); ok {
		return l / ri
	}
	pl, pr := purify(l, r)
	return pl.Div(pr)
}

// Mod calculates the remainder of dividing this Integer by another Number
func (l Integer) Mod(r Number) Number {
	if ri, ok := r.(Integer); ok {
		return l % ri
	}
	pl, pr := purify(l, r)
	return pl.Mod(pr)
}

// IsNaN tells you that this Integer is, in fact, a Number
func (Integer) IsNaN() bool {
	return false
}

// IsPosInf tells you that this Integer is not positive infinity
func (Integer) IsPosInf() bool {
	return false
}

// IsNegInf tells you that this Integer is not negative infinity
func (Integer) IsNegInf() bool {
	return false
}

// String converts this Integer to a string
func (l Integer) String() string {
	return fmt.Sprintf("%d", l)
}

// Cmp compares this BigInt to another Number
func (l *BigInt) Cmp(r Number) Comparison {
	if ri, ok := r.(*BigInt); ok {
		return Comparison(l.Int.Cmp(ri.Int))
	}
	lp, rp := purify(l, r)
	return lp.Cmp(rp)
}

// Add adds this BigInt to another Number
func (l *BigInt) Add(r Number) Number {
	if ri, ok := r.(*BigInt); ok {
		return &BigInt{
			Int: new(big.Int).Add(l.Int, ri.Int),
		}
	}
	lp, rp := purify(l, r)
	return lp.Add(rp)
}

// Sub subtracts another Number from this BigInt
func (l *BigInt) Sub(r Number) Number {
	if ri, ok := r.(*BigInt); ok {
		return &BigInt{
			Int: new(big.Int).Sub(l.Int, ri.Int),
		}
	}
	lp, rp := purify(l, r)
	return lp.Sub(rp)
}

// Mul multiples this BigInt by another Number
func (l *BigInt) Mul(r Number) Number {
	if ri, ok := r.(*BigInt); ok {
		return &BigInt{
			Int: new(big.Int).Mul(l.Int, ri.Int),
		}
	}
	lp, rp := purify(l, r)
	return lp.Mul(rp)
}

// Div divides this BigInt by another Number
func (l *BigInt) Div(r Number) Number {
	if ri, ok := r.(*BigInt); ok {
		return &BigInt{
			Int: new(big.Int).Quo(l.Int, ri.Int),
		}
	}
	lp, rp := purify(l, r)
	return lp.Div(rp)
}

// Mod calculates the remainder of dividing this BigInt by another Number
func (l *BigInt) Mod(r Number) Number {
	if ri, ok := r.(*BigInt); ok {
		return &BigInt{
			Int: new(big.Int).Rem(l.Int, ri.Int),
		}
	}
	lp, rp := purify(l, r)
	return lp.Mod(rp)
}

// IsNaN tells you that this BigInt is, in fact, a Number
func (*BigInt) IsNaN() bool {
	return false
}

// IsPosInf tells you that this BigInt is not positive infinity
func (*BigInt) IsPosInf() bool {
	return false
}

// IsNegInf tells you that this BigInt is not negative infinity
func (*BigInt) IsNegInf() bool {
	return false
}

// String converts this BigInt to a string
func (l *BigInt) String() string {
	return l.Int.String()
}
