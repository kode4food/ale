package data

import (
	"fmt"
	"math/big"
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
	CouldNotPromote = "could not purify: %v and %v"
)

// purify performs automatic contagion of operands
func purify(l, r Number) (Number, Number) {
	switch lt := l.(type) {
	case Integer:
		switch r.(type) {
		case Float:
			return Float(lt), r
		case *BigInt:
			return &BigInt{Int: big.NewInt(int64(lt))}, r
		case *Ratio:
			lr := &Ratio{
				Rat: new(big.Rat).SetFrac64(int64(lt), 1),
			}
			return lr, r
		}

	case Float:
		switch rt := r.(type) {
		case Integer:
			return l, Float(rt)
		case *BigInt:
			bf := new(big.Float).SetInt(rt.Int)
			f, _ := bf.Float64()
			return l, Float(f)
		case *Ratio:
			lf, _ := rt.Float64()
			return Float(lf), r
		}

	case *BigInt:
		switch rt := r.(type) {
		case Integer:
			return l, &BigInt{Int: big.NewInt(int64(rt))}
		case Float:
			bf := new(big.Float).SetInt(lt.Int)
			f, _ := bf.Float64()
			return Float(f), r
		case *Ratio:
			lr := &Ratio{
				Rat: new(big.Rat).SetInt(lt.Int),
			}
			return lr, r
		}

	case *Ratio:
		switch rt := r.(type) {
		case Integer:
			rr := &Ratio{
				Rat: new(big.Rat).SetFrac64(int64(rt), 1),
			}
			return l, rr
		case Float:
			f, _ := lt.Float64()
			return Float(f), r
		case *BigInt:
			ri := &Ratio{
				Rat: new(big.Rat).SetInt(rt.Int),
			}
			return l, ri
		}
	}

	panic(fmt.Errorf(CouldNotPromote, l, r))
}
