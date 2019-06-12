package builtin

import "gitlab.com/kode4food/ale/data"

// Cons adds a value to the beginning of the provided Sequence
func Cons(args ...data.Value) data.Value {
	car := args[0]
	cdr := args[1]
	if p, ok := cdr.(data.Prepender); ok {
		return p.Prepend(car)
	}
	return data.NewCons(car, cdr)
}

// Car returns the CAR register of the provided Pair
func Car(args ...data.Value) data.Value {
	return args[0].(data.Pair).Car()
}

// Cdr returns the CDR register of the provided Pair
func Cdr(args ...data.Value) data.Value {
	return args[0].(data.Pair).Cdr()
}

// IsPair returns whether or not the provided value is a Pair
func IsPair(args ...data.Value) data.Value {
	_, ok := args[0].(data.Pair)
	return data.Bool(ok)
}
