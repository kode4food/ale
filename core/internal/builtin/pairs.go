package builtin

import "github.com/kode4food/ale/data"

// Cons adds a value to the beginning of the provided Sequence
var Cons = data.Applicative(func(args ...data.Value) data.Value {
	car := args[0]
	cdr := args[1]
	if p, ok := cdr.(data.PrependerSequence); ok {
		return p.Prepend(car)
	}
	return data.NewCons(car, cdr)
}, 2)

// Car returns the first element of the provided Pair
var Car = data.Applicative(func(args ...data.Value) data.Value {
	return args[0].(data.Pair).Car()
}, 1)

// Cdr returns the first element of the provided Pair
var Cdr = data.Applicative(func(args ...data.Value) data.Value {
	return args[0].(data.Pair).Cdr()
}, 1)

// IsPair returns whether the provided value is a Pair
var IsPair = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.Pair)
	return data.Bool(ok)
}, 1)

// IsCons returns whether the provide value is a Cons cell
var IsCons = data.Applicative(func(args ...data.Value) data.Value {
	_, ok := args[0].(data.Cons)
	return data.Bool(ok)
}, 1)
