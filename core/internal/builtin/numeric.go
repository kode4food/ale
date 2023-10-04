package builtin

import "github.com/kode4food/ale/data"

// Neq returns true if any of the numbers is not equal to the others
var Neq = data.Applicative(func(args ...data.Value) data.Value {
	var res = args[0].(data.Number)
	for _, n := range args[1:] {
		if res.Cmp(n.(data.Number)) != data.EqualTo {
			return data.True
		}
	}
	return data.False
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
