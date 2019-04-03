package helpers

import "gitlab.com/kode4food/ale/api"

// A constructs an Associative
func A(args ...api.Vector) api.Associative {
	return api.NewAssociative(args...)
}

// B constructs a Bool
func B(value bool) api.Bool {
	return api.Bool(value)
}

// F constructs a Float
func F(f float64) api.Float {
	return api.Float(f)
}

// I constructs an Integer
func I(i int64) api.Integer {
	return api.Integer(i)
}

// K constructs a Keyword
func K(s string) api.Keyword {
	return api.Keyword(s)
}

// L constructs a List
func L(args ...api.Value) *api.List {
	return api.NewList(args...)
}

// N constructs a Name
func N(s string) api.Name {
	return api.Name(s)
}

// S constructs a String
func S(s string) api.String {
	return api.String(s)
}

// V constructs a Vector
func V(args ...api.Value) api.Vector {
	return api.Vector(args)
}
