package helpers

import (
	"fmt"

	"github.com/kode4food/ale/data"
)

// B constructs a Bool
func B(value bool) data.Bool {
	return data.Bool(value)
}

// C constructs a Cons
func C(k, v data.Value) data.Cons {
	return data.NewCons(k, v)
}

// F constructs a Float
func F(f float64) data.Float {
	return data.Float(f)
}

// I constructs an Integer
func I(i int64) data.Integer {
	return data.Integer(i)
}

// K constructs a Keyword
func K(s string) data.Keyword {
	return data.Keyword(s)
}

// L constructs a List
func L(args ...data.Value) data.List {
	return data.NewList(args...)
}

// LS constructs a Local Symbol
func LS(n string) data.Symbol {
	return data.NewLocalSymbol(N(n))
}

// N constructs a Name
func N(s string) data.Name {
	return data.Name(s)
}

// O constructs an Object from Pairs
func O(p ...data.Pair) data.Object {
	return data.NewObject(p...)
}

// R constructs a Ratio
func R(num, den int64) data.Number {
	return data.MustParseRatio(fmt.Sprintf("%d/%d", num, den))
}

// S constructs a String
func S(s string) data.String {
	return data.String(s)
}

// V constructs a Vector
func V(args ...data.Value) data.Vector {
	return data.NewVector(args...)
}
