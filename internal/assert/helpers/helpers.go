package helpers

import "github.com/kode4food/ale/data"

// B constructs a Bool
func B(value bool) data.Bool {
	return data.Bool(value)
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

// N constructs a Name
func N(s string) data.Name {
	return data.Name(s)
}

// O constructs an Object
func O(args ...data.Value) data.Object {
	return data.NewObject(args...)
}

// S constructs a String
func S(s string) data.String {
	return data.String(s)
}

// V constructs a Vector
func V(args ...data.Value) data.Vector {
	return data.NewVector(args...)
}
