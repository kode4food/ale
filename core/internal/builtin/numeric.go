package builtin

import "github.com/kode4food/ale/data"

// Sub will subtract one number from the previous, in turn
var Sub = data.Applicative(func(args ...data.Value) data.Value {
	if len(args) == 1 {
		return data.Integer(-1).Mul(args[0].(data.Number))
	}
	var res = args[0].(data.Number)
	var rest = args[1:]
	for _, n := range rest {
		res = res.Sub(n.(data.Number))
	}
	return res
}, 1, data.OrMore)

// Div will divide one number by the next, in turn
var Div = data.Applicative(func(args ...data.Value) data.Value {
	var res = args[0].(data.Number)
	for _, n := range args[1:] {
		res = res.Div(n.(data.Number))
	}
	return res
}, 1, data.OrMore)

// Mod will produce the remainder of dividing one number by the next, in turn
var Mod = data.Applicative(func(args ...data.Value) data.Value {
	var res = args[0].(data.Number)
	for _, n := range args[1:] {
		res = res.Mod(n.(data.Number))
	}
	return res
}, 1, data.OrMore)

// Eq returns whether the provided numbers are equal
var Eq = data.Applicative(func(args ...data.Value) data.Value {
	var res = args[0].(data.Number)
	for _, n := range args[1:] {
		if res.Cmp(n.(data.Number)) != data.EqualTo {
			return data.False
		}
	}
	return data.True
}, 1, data.OrMore)

// Neq returns true if any of the numbers is not equal to the others
var Neq = data.Applicative(func(args ...data.Value) data.Value {
	if Eq.Call(args...) == data.True {
		return data.False
	}
	return data.True
}, 1, data.OrMore)

// Gt returns true if each number is greater than the previous
var Gt = data.Applicative(func(args ...data.Value) data.Value {
	var l = args[0].(data.Number)
	for _, v := range args[1:] {
		r := v.(data.Number)
		if l.Cmp(r) != data.GreaterThan {
			return data.False
		}
		l = r
	}
	return data.True
}, 1, data.OrMore)

// Gte returns true if each number is greater than or equal to the previous
var Gte = data.Applicative(func(args ...data.Value) data.Value {
	var l = args[0].(data.Number)
	for _, v := range args[1:] {
		r := v.(data.Number)
		cmp := l.Cmp(r)
		if !(cmp == data.GreaterThan || cmp == data.EqualTo) {
			return data.False
		}
		l = r
	}
	return data.True
}, 1, data.OrMore)

// Lt returns true if each number is less than the previous
var Lt = data.Applicative(func(args ...data.Value) data.Value {
	var l = args[0].(data.Number)
	for _, v := range args[1:] {
		r := v.(data.Number)
		if l.Cmp(r) != data.LessThan {
			return data.False
		}
		l = r
	}
	return data.True
}, 1, data.OrMore)

// Lte returns true if each number is less than or equal to the previous
var Lte = data.Applicative(func(args ...data.Value) data.Value {
	var l = args[0].(data.Number)
	for _, v := range args[1:] {
		r := v.(data.Number)
		cmp := l.Cmp(r)
		if !(cmp == data.LessThan || cmp == data.EqualTo) {
			return data.False
		}
		l = r
	}
	return data.True
}, 1, data.OrMore)

// IsPosInf returns true if the provided number represents positive infinity
var IsPosInf = data.Applicative(func(args ...data.Value) data.Value {
	if num, ok := args[0].(data.Number); ok {
		return data.Bool(num.IsPosInf())
	}
	return data.False
}, 1)

// IsNegInf returns true if the provided number represents negative infinity
var IsNegInf = data.Applicative(func(args ...data.Value) data.Value {
	if num, ok := args[0].(data.Number); ok {
		return data.Bool(num.IsNegInf())
	}
	return data.False
}, 1)

// IsNaN returns true if the provided value is not a number
var IsNaN = data.Applicative(func(args ...data.Value) data.Value {
	if num, ok := args[0].(data.Number); ok {
		return data.Bool(num.IsNaN())
	}
	return data.False
}, 1)
