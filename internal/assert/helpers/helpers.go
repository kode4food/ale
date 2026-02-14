package helpers

import (
	"fmt"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/data"
)

// B constructs a Bool
func B(v bool) data.Bool {
	return data.Bool(v)
}

// C constructs a Cons
func C(k, v ale.Value) *data.Cons {
	return data.NewCons(k, v)
}

// F constructs a Float
func F(f float64) data.Float {
	return data.Float(f)
}

// I will construct an Integer
func I(i int64) data.Integer {
	return data.Integer(i)
}

// K constructs a kwd
func K(s string) data.Keyword {
	return data.Keyword(s)
}

// L constructs a List
func L(args ...ale.Value) *data.List {
	return data.NewList(args...)
}

// LS constructs a Local Symbol
func LS(n string) data.Symbol {
	return data.Local(n)
}

// QS constructs a Qualified Symbol
func QS(d, l string) data.Symbol {
	return data.NewQualifiedSymbol(data.Local(l), data.Local(d))
}

// O constructs an Object from Pairs
func O(p ...data.Pair) *data.Object {
	return data.NewObject(p...)
}

// R constructs a Ratio
func R(num, den int64) data.Number {
	res, err := data.ParseRatio(fmt.Sprintf("%d/%d", num, den))
	if err != nil {
		panic(err)
	}
	return res
}

// S constructs a String
func S(s string) data.String {
	return data.String(s)
}

// V constructs a Vector
func V(args ...ale.Value) data.Vector {
	return args
}
